# Echo logger with rotator

## Usage

```
package main

import ()
	"github.com/labstack/echo"
	"github.com/silentred/rotator"
)

// Echo is the web engine
var Echo *echo.Echo

func init() {
	Echo = echo.New()
}

func main() {
	initLogger()
    Echo.Logger.Info("test")
	Echo.Start(":8090")
}

func initLogger() {
	path := "storage/log"
	appName := "app"
	limitSize := 100 << 20 // 100MB

	r := rotator.NewFileSizeRotator(path, appName, "log", limitSize)
	Echo.Logger.SetOutput(r) 
	Echo.Logger.SetLevel(elog.WARN)
}

```
