/**
 * @Author: steven
 * @Description:
 * @File: main
 * @Date: 22/12/23 09.25
 */

package main

import (
	"github.com/evorts/kevlars/scaffold"
	"github.com/evortstech/kevlars/examples/simple_rest/app"
	"github.com/labstack/echo/v4"
)

func main() {
	scaffold.
		NewApp(scaffold.WithScope("restful_api")).
		WithDatabases().
		RunRestApiUseEcho(func(sa *scaffold.Application, e *echo.Echo) {
			app.Start(sa, e)
		})
}
