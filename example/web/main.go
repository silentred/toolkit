package main

import (
	"flag"

	redis "gopkg.in/redis.v5"

	"github.com/labstack/echo/v4"

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
	return nil
}

func initRoute(app service.Application) error {
	if web, ok := app.(service.WebApplication); ok {
		web.GetRouter().GET("/", func(ctx echo.Context) error {
			var ret string
			if _, ok := app.Get("mysql").(*service.MysqlManager); ok {
				ret += "Mysql init \n"
			}
			if _, ok := app.Get("redis").(*redis.Client); ok {
				ret += "Redis init \n"
			}
			if _, ok := app.Get("app.web").(*service.WebApp); ok {
				ret += "WebApp init \n"
			}

			ret += "hello world \n"
			return ctx.String(200, ret)
		})
	}
	return nil
}
