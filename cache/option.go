/**
 * @Author: steven
 * @Description:
 * @File: Option
 * @Version: 1.0.0
 * @Date: 20/06/23 11.44
 */

package cache

import "github.com/evorts/kevlars/telemetry"

type Option interface {
	apply(*manager)
}

type option func(*manager)

func (o option) apply(m *manager) {
	o(m)
}

func WithPrefix(v string) Option {
	return option(func(m *manager) {
		m.prefix = v
	})
}

func WithTLSFile(v bool, certFile, keyFile, serverName string) Option {
	return option(func(m *manager) {
		m.useTLS = v
		m.certFile = certFile
		m.keyFile = keyFile
		m.serverName = serverName
	})
}

func WithTLSB64(v bool, certB64, keyB64, serverName string) Option {
	return option(func(m *manager) {
		m.useTLS = v
		m.certB64 = certB64
		m.keyB64 = keyB64
		m.serverName = serverName
	})
}

func WithTelemetry(tm telemetry.Manager) Option {
	return option(func(m *manager) {
		m.tm = tm
	})
}

func WithScope(v string) Option {
	return option(func(m *manager) {
		m.scope = v
	})
}
