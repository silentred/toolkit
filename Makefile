# Ensure GOPATH is set before running build process.
ifeq "$(GOPATH)" ""
  $(error Please set the environment variable GOPATH before running `make`)
endif

CURDIR := $(shell pwd)
GO        := go
GOBUILD   := $(GO) build
GOTEST    := $(GO) test

OS        := "`uname -s`"
LINUX     := "Linux"
MAC       := "Darwin"
PACKAGES  := $$(go list ./...| grep -vE 'vendor|tests')
FILES     := $$(find . -name '*.go' | grep -vE 'vendor')
TARGET	  := "toolkit"
LDFLAGS   += -X "main.BuildTS=$(shell date -u '+%Y-%m-%d %I:%M:%S')"
LDFLAGS   += -X "main.GitHash=$(shell git rev-parse HEAD)"

test:
	$(GOTEST) $(PACKAGES) -cover

build:
	$(GOBUILD) -ldflags '$(LDFLAGS)' -o $(TARGET)

dev: test build

clean:
	rm $(TARGET)