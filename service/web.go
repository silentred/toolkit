package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
)

var (
	_ WebApplication = &WebApp{}
)

// WebApplication interface represents a web application
type WebApplication interface {
	Application
	GetRouter() *echo.Echo
	SetRouter(*echo.Echo)
	ListenAndServe()
}

// WebApp is the concrete type of WebApplication
type WebApp struct {
	App
	Router *echo.Echo
}

// NewWebApp returns a new web app
func NewWebApp() *WebApp {
	app := &WebApp{
		App:    NewApp(),
		Router: echo.New(),
	}
	app.Set("app.web", app, new(WebApplication))
	return app
}

// Initialize web application
func (app *WebApp) Initialize() {
	app.RegisterHook(LoggerHook, initRouterLogger)
	initialize(app)
}

// SetRouter sets router
func (app *WebApp) SetRouter(r *echo.Echo) {
	app.Router = r
}

// GetRouter gets router
func (app *WebApp) GetRouter() *echo.Echo {
	return app.Router
}

// ListenAndServe the web application
func (app *WebApp) ListenAndServe() {
	graceStart(app)
	runHooks(ShutdownHook, app)
}

func graceStart(app *WebApp) error {
	// Start server
	go func() {
		if err := app.Router.Start(fmt.Sprintf("%s:%d", app.Config.Host, app.Config.Port)); err != nil {
			log.Fatal(err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 3 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := app.Router.Shutdown(ctx); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func initRouterLogger(app Application) error {
	if webApp, ok := app.(*WebApp); ok {
		webApp.GetRouter().Logger = app.DefaultLogger()
		return nil
	}
	return fmt.Errorf("initializing Echo.Logger: app is not *WebApp")
}
