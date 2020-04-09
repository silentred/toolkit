package filter

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/silentred/toolkit/util"
)

var (
	ErrToken = util.NewError(400403, "invalid token")
)

type (
	// TokenConfig defines the config for Auth middleware.
	TokenConfig struct {
		Logger util.Logger
		Secret string
	}
)

// AuthToken middleware
func AuthToken(logger util.Logger, secret string) echo.MiddlewareFunc {
	config := TokenConfig{
		Logger: logger,
		Secret: secret,
	}
	return newTokenFilter(config)
}

func newTokenFilter(config TokenConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			// verify the token
			var token = c.QueryParam("token")
			var nonce = c.QueryParam("nonce")
			var method = c.Request().Method
			var ts = c.QueryParam("ts")
			var hashStr = fmt.Sprintf("%s,%s,%s", method, nonce, ts)

			bytes := md5.Sum(util.Slice(hashStr))
			md5Str := fmt.Sprintf("%x", util.String(bytes[:]))
			rightSign := base64.StdEncoding.EncodeToString(util.Slice(md5Str))

			if rightSign != token {
				config.Logger.Errorf("str:%s ; token:%s ; rightSign:%s", hashStr, token, rightSign)
				return c.JSON(200, ErrToken)
			}

			if err = next(c); err != nil {
				c.Error(err)
			}
			return
		}
	}
}
