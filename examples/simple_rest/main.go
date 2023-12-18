/**
 * @Author: steven
 * @Description:
 * @File: main
 * @Date: 22/12/23 09.25
 */

package main

import (
	"github.com/evorts/kevlars/scaffold"
	"github.com/labstack/echo/v4"
)

func main() {
	scaffold.
		NewApp(scaffold.WithScope("restful_api")).
		WithDatabases().
		RunRestApiUseEcho(func(app *scaffold.Application, e *echo.Echo) {
			// do something such as routing and needed process
			// all resources that instantiate above (e.g. WithDatabase) are available under `app`
			// for example:
			// if we want to get config value:
			// app.Config().GetString("key")
			// if we want to use database connection:
			// app.DefaultDB().Exec(ctx, q, args...)
		})
}
