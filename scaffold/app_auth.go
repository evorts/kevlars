/**
 * @Author: steven
 * @Description:
 * @File: app_auth
 * @Date: 13/05/24 18.14
 */

package scaffold

import (
	"github.com/evorts/kevlars/auth"
	"io/fs"
)

type IAuth interface {
	WithClientAuthorization(migrate func() fs.FS) IApplication

	ClientAuth() auth.ClientManager
}

func (app *Application) WithClientAuthorization(migrate func() fs.FS) IApplication {
	app.authClient = auth.NewClientManager(
		app.DefaultDB(),
		auth.ClientWithDatabaseRead(app.DefaultDBR()),
	)
	if migrate != nil {
		app.authClient.AddOptions(auth.ClientWithExecuteMigration(migrate(), func() bool {
			return app.Config().GetBool("auth.client.dynamic_migration_enabled")
		}))
	}
	if v := app.Config().GetBool("auth.client.lazy_load_data"); v {
		app.authClient.AddOptions(auth.ClientWithLazyLoadData(app.Config().GetBool("auth.client.lazy_load_data")))
	}
	app.authClient.MustInit()
	return app
}

func (app *Application) ClientAuth() auth.ClientManager {
	return app.authClient
}
