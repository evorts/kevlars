/**
 * @Author: steven
 * @Description:
 * @File: smtp_option
 * @Date: 31/05/24 08.23
 */

package csmtp

import (
	"github.com/evorts/kevlars/common"
	"time"
)

func SmtpWithTimeout(timeout *time.Duration) common.Option[client] {
	return common.OptionFunc[client](func(c *client) {
		c.timeout = timeout
	})
}
