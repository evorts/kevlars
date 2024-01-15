/**
 * @Author: steven
 * @Description:
 * @File: echo_static
 * @Date: 15/01/24 23.09
 */

package midware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/random"
)

func EchoStaticMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			staticSessionId := ctx.Request().Header.Get("X-Static-Id")
			if len(staticSessionId) > 0 {
				// check redis if session id exist
			}
			if len(staticSessionId) < 1 {
				// generate new response header when it's not exist
				staticSessionId = random.String(32)
			}
			// @todo: save to redis
			ctx.Response().Before(func() {
				ctx.Response().Header().Set("X-Static-Id", staticSessionId)
			})
			return next(ctx)
		}
	}
}
