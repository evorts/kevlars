/**
 * @Author: steven
 * @Description:
 * @File: metric_datadog
 * @Date: 18/12/23 00.55
 */

package telemetry

import (
	"errors"
	"fmt"
	"github.com/DataDog/datadog-go/v5/statsd"
	"github.com/evorts/kevlars/logger"
)

type metricDataDog struct {
	key          string // used for metrics map
	name         string // used to send to datadog
	env          string
	sampleRate   float64
	agentAddress string

	client *statsd.Client
	log    logger.Manager
}

func MetricDatadogWithEnv(env string) MetricOption[*metricDataDog] {
	return func(m *metricDataDog) {
		m.env = env
	}
}

func MetricDatadogWithSampleRate(sampleRate float64) MetricOption[*metricDataDog] {
	return func(m *metricDataDog) {
		m.sampleRate = sampleRate
	}
}

func MetricDatadogWithAgentAddress(address string) MetricOption[*metricDataDog] {
	return func(m *metricDataDog) {
		m.agentAddress = address
	}
}

func (m *metricDataDog) IncrSuccess(metricName string) {
	m.Incr(metricName, []string{MetricResponseCodeSuccess.String()})
}

func (m *metricDataDog) IncrFail(metricName string, _ error) {
	m.Incr(metricName, []string{MetricResponseCodeFail.String()})
}

func (m *metricDataDog) DecrSuccess(metricName string) {
	m.Decr(metricName, []string{MetricResponseCodeSuccess.String()})
}

func (m *metricDataDog) DecrFail(metricName string, _ error) {
	m.Decr(metricName, []string{MetricResponseCodeFail.String()})
}

func (m *metricDataDog) IncrHTTP(method, metricName string, httpStatusCode int) {
	// MetricName Syntax:{HTTPMethod}_{URL}
	mn := fmt.Sprintf("%s_%s", method, metricName)
	rsCode := fmt.Sprintf("%s:%d", MetricResponseCodeBase, httpStatusCode)
	m.Incr(mn, []string{rsCode})
}

func (m *metricDataDog) DecrHTTP(method, metricName string, httpStatusCode int) {
	// MetricName Syntax:{HTTPMethod}_{URL}
	mn := fmt.Sprintf("%s_%s", method, metricName)
	rsCode := fmt.Sprintf("%s:%d", MetricResponseCodeBase, httpStatusCode)
	m.Decr(mn, []string{rsCode})
}

func (m *metricDataDog) Incr(metricName string, tags []string) {
	err := m.client.Incr(metricName, tags, m.sampleRate)
	m.log.DebugWhen(err != nil, func(debug func(messages ...interface{})) {
		debug("datadog incr err:", err.Error())
	})
}

func (m *metricDataDog) Decr(metricName string, tags []string) {
	err := m.client.Decr(metricName, tags, m.sampleRate)
	m.log.DebugWhen(err != nil, func(messages func(...interface{})) {
		messages("datadog decr err:", err.Error())
	})
}

func (m *metricDataDog) Count(metricName string, value int64, tags []string) {
	err := m.client.Count(metricName, value, tags, m.sampleRate)
	m.log.DebugWhen(err != nil, func(messages func(...interface{})) {
		messages("datadog count err:", err.Error())
	})
}

func (m *metricDataDog) Gauge(metricName string, value float64, tags []string) {
	err := m.client.Gauge(metricName, value, tags, m.sampleRate)
	m.log.DebugWhen(err != nil, func(messages func(...interface{})) {
		messages("datadog gauge err:", err.Error())
	})
}

func (m *metricDataDog) Histogram(metricName string, value float64, tags []string) {
	err := m.client.Histogram(metricName, value, tags, m.sampleRate)
	m.log.DebugWhen(err != nil, func(messages func(...interface{})) {
		messages("datadog histogram err:", err.Error())
	})
}

func (m *metricDataDog) Distribution(metricName string, value float64, tags []string) {
	err := m.client.Distribution(metricName, value, tags, m.sampleRate)
	m.log.DebugWhen(err != nil, func(messages func(...interface{})) {
		messages("datadog distribution err:", err.Error())
	})
}

func (m *metricDataDog) init() error {
	if len(m.agentAddress) < 1 {
		return errors.New("agent address not defined")
	}
	statsd, err := statsd.New(m.agentAddress)
	if err != nil {
		return fmt.Errorf("failed to initialize Datadog Metric %w", err)
	}
	m.client = statsd
	return nil
}

func (m *metricDataDog) Key() string {
	return m.key
}

func (m *metricDataDog) Name() string {
	return m.name
}

func (m *metricDataDog) MustInit() Metric {
	if err := m.init(); err != nil {
		panic(err)
	}
	return m
}

func NewDatadogMetrics(key, name string, log logger.Manager, opts ...MetricOption[*metricDataDog]) Metric {
	m := &metricDataDog{key: key, name: name, log: log}
	for _, opt := range opts {
		opt(m)
	}
	return m
}
