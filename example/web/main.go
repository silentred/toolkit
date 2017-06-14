package main

import (
	"flag"
	"log"

	"github.com/labstack/echo"

	"github.com/silentred/toolkit/db"
	"github.com/silentred/toolkit/service"
)

func main() {
	flag.Parse()

	app := service.NewWebApp()
	app.RegisterHook(service.ConfigHook, initConfig)
	app.RegisterHook(service.RouterHook, initRoute)
	app.RegisterHook(service.ServiceHook, initService)
	app.Initialize()
	app.ListenAndServe()
}

func initConfig(app service.Application) error {
	return nil
}

func initService(app service.Application) error {
	mm, err := db.NewMysqlManager(app, app.GetConfig().Mysql)
	if err != nil {
		log.Fatal(err)
	}
	app.Set("mysql", mm, nil)

	redis := db.NewRedisClient(app.GetConfig().Redis)
	app.Set("redis", redis, nil)

	return nil
}

func initRoute(app service.Application) error {
	if web, ok := app.(service.WebApplication); ok {
		web.GetRouter().GET("/", func(ctx echo.Context) error {
			return ctx.String(200, "hello world")
		})
	}
	return nil
}
