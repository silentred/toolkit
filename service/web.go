package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo"
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
	return app
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
	app.RegisterHook(LoggerHook, initRouterLogger)
	app.Initialize()
	app.graceStart()
	app.runHooks(ShutdownHook)
}

func (app *WebApp) graceStart() error {
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

func initRouterLogger(app Application) error {
	if webApp, ok := app.(WebApplication); ok {
		webApp.GetRouter().Logger = app.DefaultLogger()
		return nil
	}
	return fmt.Errorf("app is not WebApplication")
}
