/**
 * @Author: steven
 * @Description:
 * @File: option
 * @Version: 1.0.0
 * @Date: 10/05/23 17.18
 */

package logger

import "github.com/evorts/kevlars/ctime"

type Option interface {
	apply(m *manager)
}

type option func(m *manager)

func (o option) apply(m *manager) {
	o(m)
}

func WithServiceName(v string) Option {
	return option(func(m *manager) {
		m.name = v
	})
}

func WithTZTimeFormatter(v ctime.TimeZone) Option {
	return option(func(m *manager) {
		if f, err := newTZFormatter(v); err == nil {
			m.f = f
			m.l.SetFormatter(m.f)
		}
	})
}
