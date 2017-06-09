package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"time"

	"flag"

	"encoding/json"

	"path/filepath"

	"github.com/labstack/echo"
	elog "github.com/labstack/gommon/log"
	"github.com/silentred/echorus"
	cfg "github.com/silentred/toolkit/config"
	"github.com/silentred/toolkit/util"
	"github.com/silentred/toolkit/util/container"
	"github.com/silentred/toolkit/util/rotator"
	"github.com/silentred/toolkit/util/strings"
	"github.com/spf13/viper"
)

var (
	// AppMode is App's running envirenment. Valid values are dev and prod
	AppMode    string
	ConfigFile string
	LogPath    string
)

func init() {
	flag.StringVar(&AppMode, "mode", "", "RunMode of the application: dev or prod")
	flag.StringVar(&ConfigFile, "cfg", "", "absolute path of config file")
	flag.StringVar(&LogPath, "logPath", ".", "logPath is where log file will be")
}

type Application interface {
}

type WebApp interface {
}

// HookFunc when app starting and tearing down
type HookFunc func(*App) error

// App represents the application
type App struct {
	Store    *container.Map
	Injector container.Injector
	Router   *echo.Echo

	// TODO use interface for logger
	loggers map[string]*echorus.Echorus
	Config  cfg.AppConfig

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
		loggers:  make(map[string]*echorus.Echorus),
	}
	// register App itself
	app.Set("app", app, nil)
	return app
}

// Logger of name
func (app *App) Logger(name string) *echorus.Echorus {
	if name == "" {
		return app.loggers["default"]
	}
	if l, ok := app.loggers[name]; ok {
		return l
	}

	return nil
}

func (app *App) SetLogger(name string, logger *echorus.Echorus) bool {
	if _, ok := app.loggers[name]; ok {
		return false
	}
	app.loggers[name] = logger
	return true
}

// DefaultLogger gets default logger
func (app *App) DefaultLogger() *echorus.Echorus {
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

// InitConfig in format of toml
func (app *App) InitConfig() {
	// use viper to resolve config.toml
	if ConfigFile == "" {
		var configName = app.getConfigFile()
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

	// TODO: session config
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

	app.Config = config

	// hook
	app.runConfigHooks()
}

func (app *App) getConfigFile() string {
	var configName = "config"
	if AppMode != "" {
		configName = fmt.Sprintf("%s.%s", "config", AppMode)
	}
	return configName
}

func (app *App) InitLogger() {
	var level elog.Lvl
	switch app.Config.Mode {
	case cfg.ModeProd:
		level = elog.INFO
	default:
		level = elog.DEBUG
	}
	// new default Logger
	defaultLogger := NewLogger(app.Config.Name, level, app.Config.Log)

	// set logger
	app.loggers["default"] = defaultLogger
	app.Route.Logger = defaultLogger

	// hook
	app.runLoggerHooks()
}

func (app *App) runConfigHooks() {
	var err error
	for _, f := range app.configHooks {
		if f != nil {
			err = f(app)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func (app *App) runLoggerHooks() {
	var err error
	for _, f := range app.loggerHooks {
		if f != nil {
			err = f(app)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func (app *App) initService() {
	// hoook
	var err error
	for _, f := range app.serviceHooks {
		if f != nil {
			err = f(app)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func (app *App) initRoute() {
	// hook
	var err error
	for _, f := range app.routeHooks {
		if f != nil {
			err = f(app)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

}

func (app *App) shutdown() {
	var err error
	for _, f := range app.shutdownHooks {
		if f != nil {
			err = f(app)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

// RegisterConfigHook at initConfig
func (app *App) RegisterConfigHook(hooks ...HookFunc) {
	app.configHooks = append(app.configHooks, hooks...)
}

func (app *App) RegisterLoggerHook(hooks ...HookFunc) {
	app.loggerHooks = append(app.loggerHooks, hooks...)
}

func (app *App) RegisterServiceHook(hooks ...HookFunc) {
	app.serviceHooks = append(app.serviceHooks, hooks...)
}

func (app *App) RegisterRouteHook(hooks ...HookFunc) {
	app.routeHooks = append(app.routeHooks, hooks...)
}

func (app *App) RegisterShutdownHook(hooks ...HookFunc) {
	app.shutdownHooks = append(app.shutdownHooks, hooks...)
}

func (app *App) Init() {
	app.InitConfig()
	app.InitLogger()
	app.initService()
	app.initRoute()
}

// Start running the application
func (app *App) Start() {
	app.Init()
	//app.route.Start(fmt.Sprintf(":%d", app.config.Port))
	app.graceStart()
	app.shutdown()
}

func (app *App) graceStart() error {
	// Start server
	go func() {
		if err := app.Route.Start(fmt.Sprintf(":%d", app.Config.Port)); err != nil {
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
	if err := app.Route.Shutdown(ctx); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// NewLogger return a new
func NewLogger(appName string, level elog.Lvl, config cfg.LogConfig) *echorus.Echorus {
	// new default Logger
	var writer io.Writer
	var spliter rotator.Spliter
	var err error

	if config.Suffix == "" {
		config.Suffix = "log"
	}

	switch config.Providor {
	case cfg.ProvidorFile:
		if config.RotateEnable {
			switch config.RotateMode {
			case cfg.RotateByDay:
				spliter = rotator.NewDaySpliter()
			case cfg.RotateBySize:
				limitSize, err := strings.ParseByteSize(config.RotateLimit) // 100 MB
				if err != nil {
					log.Fatal(err)
				}
				spliter = rotator.NewSizeSpliter(uint64(limitSize))
			default:
				log.Fatalf("invalid RotateMode: %s", config.RotateMode)
			}

			writer = rotator.NewFileRotator(config.LogPath, appName, config.Suffix, spliter)
		} else {
			writer, err = os.Open(filepath.Join(config.LogPath, appName+"."+config.Suffix))
			if err != nil {
				log.Fatal(err)
			}
		}
	default:
		writer = os.Stdout
	}

	logger := echorus.NewLogger()
	logger.SetPrefix(appName)
	logger.SetFormat(echorus.TextFormat)
	logger.SetOutput(writer)
	logger.SetLevel(level)

	return logger
}
