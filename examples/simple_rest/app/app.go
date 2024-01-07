/**
 * @Author: steven
 * @Description:
 * @File: app
 * @Date: 08/01/24 06.15
 */

package app

import (
	"github.com/evorts/kevlars/contracts"
	"github.com/evorts/kevlars/scaffold"
	"github.com/labstack/echo/v4"
	"strconv"
)

func Start(app *scaffold.Application, e *echo.Echo) {
	tdm := NewTodo(app.DefaultDB())
	// create new todo item
	e.PUT("/todo", func(c echo.Context) error {
		body := new(CreateTodoRequest)
		if err := c.Bind(body); err != nil {
			return c.JSON(contracts.NewResponseBadRequest("invalid request"))
		}
		id, err := tdm.Create(c.Request().Context(), Todo{
			Title:       body.Title,
			Description: body.Description,
		})
		if err != nil {
			return c.JSON(contracts.NewResponseInternalServerError(err.Error()))
		}
		return c.JSON(contracts.NewResponseOK("OK", map[string]int{"id": id}))
	})
	// update todo item
	e.POST("/todo", func(c echo.Context) error {
		body := new(UpdateTodoRequest)
		if err := c.Bind(body); err != nil || body.Id < 1 {
			return c.JSON(contracts.NewResponseBadRequest("invalid request"))
		}
		err := tdm.UpdateById(c.Request().Context(), Todo{
			Id:          body.Id,
			Title:       body.Title,
			Description: body.Description,
		})
		if err != nil {
			return c.JSON(contracts.NewResponseInternalServerError(err.Error()))
		}
		return c.JSON(contracts.NewResponseOK("OK", make(map[string]string)))
	})
	e.GET("/todo/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(contracts.NewResponseBadRequest("invalid request"))
		}
		rs, errGet := tdm.GetById(c.Request().Context(), id)
		if errGet != nil {
			return c.JSON(contracts.NewResponseInternalServerError(errGet.Error()))
		}
		return c.JSON(contracts.NewResponseOK("OK", rs))
	})
	e.DELETE("/todo/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(contracts.NewResponseBadRequest("invalid request"))
		}
		err = tdm.DeleteById(c.Request().Context(), id)
		if err != nil {
			return c.JSON(contracts.NewResponseInternalServerError(err.Error()))
		}
		return c.JSON(contracts.NewResponseOK("OK", make(map[string]string)))
	})
}
