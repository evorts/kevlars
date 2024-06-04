/**
 * @Author: steven
 * @Description:
 * @File: option
 * @Version: 1.0.0
 * @Date: 29/08/23 23.23
 */

package ctime

import "github.com/evorts/kevlars/common"

func WithCustomTimeZone(v TimeZone) common.Option[manager] {
	return common.OptionFunc[manager](func(m *manager) {
		m.tz = v
	})
}
