/**
 * @Author: steven
 * @Description:
 * @File: options
 * @Date: 20/12/23 23.29
 */

package scaffold

import (
	"context"
	"time"
)

type Option interface {
	apply(app *Application)
}

type optionFunc func(app *Application)

func (o optionFunc) apply(app *Application) {
	o(app)
}

func WithName(v string) Option {
	return optionFunc(func(app *Application) {
		app.name = v
	})
}

func WithVersion(v string) Option {
	return optionFunc(func(app *Application) {
		app.version = v
	})
}

func WithScope(v string) Option {
	return optionFunc(func(app *Application) {
		app.scope = v
	})
}

func WithPortRestConfigPath(path string) Option {
	return optionFunc(func(app *Application) {
		app.portRest = app.Config().GetInt(path)
	})
}

func WithPortGrpcConfigPath(path string) Option {
	return optionFunc(func(app *Application) {
		app.portGrpc = app.Config().GetInt(path)
	})
}

func WithStartContext(v context.Context) Option {
	return optionFunc(func(app *Application) {
		app.startContext = v
	})
}

func WithGracefulTimeout(v time.Duration) Option {
	return optionFunc(func(app *Application) {
		app.gracefulTimeout = v
	})
}

func WithParseRequestId(v bool) Option {
	return optionFunc(func(app *Application) {
		app.parseRequestId = v
	})
}

func WithCustomRequestIdGenerator(v func() string) Option {
	return optionFunc(func(app *Application) {
		app.requestIdGenerator = v
	})
}
