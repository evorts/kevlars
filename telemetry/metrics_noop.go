/**
 * @Author: steven
 * @Description:
 * @File: metrics_noop
 * @Date: 07/01/24 13.00
 */

package telemetry

import "github.com/evorts/kevlars/logger"

type metricNoop struct {
	action MetricPushAction
}

func (m *metricNoop) StartBy(provider, name string) MetricPushAction {
	return m.action
}

func (m *metricNoop) Start(name string) MetricPushAction {
	return m.action
}

func (m *metricNoop) StartDefault(name string) MetricPushAction {
	return m.action
}

func (m *metricNoop) Enabled() bool {
	return false
}

func (m *metricNoop) Count(metricName string, value int64, tags []string) {
	// do nothing
}

func (m *metricNoop) Gauge(metricName string, value float64, tags []string) {
	// do nothing
}

func (m *metricNoop) Histogram(metricName string, value float64, tags []string) {
	// do nothing
}

func (m *metricNoop) Distribution(metricName string, value float64, tags []string) {
	// do nothing
}

func (m *metricNoop) MustInit() MetricsManager {
	return m
}

func NewMetricNoop() MetricsManager {
	return &metricNoop{action: NewMetricActionNoop(logger.NewNoop())}
}
