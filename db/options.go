/**
 * @Author: steven
 * @Description:
 * @File: options
 * @Version: 1.0.0
 * @Date: 20/09/23 11.27
 */

package db

import "github.com/evorts/kevlars/telemetry"

type Option interface {
	apply(m *manager)
}

type option func(m *manager)

func (o option) apply(m *manager) {
	o(m)
}

func WithMaxOpenConnection(v int) Option {
	return option(func(m *manager) {
		m.maxOpenConnection = v
	})
}

func WithMaxIdleConnection(v int) Option {
	return option(func(m *manager) {
		m.maxIdleConnection = v
	})
}

func WithScope(v string) Option {
	return option(func(m *manager) {
		m.scope = v
	})
}

func WithOTelConnect(v bool) Option {
	return option(func(m *manager) {
		m.oTelOpenConnect = v
	})
}

func WithTelemetry(v telemetry.Manager) Option {
	return option(func(m *manager) {
		m.tm = v
	})
}

func WithTelemetryEnabled(v bool) Option {
	return option(func(m *manager) {
		m.telemetryEnabled = v
	})
}
