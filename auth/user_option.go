/**
 * @Author: steven
 * @Description:
 * @File: auth_option
 * @Date: 18/05/24 14.50
 */

package auth

import (
	"github.com/evorts/kevlars/common"
	"github.com/evorts/kevlars/db"
	"github.com/evorts/kevlars/logger"
)

func UserAuthWithLogger(v logger.Manager) common.Option[clientManager] {
	return common.OptionFunc[clientManager](func(c *clientManager) {
		c.log = v
	})
}

func UserAuthWithDatabaseRead(db db.Manager) common.Option[clientManager] {
	return common.OptionFunc[clientManager](func(c *clientManager) {
		c.dbr = db
	})
}
