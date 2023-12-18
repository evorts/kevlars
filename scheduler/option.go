/**
 * @Author: steven
 * @Description:
 * @File: option
 * @Date: 08/12/23 07.00
 */

package scheduler

import "github.com/evorts/kevlars/logger"

type Option interface {
	apply(m *manager)
}

type option func(*manager)

func (o option) apply(m *manager) {
	o(m)
}

func WithLogger(log logger.Manager) Option {
	return option(func(m *manager) {
		m.log = log
	})
}

func WithName(v string) Option {
	return option(func(m *manager) {
		m.name = v
	})
}
