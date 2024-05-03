/**
 * @Author: steven
 * @Description:
 * @File: valkey_option
 * @Date: 30/04/24 18.53
 */

package inmemory

import "github.com/evorts/kevlars/telemetry"

func ValKeyWithPassword(v string) Option[valkeyManager] {
	return option[valkeyManager](func(m *valkeyManager) {
		m.pwd = v
	})
}

func ValKeyWithDB(v int) Option[valkeyManager] {
	return option[valkeyManager](func(m *valkeyManager) {
		m.db = v
	})
}

func ValKeyWithPrefix(prefix string) Option[valkeyManager] {
	return option[valkeyManager](func(m *valkeyManager) {
		m.prefix = prefix
	})
}

func ValKeyWithTLSFile(v bool, certFile, keyFile, serverName string) Option[valkeyManager] {
	return option[valkeyManager](func(m *valkeyManager) {
		m.useTLS = v
		m.certFile = certFile
		m.keyFile = keyFile
		m.serverName = serverName
	})
}

func ValKeyWithTLSB64(v bool, certB64, keyB64, serverName string) Option[valkeyManager] {
	return option[valkeyManager](func(m *valkeyManager) {
		m.useTLS = v
		m.certB64 = certB64
		m.keyB64 = keyB64
		m.serverName = serverName
	})
}

func ValKeyWithTelemetry(v telemetry.Manager) Option[valkeyManager] {
	return option[valkeyManager](func(m *valkeyManager) {
		m.tm = v
	})
}

func ValKeyWithScope(v string) Option[valkeyManager] {
	return option[valkeyManager](func(m *valkeyManager) {
		m.scope = v
	})
}

func NewValKeyOptions() []Option[valkeyManager] {
	return make([]Option[valkeyManager], 0)
}
