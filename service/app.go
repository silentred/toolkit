package service

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	elog "github.com/labstack/gommon/log"
	cfg "github.com/silentred/toolkit/config"
	"github.com/silentred/toolkit/util"
	"github.com/silentred/toolkit/util/container"
	"github.com/spf13/viper"
)

// HookFunc when app starting and tearing down
type HookFunc func(Application) error

// HookType for hook
type HookType byte

const (
	// ConfigHook is a hook type for config
	ConfigHook HookType = iota
	// LoggerHook is a hook type for logger
	LoggerHook
	// ServiceHook is a hook type for service
	ServiceHook
	// RouterHook is a hook type for router
	RouterHook
	// ShutdownHook is a hook type for shutting down
	ShutdownHook
)

var (
	// AppMode is App's running envirenment. Valid values are dev and prod
	AppMode string
	// ConfigFile is the absolute path of config file
	ConfigFile string
	// LogPath is where log file will be
	LogPath string

	_ Application = &App{}
)

func init() {
	flag.StringVar(&AppMode, "mode", "", "RunMode of the application: dev or prod")
	flag.StringVar(&ConfigFile, "cfg", "config.toml", "absolute path of config file")
	flag.StringVar(&LogPath, "logPath", ".", "logPath is where log file will be")
}

// Application interface represents a service application
type Application interface {
	// Container
	Set(key string, object interface{}, ifacePtr interface{})
	Get(key string) interface{}
	Inject(object interface{}) error
	// hook
	GetHook(HookType) *[]HookFunc
	// load from file
	LoadConfig(mode string) *cfg.AppConfig
	SetConfig(*cfg.AppConfig)
	GetConfig() *cfg.AppConfig
	// logger
	DefaultLogger() util.Logger
	Logger(name string) (util.Logger, error)
	SetLogger(string, util.Logger)
	// init
	Initialize()

	RegisterHook(HookType, ...HookFunc)
}

// App represents the application
type App struct {
	Store    *container.Map
	Injector container.Injector
	//Router   *echo.Echo

	//loggers map[string]util.Logger
	Config *cfg.AppConfig

	configHooks   []HookFunc
	loggerHooks   []HookFunc
	serviceHooks  []HookFunc
	routeHooks    []HookFunc
	shutdownHooks []HookFunc
}

// NewApp gets a new application
func NewApp() App {
	app := App{
		Store:    &container.Map{},
		Injector: container.NewInjector(),
		//loggers:  make(map[string]util.Logger),
	}
	// register App itself
	app.Set("app", app, new(Application))
	// register default hooks
	app.RegisterHook(ConfigHook, initConfig)
	app.RegisterHook(LoggerHook, initLogger)
	app.RegisterHook(ServiceHook, initMySQL, initRedis)

	return app
}

// Logger of name
func (app *App) Logger(name string) (util.Logger, error) {
	var key = fmt.Sprintf("logger.%s", name)
	if l, ok := app.Get(key).(util.Logger); ok {
		return l, nil
	}
	return nil, fmt.Errorf("not found key: %s", key)
}

// SetLogger set logger
func (app *App) SetLogger(name string, logger util.Logger) {
	var key = fmt.Sprintf("logger.%s", name)
	app.Set(key, logger, nil)
}

// DefaultLogger gets default logger
func (app *App) DefaultLogger() util.Logger {
	l, err := app.Logger("default")
	if err != nil {
		log.Fatalf("default logger is not present, err: %v", err)
	}
	return l
}

// Set object into app.Store and Map it into app.Injector
func (app *App) Set(key string, object interface{}, ifacePtr interface{}) {
	app.Store.Set(key, object)
	if ifacePtr != nil {
		app.Injector.MapTo(object, ifacePtr)
	} else {
		app.Injector.Map(object)
	}
}

// Get object from app.Store
func (app *App) Get(key string) interface{} {
	return app.Store.Get(key)
}

// Inject dependencies to the object. Please MAKE SURE that the dependencies should be stored at app.Injector
// before this method is called. Please use app.Set() to make this happen.
func (app *App) Inject(object interface{}) error {
	return app.Injector.Apply(object)
}

// GetHook returns hook slice by type
func (app *App) GetHook(ht HookType) *[]HookFunc {
	var hook *[]HookFunc

	switch ht {
	case ConfigHook:
		hook = &app.configHooks
	case LoggerHook:
		hook = &app.loggerHooks
	case ServiceHook:
		hook = &app.serviceHooks
	case RouterHook:
		hook = &app.routeHooks
	case ShutdownHook:
		hook = &app.shutdownHooks
	}

	return hook
}

// SetConfig sets config ptr
func (app *App) SetConfig(config *cfg.AppConfig) {
	app.Config = config
}

// GetConfig gets config ptr
func (app *App) GetConfig() *cfg.AppConfig {
	return app.Config
}

// Initialize application
func (app *App) Initialize() {
	initialize(app)
}

// LoadConfig by mode from file
func (app *App) LoadConfig(mode string) *cfg.AppConfig {
	// use viper to resolve config.toml
	if ConfigFile == "" {
		var configName = getConfigFile(mode)
		viper.AddConfigPath(".")
		viper.AddConfigPath(util.SelfDir())
		viper.SetConfigName(configName)
	} else {
		viper.SetConfigFile(ConfigFile)
	}

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	// make AppConfig; set data from viper
	config := cfg.AppConfig{}
	appConfig(&config)

	// log config
	loggerConfig(&config)

	// mysql config
	mysqlConfig(&config)

	// redis config
	redisConfig(&config)

	return &config
}

func appConfig(c *cfg.AppConfig) {
	c.Name = viper.GetString("app.name")
	c.Mode = viper.GetString("app.runMode")
	c.Port = viper.GetInt("app.port")
}

func loggerConfig(c *cfg.AppConfig) {
	l := cfg.LogConfig{
		Name:         "default",
		LogPath:      viper.GetString("app.logPath"),
		Providor:     viper.GetString("app.logProvider"),
		RotateEnable: viper.GetBool("app.logRotate"),
		RotateMode:   viper.GetString("app.logRotateType"),
		RotateLimit:  viper.GetString("app.logLimit"),
		Suffix:       viper.GetString("app.logExt"),
	}
	c.Log = l
}

func mysqlConfig(c *cfg.AppConfig) {
	mysql := cfg.MysqlConfig{}
	mysqlConfig := viper.Get("mysql")
	mysqlConfigBytes, err := json.Marshal(mysqlConfig)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(mysqlConfigBytes, &mysql.Instances)
	if err != nil {
		log.Fatal(err)
	}
	mysql.Ping = viper.GetBool("mysql_manager.ping")
	mysql.InitMySQL = viper.GetBool("mysql_manager.init")
	c.Mysql = mysql
}

func redisConfig(c *cfg.AppConfig) {
	redis := cfg.RedisInstance{
		Host: viper.GetString("redis.host"),
		Port: viper.GetInt("redis.port"),
		Db:   viper.GetInt("redis.db"),
		Pwd:  viper.GetString("redis.password"),
	}
	redisConfig := cfg.RedisConfig{
		Ping:          viper.GetBool("redis_manager.ping"),
		InitRedis:     viper.GetBool("redis_manager.init"),
		RedisInstance: redis,
	}

	c.Redis = redisConfig
}

func getConfigFile(mode string) string {
	var configName = "config"
	if mode != "" {
		configName = fmt.Sprintf("%s.%s", "config", mode)
	}
	return configName
}

// initConfig loads config from toml file and setConfig
func initConfig(app Application) error {
	app.SetConfig(app.LoadConfig(AppMode))
	return nil
}

func initLogger(app Application) error {
	var level elog.Lvl
	switch app.GetConfig().Mode {
	case cfg.ModeProd:
		level = elog.INFO
	default:
		level = elog.DEBUG
	}
	// new default Logger
	defaultLogger := util.NewLogger(app.GetConfig().Name, level, app.GetConfig().Log)
	app.SetLogger("default", defaultLogger)

	return nil
}

func runHooks(ht HookType, app Application) {
	var err error
	var hook = app.GetHook(ht)
	if *hook != nil {
		for _, f := range *hook {
			if f != nil {
				if err = f(app); err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}

// RegisterHook in application's starting process
func (app *App) RegisterHook(ht HookType, hooks ...HookFunc) {
	var hook = app.GetHook(ht)
	*hook = append(*hook, hooks...)
}

// initialize the objects
func initialize(app Application) {
	runHooks(ConfigHook, app)
	runHooks(LoggerHook, app)
	runHooks(ServiceHook, app)
	runHooks(RouterHook, app)
}
