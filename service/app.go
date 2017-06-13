package service

import (
	"fmt"
	"log"

	"flag"

	"encoding/json"

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
	ConfigHook HookType = iota
	LoggerHook
	ServiceHook
	RouterHook
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
	// load from file
	LoadConfig(mode string) *cfg.AppConfig
	SetConfig(*cfg.AppConfig)
	GetConfig() *cfg.AppConfig
	// logger
	DefaultLogger() (util.Logger, error)
	Logger(name string) (util.Logger, error)
	SetLogger(string, util.Logger)
	// initialize
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
func (app *App) DefaultLogger() (util.Logger, error) {
	return app.Logger("logger.default")
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

// SetConfig sets config ptr
func (app *App) SetConfig(config *cfg.AppConfig) {
	app.Config = config
}

// GetConfig gets config ptr
func (app *App) GetConfig() *cfg.AppConfig {
	return app.Config
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
	config.Name = viper.GetString("app.name")
	config.Mode = viper.GetString("app.runMode")
	config.Port = viper.GetInt("app.port")

	// log config
	loggerConfig(&config)

	// mysql config
	mysqlConfig(&config)

	// redis config
	redisConfig(&config)

	return &config
}

func loggerConfig(c *cfg.AppConfig) {
	l := cfg.LogConfig{}
	l.Name = "default"
	l.LogPath = viper.GetString("app.logPath")
	l.Providor = viper.GetString("app.logProvider")
	l.RotateEnable = viper.GetBool("app.logRotate")
	l.RotateMode = viper.GetString("app.logRotateType")
	l.RotateLimit = viper.GetString("app.logLimit")
	l.Suffix = viper.GetString("app.logExt")
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
	c.Mysql = mysql
}

func redisConfig(c *cfg.AppConfig) {
	redis := cfg.RedisInstance{}
	redis.Host = viper.GetString("redis.host")
	redis.Port = viper.GetInt("redis.port")
	redis.Db = viper.GetInt("redis.db")
	redis.Pwd = viper.GetString("redis.password")
	redis.Ping = viper.GetBool("redis.ping")
	c.Redis = redis
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

func initMySQL(app Application) error {
	return nil
}

func initRedis(app Application) error {
	return nil
}

func (app *App) runHooks(ht HookType) {
	var err error
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

	if *hook != nil {
		for _, f := range *hook {
			if f != nil {
				err = f(app)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}

// RegisterHook in application's starting process
func (app *App) RegisterHook(ht HookType, hooks ...HookFunc) {
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

	*hook = append(*hook, hooks...)
}

// Initialize the objects
func (app *App) Initialize() {
	app.runHooks(ConfigHook)
	app.runHooks(LoggerHook)
	app.runHooks(ServiceHook)
	app.runHooks(RouterHook)
}
