/**
 * @Author: steven
 * @Description:
 * @File: app_context
 * @Date: 10/05/24 19.03
 */

package scaffold

import (
	"context"
)

type IContext interface {
	WithValue(key any, value any) IApplication
	WithContext(ctx context.Context) IApplication

	Context() context.Context
}

func (app *Application) WithValue(key any, value any) IApplication {
	app.startContext = context.WithValue(app.startContext, key, value)
	return app
}

func (app *Application) WithContext(ctx context.Context) IApplication {
	app.startContext = ctx
	return app
}

func (app *Application) Context() context.Context {
	return app.startContext
}
