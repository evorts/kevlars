/**
 * @Author: steven
 * @Description:
 * @File: option
 * @Version: 1.0.0
 * @Date: 29/08/23 23.23
 */

package ts

type Option interface {
	apply(m *manager)
}

type option func(m *manager)

func (o option) apply(m *manager) {
	o(m)
}

func WithCustomTimeZone(v TimeZone) Option {
	return option(func(m *manager) {
		m.tz = v
	})
}
