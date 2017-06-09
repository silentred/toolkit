package main

import (
	"flag"
	"log"

	"github.com/labstack/echo"
	"github.com/silentred/kassadin"
	"github.com/silentred/kassadin/db"
)

func main() {
	flag.Parse()

	app := kassadin.NewApp()
	app.RegisterConfigHook(initConfig)
	app.RegisterRouteHook(initRoute)
	app.RegisterServiceHook(initService)
	app.Start()
}

func initConfig(app *kassadin.App) error {
	return nil
}

func initService(app *kassadin.App) error {
	mm, err := db.NewMysqlManager(app, app.Config.Mysql)
	if err != nil {
		log.Fatal(err)
	}
	app.Set("mysql", mm, nil)

	redis := db.NewRedisClient(app.Config.Redis)
	app.Set("redis", redis, nil)

	return nil
}

func initRoute(app *kassadin.App) error {
	app.Route.GET("/", func(ctx echo.Context) error {
		return ctx.String(200, "hello world")
	})

	return nil
}
