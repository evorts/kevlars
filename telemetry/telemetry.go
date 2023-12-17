/**
 * @Author: steven
 * @Description:
 * @File: telemetry
 * @Date: 18/12/23 00.55
 */

package telemetry

import (
	"context"
	"errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	traceSdk "go.opentelemetry.io/otel/sdk/trace"
	semConv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

type Manager interface {
	Init() error
	MustInit() Manager
	Enabled() bool
	Tracer() trace.Tracer

	// NewTracer returns a Tracer with the given name and options. If a Tracer for
	// the given name and options does not exist it is created, otherwise the
	// existing Tracer is returned.
	//
	// If name is empty, DefaultTracerName is used instead.
	//
	// This method is safe to be called concurrently.
	NewTracer(name string, opts ...trace.TracerOption) trace.Tracer

	// ForceFlush immediately exports all spans that have not yet been exported for
	// all the registered span processors.
	ForceFlush(ctx context.Context) error
	// Shutdown shuts down the span processors in the order they were registered.
	Shutdown(ctx context.Context) error
}

type SpanProcessor string

const (
	SpanProcessorBatch        SpanProcessor = "batch"
	SpanProcessorSimple       SpanProcessor = "simple"
	SpanProcessorSimpleCustom SpanProcessor = "simple_custom"
)

type Exporter interface {
	traceSdk.SpanExporter

	Init() error
	SpanProcessorType() SpanProcessor
}

type manager struct {
	enabled bool

	serviceName    string
	serviceVersion string
	env            string

	exporters []Exporter
	tp        *traceSdk.TracerProvider
	tracer    trace.Tracer
	metric    metric.MeterProvider
}

func (m *manager) Enabled() bool {
	return m.enabled
}

func (m *manager) Tracer() trace.Tracer {
	return m.tracer
}

// NewTracer returns a Tracer with the given name and options. If a Tracer for
// the given name and options does not exist it is created, otherwise the
// existing Tracer is returned.
//
// If name is empty, DefaultTracerName is used instead.
//
// This method is safe to be called concurrently.
func (m *manager) NewTracer(name string, opts ...trace.TracerOption) trace.Tracer {
	return m.tp.Tracer(name, opts...)
}

// ForceFlush immediately exports all spans that have not yet been exported for
// all the registered span processors.
func (m *manager) ForceFlush(ctx context.Context) error {
	return m.tp.ForceFlush(ctx)
}

// Shutdown shuts down the span processors in the order they were registered.
func (m *manager) Shutdown(ctx context.Context) error {
	return m.tp.Shutdown(ctx)
}

func (m *manager) MustInit() Manager {
	if err := m.Init(); err != nil {
		panic(err)
	}
	return m
}

func (m *manager) Init() error {
	if len(m.exporters) < 1 {
		return errors.New("no exporters defined")
	}
	providers := make([]traceSdk.TracerProviderOption, 0)
	for _, exporter := range m.exporters {
		if err := exporter.Init(); err != nil {
			return err
		}
		var processor traceSdk.SpanProcessor
		if exporter.SpanProcessorType() == SpanProcessorSimple {
			processor = traceSdk.NewSimpleSpanProcessor(exporter)
		} else if exporter.SpanProcessorType() == SpanProcessorBatch {
			processor = traceSdk.NewBatchSpanProcessor(exporter)
		}
		if processor == nil {
			continue
		}
		providers = append(providers, traceSdk.WithSpanProcessor(processor))
	}
	if len(providers) < 1 {
		return errors.New("no providers defined")
	}
	providers = append(
		providers,
		// Record information about this application in a Resource.
		traceSdk.WithResource(
			resource.NewWithAttributes(
				semConv.SchemaURL,
				semConv.ServiceNameKey.String(m.serviceName),
				semConv.ServiceVersionKey.String(m.serviceVersion),
				attribute.String("env", m.env),
			),
		),
	)
	m.tp = traceSdk.NewTracerProvider(providers...)
	otel.SetTracerProvider(m.tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{}, propagation.Baggage{},
		),
	)
	m.tracer = m.tp.Tracer(
		m.serviceName,
		trace.WithInstrumentationVersion(m.serviceVersion),
	)
	return nil
}

func NewTelemetryManager(env, serviceName, serviceVersion string, enabled bool, exporters ...Exporter) Manager {
	return &manager{env: env, serviceName: serviceName, serviceVersion: serviceVersion, enabled: enabled, exporters: exporters}
}
