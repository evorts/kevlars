/**
 * @Author: steven
 * @Description:
 * @File: metrics
 * @Date: 18/12/23 00.54
 */

package telemetry

import (
	"errors"
	"github.com/evorts/kevlars/logger"
	"time"
)

type MetricPushAction interface {
	Push(tags ...string)
}

type metricActionNoop struct {
	log logger.Manager
}

func (m *metricActionNoop) Push(_ ...string) {
	m.log.Debug("metric action push noop triggered")
}

type metricPushAction struct {
	name    string
	tags    []string
	startAt time.Time
	metrics []Metric
	log     logger.Manager
}

func (m *metricPushAction) appendTags(tags ...string) {
	m.tags = append(m.tags, tags...)
}

func (m *metricPushAction) Push(tags ...string) {
	m.appendTags(tags...)
	for _, metric := range m.metrics {
		m.log.InfoWithProps(map[string]interface{}{
			"name": m.name,
			"tags": tags,
		}, "sending metrics")
		metric.Count(m.name, 1, m.tags)
		elapsedTime := time.Since(m.startAt).Milliseconds()
		metric.Histogram(m.name+".histogram", float64(elapsedTime), m.tags)
	}
}

func (m *metricPushAction) addMetrics(metrics ...Metric) MetricPushAction {
	m.metrics = append(m.metrics, metrics...)
	return m
}

func NewMetricActionNoop(log logger.Manager) MetricPushAction {
	return &metricActionNoop{log: log}
}

func newMetricPushAction(log logger.Manager, name string, metrics []Metric, tags ...string) MetricPushAction {
	return &metricPushAction{
		name:    name,
		tags:    tags,
		startAt: time.Now(),
		metrics: metrics,
		log:     log,
	}
}

type MetricsManager interface {
	StartBy(provider, name string) MetricPushAction
	Start(name string) MetricPushAction
	StartDefault(name string) MetricPushAction
	Enabled() bool

	MetricActions

	MustInit() MetricsManager
}

type metricsManager struct {
	enabled     bool
	env         string
	serviceName string
	metrics     map[string]Metric
	log         logger.Manager
}

var defaultProvider = MetricDatadog

func (m *metricsManager) getProvider(name string) (Metric, error) {
	if p, ok := m.metrics[name]; ok {
		return p, nil
	}
	return nil, errors.New("metric provider not found")
}
func (m *metricsManager) execBy(provider string, run func(metric Metric)) {
	if p, err := m.getProvider(provider); err == nil {
		run(p)
	}
}
func (m *metricsManager) execAll(run func(metric Metric)) {
	for _, mt := range m.metrics {
		run(mt)
	}
}

func (m *metricsManager) Enabled() bool {
	return m.enabled
}

func (m *metricsManager) Count(metricName string, value int64, tags []string) {
	m.execAll(func(metric Metric) {
		metric.Count(metricName, value, tags)
	})
}

func (m *metricsManager) Gauge(metricName string, value float64, tags []string) {
	m.execAll(func(metric Metric) {
		metric.Gauge(metricName, value, tags)
	})
}

func (m *metricsManager) Histogram(metricName string, value float64, tags []string) {
	m.execAll(func(metric Metric) {
		metric.Histogram(metricName, value, tags)
	})
}

func (m *metricsManager) Distribution(metricName string, value float64, tags []string) {
	m.execAll(func(metric Metric) {
		metric.Distribution(metricName, value, tags)
	})
}

// StartBy specific metric provider only
func (m *metricsManager) StartBy(provider, name string) MetricPushAction {
	p, err := m.getProvider(provider)
	if !m.enabled || err != nil {
		return NewMetricActionNoop(m.log)
	}
	actionName := m.serviceName + "." + name
	tags := []string{"src_env:" + m.env, "name:" + actionName}
	return newMetricPushAction(m.log, actionName, []Metric{p}, tags...)
}

// StartDefault using default metric provider
func (m *metricsManager) StartDefault(name string) MetricPushAction {
	p, err := m.getProvider(defaultProvider.String())
	if !m.enabled || err != nil {
		return NewMetricActionNoop(m.log)
	}
	actionName := m.serviceName + "." + name
	tags := []string{"src_env:" + m.env, "name:" + actionName}
	return newMetricPushAction(m.log, actionName, []Metric{p}, tags...)
}

// Start using all instantiated providers
func (m *metricsManager) Start(name string) MetricPushAction {
	if !m.enabled {
		return NewMetricActionNoop(m.log)
	}
	actionName := m.serviceName + "." + name
	tags := []string{"src_env:" + m.env, "name:" + actionName}
	metrics := make([]Metric, 0)
	for _, metric := range m.metrics {
		metrics = append(metrics, metric)
	}
	return newMetricPushAction(m.log, actionName, metrics, tags...)
}

func (m *metricsManager) MustInit() MetricsManager {
	if len(m.metrics) > 0 {
		m.enabled = true
	}
	return m
}

func NewMetricsManager(log logger.Manager, serviceName string, metrics ...Metric) MetricsManager {
	m := make(map[string]Metric)
	for _, metric := range metrics {
		m[metric.Key()] = metric
	}
	return &metricsManager{metrics: m, serviceName: serviceName, log: log}
}
