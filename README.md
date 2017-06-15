# toolkit

Toolkit of web API and RPC services

![](https://img.shields.io/badge/language-golang-blue.svg)
![](https://img.shields.io/badge/license-MIT-000000.svg)
![](https://img.shields.io/github/tag/silentred/toolkit.svg)
[![codebeat badge](https://codebeat.co/badges/644e898b-f0cb-4a05-b701-cdfb790c37e5)](https://codebeat.co/projects/github-com-silentred-toolkit-master)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/ae66252a0e4b45719e08037088d07863)](https://www.codacy.com/app/silentred/toolkit?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=silentred/toolkit&amp;utm_campaign=Badge_Grade)

## Build

```
go get github.com/silentred/toolkit
cd $GOPATH/src/github.com/silentred/toolkit
make build
./toolkit -v
 _____           _ _    _ _
|_   _|__   ___ | | | _(_) |_
  | |/ _ \ / _ \| | |/ / | __|
  | | (_) | (_) | |   <| | |_
  |_|\___/ \___/|_|_|\_\_|\__|

Version: v0.0.1
GitHash: 7975f569095d28c979c70f192533fa35a507b0c7
BuildTS: 2017-06-15 03:44:02

NAME:
   toolkit - A toolkit of the Toolkit

USAGE:
   toolkit [global options] command [command options] [arguments...]

VERSION:
   v0.0.1

COMMANDS:
     new, n   Create a new project
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

## Usage

### Create new project

```
$ ./toolkit new myapp
# Step1: creating dirs and files
create file: /Users/jason/Desktop/projects/GOPATH/src/gitlab.luojilab.com/igetserver/myapp/main.go
create file: /Users/jason/Desktop/projects/GOPATH/src/gitlab.luojilab.com/igetserver/myapp/service/echo.go
create file: /Users/jason/Desktop/projects/GOPATH/src/gitlab.luojilab.com/igetserver/myapp/Makefile
create file: /Users/jason/Desktop/projects/GOPATH/src/gitlab.luojilab.com/igetserver/myapp/config.toml
create file successful
# Step2: run following commands to start the app
cd $GOPATH/gitlab.luojilab.com/igetserver/myapp
git init && git add * && git commit -m "init commit"
make build && ./myapp
Have fun!
```

Follow the guide and run the application.

### Run and watch

Coming later.
