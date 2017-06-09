package filter

import (
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/silentred/echorus"
)

type (
	// LoggerConfig defines the config for Logger middleware.
	LoggerConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper Skipper
		Logger  *echorus.Echorus
		Format  logrus.Formatter
	}
)

func NewConfig(logger *echorus.Echorus) LoggerConfig {
	return LoggerConfig{
		Skipper: DefaultSkipper,
		Logger:  logger,
		Format:  echorus.TextFormat,
	}
}

// Logger returns a middleware that logs HTTP requests.
func Logger(logger *echorus.Echorus) echo.MiddlewareFunc {
	return LoggerWithConfig(NewConfig(logger))
}

// LoggerWithConfig returns a Logger middleware with config.
func LoggerWithConfig(config LoggerConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultSkipper
	}
	if config.Format != nil {
		config.Logger.SetFormat(config.Format)
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if config.Skipper(c) {
				return next(c)
			}

			req := c.Request()
			res := c.Response()
			start := time.Now()
			if err = next(c); err != nil {
				c.Error(err)
			}

			stop := time.Now()

			p := req.URL.Path
			if p == "" {
				p = "/"
			}

			cl := req.Header.Get(echo.HeaderContentLength)
			if cl == "" {
				cl = "0"
			}
			json := log.JSON{
				"time_unix":   strconv.FormatInt(time.Now().Unix(), 10),
				"remote_ip":   c.RealIP(),
				"host":        req.Host,
				"uri":         req.RequestURI,
				"method":      req.Method,
				"path":        p,
				"user_agent":  req.UserAgent(),
				"status":      res.Status,
				"latency":     strconv.FormatInt(int64(stop.Sub(start)), 10),
				"latency_str": stop.Sub(start).String(),
				"bytes_in":    cl,
				"bytes_out":   strconv.FormatInt(res.Size, 10),
			}

			config.Logger.Infoj(json)
			return
		}
	}
}
