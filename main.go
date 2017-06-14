package main

import (
	"fmt"
	"os"

	cli "gopkg.in/urfave/cli.v1"
)

var (
	// Version of toolkit
	Version = "v0.0.1"
	// BuildTS is timestamp of build
	BuildTS = "None"
	// GitHash is commit version of build
	GitHash = "None"

	logo = `
 _____           _ _    _ _
|_   _|__   ___ | | | _(_) |_
  | |/ _ \ / _ \| | |/ / | __|
  | | (_) | (_) | |   <| | |_
  |_|\___/ \___/|_|_|\_\_|\__| 

Version: %s
GitHash: %s
BuildTS: %s

`

	sourcePath string
)

func main() {
	cli.VersionPrinter = versionPrinter

	app := cli.NewApp()
	app.Version = Version
	app.Usage = "A toolkit of the Toolkit"

	app.Commands = []cli.Command{
		cli.Command{
			Name:      "new",
			ShortName: "n",
			Usage:     "Create a new project",
			UsageText: "For example: toolkit new app_name",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "src",
					Usage:       "Source path in the GOPATH. For example: github.com/igetserver",
					Value:       "github.com/igetserver",
					Destination: &sourcePath,
				},
			},
			Action: newAction,
		},
	}
	app.Run(os.Args)
}

func versionPrinter(ctx *cli.Context) {
	fmt.Fprintf(ctx.App.Writer, logo, ctx.App.Version, GitHash, BuildTS)
	cli.ShowAppHelp(ctx)
}

func newAction(ctx *cli.Context) error {
	fmt.Println(sourcePath)
	return nil
}
