/**
 * @Author: steven
 * @Description:
 * @File: redis_option
 * @Date: 30/04/24 18.43
 */

package inmemory

import "github.com/evorts/kevlars/telemetry"

func RedisWithPassword(v string) Option[redisManager] {
	return option[redisManager](func(m *redisManager) {
		m.pwd = v
	})
}

func RedisWithDB(v int) Option[redisManager] {
	return option[redisManager](func(m *redisManager) {
		m.db = v
	})
}

func RedisWithPrefix(v string) Option[redisManager] {
	return option[redisManager](func(m *redisManager) {
		m.prefix = v
	})
}

func RedisWithTLSFile(v bool, certFile, keyFile, serverName string) Option[redisManager] {
	return option[redisManager](func(m *redisManager) {
		m.useTLS = v
		m.certFile = certFile
		m.keyFile = keyFile
		m.serverName = serverName
	})
}

func RedisWithTLSB64(v bool, certB64, keyB64, serverName string) Option[redisManager] {
	return option[redisManager](func(m *redisManager) {
		m.useTLS = v
		m.certB64 = certB64
		m.keyB64 = keyB64
		m.serverName = serverName
	})
}

func RedisWithTelemetry(tm telemetry.Manager) Option[redisManager] {
	return option[redisManager](func(m *redisManager) {
		m.tm = tm
	})
}

func RedisWithScope(v string) Option[redisManager] {
	return option[redisManager](func(m *redisManager) {
		m.scope = v
	})
}

func NewRedisOptions() []Option[redisManager] {
	return make([]Option[redisManager], 0)
}
