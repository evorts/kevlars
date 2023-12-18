/**
 * @Author: steven
 * @Description:
 * @File: options
 * @Date: 29/09/23 10.49
 */

package rpc

import (
	"github.com/evorts/kevlars/logger"
	"github.com/evorts/kevlars/telemetry"
	"google.golang.org/grpc/credentials"
	"time"
)

type Option interface {
	apply(m *grpcClientManager)
}

type option func(m *grpcClientManager)

func (o option) apply(m *grpcClientManager) {
	o(m)
}

func WithName(v string) Option {
	return option(func(m *grpcClientManager) {
		m.name = v
	})
}

func WithLogger(v logger.Manager) Option {
	return option(func(m *grpcClientManager) {
		m.log = v
	})
}

func WithTelemetryEnabled(v bool) Option {
	return option(func(m *grpcClientManager) {
		m.telemetryEnabled = v
	})
}

func WithTelemetry(v telemetry.Manager) Option {
	return option(func(m *grpcClientManager) {
		m.tm = v
	})
}

func WithAuthorizationToken(v string) Option {
	return option(func(m *grpcClientManager) {
		m.token = v
	})
}

func WithTransportCreds(v credentials.TransportCredentials) Option {
	return option(func(m *grpcClientManager) {
		m.creds = v
	})
}

func WithTimeout(v time.Duration) Option {
	return option(func(m *grpcClientManager) {
		m.timeout = &v
	})
}

func WithMetrics(enabled bool, metric telemetry.MetricsManager) Option {
	return option(func(m *grpcClientManager) {
		m.metricEnabled = enabled
		m.metric = metric
	})
}

func WithLoggingRequestPayload(v bool, inJson bool) Option {
	return option(func(m *grpcClientManager) {
		m.logRequestPayload = v
		m.logRequestPayloadInJson = inJson
	})
}
