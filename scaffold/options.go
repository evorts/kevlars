/**
 * @Author: steven
 * @Description:
 * @File: options
 * @Date: 20/12/23 23.29
 */

package scaffold

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

func WithScope(v string) Option {
	return optionFunc(func(app *Application) {
		app.scope = v
	})
}
