package util

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/labstack/echo"
	elog "github.com/labstack/gommon/log"
	"github.com/silentred/echorus"
	cfg "github.com/silentred/toolkit/config"
	"github.com/silentred/toolkit/util/rotator"
	"github.com/silentred/toolkit/util/strings"
)

type Logger echo.Logger

// NewLogger return a new
func NewLogger(appName string, level elog.Lvl, config cfg.LogConfig) Logger {
	// new default Logger
	var writer io.Writer
	var spliter rotator.Spliter
	var err error

	if config.Suffix == "" {
		config.Suffix = "log"
	}

	switch config.Providor {
	case cfg.ProvidorFile:
		if config.RotateEnable {
			switch config.RotateMode {
			case cfg.RotateByDay:
				spliter = rotator.NewDaySpliter()
			case cfg.RotateBySize:
				limitSize, err := strings.ParseByteSize(config.RotateLimit) // 100 MB
				if err != nil {
					log.Fatal(err)
				}
				spliter = rotator.NewSizeSpliter(uint64(limitSize))
			default:
				log.Fatalf("invalid RotateMode: %s", config.RotateMode)
			}

			writer = rotator.NewFileRotator(config.LogPath, appName, config.Suffix, spliter)
		} else {
			writer, err = os.Open(filepath.Join(config.LogPath, appName+"."+config.Suffix))
			if err != nil {
				log.Fatal(err)
			}
		}
	default:
		writer = os.Stdout
	}

	logger := echorus.NewLogger()
	logger.SetPrefix(appName)
	logger.SetFormat(echorus.TextFormat)
	logger.SetOutput(writer)
	logger.SetLevel(level)

	return logger
}
