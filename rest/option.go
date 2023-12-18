/**
 * @Author: steven
 * @Description:
 * @File: option
 * @Date: 29/09/23 10.49
 */

package rest

import (
	"github.com/evorts/kevlars/logger"
	"github.com/evorts/kevlars/telemetry"
)

type Option interface {
	apply(m *manager)
}

type optionFunc func(*manager)

func (fn optionFunc) apply(m *manager) {
	fn(m)
}

func WithName(name string) Option {
	return optionFunc(func(m *manager) {
		m.name = name
	})
}

func WithBaseUrl(baseUrl string) Option {
	return optionFunc(func(m *manager) {
		m.baseUrl = baseUrl
	})
}

func WithDebugging(enabled bool) Option {
	return optionFunc(func(m *manager) {
		m.debugEnabled = enabled
		m.client.Debug = enabled
	})
}

func WithTracing(enabled bool) Option {
	return optionFunc(func(m *manager) {
		m.traceEnabled = enabled
		if enabled {
			m.client.EnableTrace()
		}
	})
}

func WithMetrics(enabled bool, metric telemetry.MetricsManager) Option {
	return optionFunc(func(m *manager) {
		m.metricEnabled = enabled
		m.metric = metric
	})
}

func WithLogger(l logger.Manager) Option {
	return optionFunc(func(m *manager) {
		m.log = l
	})
}

func WithTransport(t ...TransportOption) Option {
	return optionFunc(func(m *manager) {
		m.client.SetTransport(NewTransport(t...).BuildRoundTripper())
	})
}

func WithProxy(addr string) Option {
	return optionFunc(func(m *manager) {
		m.client = m.client.SetProxy(addr)
	})
}

func WithLogRequest(v bool) Option {
	return optionFunc(func(m *manager) {
		m.logRequestPayload = v
	})
}

func WithLogResponse(v bool) Option {
	return optionFunc(func(m *manager) {
		m.logResponse = v
	})
}
