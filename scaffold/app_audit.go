/**
 * @Author: steven
 * @Description:
 * @File: app_audit
 * @Date: 13/05/24 18.28
 */

package scaffold

import (
	"github.com/evorts/kevlars/audit"
)

type IAudit interface {
	WithAuditLog() IApplication

	AuditLog() audit.Manager
}

func (app *Application) WithAuditLog() IApplication {
	app.audit = audit.New(app.DefaultDB()).MustInit()
	return app
}

func (app *Application) AuditLog() audit.Manager {
	return app.audit
}
