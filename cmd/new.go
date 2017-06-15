package cmd

import (
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/fatih/color"
)

var (
	configTmpl = `[app]
runMode = "dev"
name = "{{.AppName}}"
port = 18080

logProvider = "file"
logPath = "/tmp"
logRotate = true
logRotateType = "day"
logLimit = "100MB"

[mysql_manager]
init = true
ping = false

[redis_manager]
init = true
ping = false

[[mysql]]
name="master"
host="localhost"
port=3306
user="root"
db="test"
password=""
read_only=false

[[mysql]]
name="slave-01"
host="localhost"
port=3306
user="root"
db="test"
password=""
read_only=true

[redis]
host="localhost"
port= 6379
db = 0
ping = false
	`

	makefileTmpl = `# Ensure GOPATH is set before running build process.
ifeq "$(GOPATH)" ""
  $(error Please set the environment variable GOPATH before running 'make')
endif

CURDIR := $(shell pwd)
GO        := go
GOBUILD   := $(GO) build
GOTEST    := $(GO) test

OS        := "` + "`uname -s`" + `"
LINUX     := "Linux"
MAC       := "Darwin"
PACKAGES  := $$(go list ./...| grep -vE 'vendor|tests')
FILES     := $$(find . -name '*.go' | grep -vE 'vendor')
TARGET	  := "{{.AppName}}"
LDFLAGS   += -X "main.BuildTS=$(shell date -u '+%Y-%m-%d %I:%M:%S')"
LDFLAGS   += -X "main.GitHash=$(shell git rev-parse HEAD)"

test:
	$(GOTEST) $(PACKAGES) -cover

build:
	$(GOBUILD) -ldflags '$(LDFLAGS)' -o $(TARGET)

dev: test build

clean:
	rm $(TARGET)
	`

	mainTmpl = `
package main

import (
	"flag"

	"github.com/labstack/echo"
	"github.com/silentred/toolkit/service"
	svc "{{.SrcPath}}/{{.AppName}}/service"
)

var (
	GitHash = "None"
	BuildTS = "None"
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
	// init hello service
	hello := &svc.HelloService{}
	app.Inject(hello)
	app.Set("hello", hello, nil)
	return nil
}

func initRoute(app service.Application) error {
	if web, ok := app.(service.WebApplication); ok {
		web.GetRouter().GET("/", func(ctx echo.Context) error {
			var ret = ctx.QueryParam("say")
			if ret == "" {
				ret = "Hello world"
			}
			if h, ok := web.Get("hello").(*svc.HelloService); ok {
				ret = h.SayHello(ret)
			}
			return ctx.String(200, ret)
		})
	}
	return nil
}
	`

	echoTmpl = `
package service

import (
	"fmt"

	"github.com/silentred/toolkit/service"
)

var (
	tmpl =` + "`" + `
( %s )
( )
 ,__, |    | 
(oo)\\|    |___
(__)\\|    |   )\\_
      |    |_w |  \\
      |    |  ||   *

		Cower....` +
		"`" + `
)

type HelloService struct {
	WebApp service.WebApplication` + "`inject:\"app.web\"`" + `
}

func (h *HelloService) SayHello(thought string) string {
	return fmt.Sprintf(tmpl, thought)
}
	`
)

type file struct {
	path string
	data string
}

// RunNew runs new command
func RunNew(path, appName string) {
	gopath, has := os.LookupEnv("GOPATH")
	if !has {
		log.Fatal("Please set GOPATH first")
	}

	makeFiles(gopath, path, appName)

	color.Yellow("# Step2: run following commands to start the app")
	color.Blue("cd $GOPATH/%s/%s", path, appName)
	color.Blue("git init && git add * && git commit -m \"init commit\" ")
	color.Blue("make build && ./%s", appName)
	color.Green("Have fun!")
}

func makeFiles(gopath, path, appName string) {
	var data = map[string]string{
		"SrcPath": path,
		"AppName": appName,
	}

	// file path -> template
	sources := map[string]*template.Template{}
	sources["Makefile"], _ = template.New("make").Parse(makefileTmpl)
	sources["config.toml"], _ = template.New("make").Parse(configTmpl)
	sources["main.go"], _ = template.New("make").Parse(mainTmpl)
	sources["service/echo.go"], _ = template.New("make").Parse(echoTmpl)

	appDir := fmt.Sprintf("%s/src/%s/%s", gopath, path, appName)
	dirs := []string{
		appDir,
		appDir + "/service",
	}
	for _, val := range dirs {
		err := os.MkdirAll(val, 0755|os.ModeSticky)
		if err != nil {
			log.Fatalf("create dir %s: %v", val, err)
		}
	}

	color.Yellow("# Step1: creating dirs and files")
	for key, val := range sources {
		path := fmt.Sprintf("%s/src/%s/%s/%s", gopath, path, appName, key)
		f, err := os.Create(path)
		if err != nil {
			log.Fatalf("create file %s: %v", path, err)
		}
		// log create file
		color.Blue("create file: %s", path)
		val.Execute(f, data)
	}

	color.Green("create file successful")
}
