/**
 * @Author: steven
 * @Description:
 * @File: model
 * @Date: 20/12/23 23.31
 */

package scaffold

import "github.com/labstack/echo/v4"

const (
	DefaultKey = "default"
)

type route struct {
	method      string
	path        string
	handlerEcho echo.HandlerFunc
}
