/**
 * @Author: steven
 * @Description:
 * @File: auth_option
 * @Date: 18/05/24 14.50
 */

package auth

import (
	"github.com/evorts/kevlars/audit"
	"github.com/evorts/kevlars/common"
	"github.com/evorts/kevlars/db"
	"github.com/evorts/kevlars/inmemory"
	"github.com/evorts/kevlars/logger"
)

func UserAuthWithLogger(v logger.Manager) common.Option[userManager] {
	return common.OptionFunc[userManager](func(c *userManager) {
		c.log = v
	})
}

func UserAuthWithDatabaseRead(db db.Manager) common.Option[userManager] {
	return common.OptionFunc[userManager](func(c *userManager) {
		c.dbr = db
	})
}

func UserAuthWithAuditManager(am audit.Manager) common.Option[userManager] {
	return common.OptionFunc[userManager](func(u *userManager) {
		u.audit = am
	})
}

func UserAuthWithInMemoryManager(im inmemory.Manager) common.Option[userManager] {
	return common.OptionFunc[userManager](func(u *userManager) {
		u.im = im
	})
}
