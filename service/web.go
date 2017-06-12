package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"flag"

	"encoding/json"

	"github.com/labstack/echo"
	elog "github.com/labstack/gommon/log"
	cfg "github.com/silentred/toolkit/config"
	"github.com/silentred/toolkit/util"
	"github.com/silentred/toolkit/util/container"
	"github.com/spf13/viper"
)

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
	AppMode    string
	ConfigFile string
	LogPath    string

	_ Application = &App{}
)

func init() {
	flag.StringVar(&AppMode, "mode", "", "RunMode of the application: dev or prod")
	flag.StringVar(&ConfigFile, "cfg", "config.toml", "absolute path of config file")
	flag.StringVar(&LogPath, "logPath", ".", "logPath is where log file will be")
}

type Application interface {
	// Container
	Set(key string, object interface{}, ifacePtr interface{})
	Get(key string) interface{}
	Inject(object interface{}) error
	// load from file
	LoadConfig(mode string) *cfg.AppConfig
	SetConfig(*cfg.AppConfig)
	GetConfig() *cfg.AppConfig

	DefaultLogger() util.Logger
	Logger(name string) util.Logger

	RegisterHook(HookType, ...HookFunc)
}

// WebApp interface for web application
type WebApp interface {
	Application
	GetRouter()
	SetRouter()
	ListenAndServe()
}

// HookFunc when app starting and tearing down
type HookFunc func(*App) error

// App represents the application
type App struct {
	Store    *container.Map
	Injector container.Injector
	Router   *echo.Echo

	loggers map[string]util.Logger
	Config  *cfg.AppConfig

	configHooks   []HookFunc
	loggerHooks   []HookFunc
	serviceHooks  []HookFunc
	routeHooks    []HookFunc
	shutdownHooks []HookFunc
}

// NewApp gets a new application
func NewApp() *App {
	app := &App{
		Store:    &container.Map{},
		Injector: container.NewInjector(),
		Router:   echo.New(),
		loggers:  make(map[string]util.Logger),
	}
	// register App itself
	app.Set("app", app, nil)
	return app
}

// Logger of name
func (app *App) Logger(name string) util.Logger {
	if name == "" {
		return app.loggers["default"]
	}
	if l, ok := app.loggers[name]; ok {
		return l
	}

	return nil
}

// SetLogger set logger
func (app *App) SetLogger(name string, logger util.Logger) bool {
	if _, ok := app.loggers[name]; ok {
		return false
	}
	app.loggers[name] = logger
	return true
}

// DefaultLogger gets default logger
func (app *App) DefaultLogger() util.Logger {
	return app.Logger("")
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
	l := cfg.LogConfig{}
	l.Name = "default"
	l.LogPath = viper.GetString("app.logPath")
	l.Providor = viper.GetString("app.logProvider")
	l.RotateEnable = viper.GetBool("app.logRotate")
	l.RotateMode = viper.GetString("app.logRotateType")
	l.RotateLimit = viper.GetString("app.logLimit")
	l.Suffix = viper.GetString("app.logExt")
	config.Log = l

	// mysql config
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
	config.Mysql = mysql

	// redis config
	redis := cfg.RedisInstance{}
	redis.Host = viper.GetString("redis.host")
	redis.Port = viper.GetInt("redis.port")
	redis.Db = viper.GetInt("redis.db")
	redis.Pwd = viper.GetString("redis.password")
	config.Redis = redis

	return &config
}

// initConfig loads config from toml file and setConfig
func initConfig(app *App) {
	app.SetConfig(app.LoadConfig(AppMode))
}

func getConfigFile(mode string) string {
	var configName = "config"
	if mode != "" {
		configName = fmt.Sprintf("%s.%s", "config", mode)
	}
	return configName
}

func initLogger(app *App) {
	var level elog.Lvl
	switch app.Config.Mode {
	case cfg.ModeProd:
		level = elog.INFO
	default:
		level = elog.DEBUG
	}
	// new default Logger
	defaultLogger := util.NewLogger(app.Config.Name, level, app.Config.Log)

	// set logger
	app.loggers["default"] = defaultLogger
	app.Router.Logger = defaultLogger

	// hook
	//app.runLoggerHooks()
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

func (app *App) Init() {
	app.runHooks(ConfigHook)
	app.runHooks(LoggerHook)
	app.runHooks(ServiceHook)
	app.runHooks(RouterHook)
}

// Start running the application
func (app *App) Start() {
	app.Init()
	//app.route.Start(fmt.Sprintf(":%d", app.config.Port))
	app.graceStart()
	app.runHooks(ShutdownHook)
}

func (app *App) graceStart() error {
	// Start server
	go func() {
		if err := app.Router.Start(fmt.Sprintf(":%d", app.Config.Port)); err != nil {
			log.Fatal(err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.Router.Shutdown(ctx); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
