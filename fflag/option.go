/**
 * @Author: steven
 * @Description:
 * @File: option
 * @Date: 17/05/24 08.38
 */

package fflag

import (
	"github.com/evorts/kevlars/common"
	"github.com/evorts/kevlars/db"
)

func WithDatabaseRead(db db.Manager) common.Option[manager] {
	return common.OptionFunc[manager](func(c *manager) {
		c.dbr = db
	})
}
