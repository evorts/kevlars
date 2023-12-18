/**
 * @Author: steven
 * @Description:
 * @File: telemetry_noop
 * @Date: 24/12/23 06.04
 */

package telemetry

import (
	"context"
	"go.opentelemetry.io/otel/trace"
)

type managerNoop struct {
	tracer trace.Tracer
}

func (m *managerNoop) Init() error {
	return nil
}

func (m *managerNoop) MustInit() Manager {
	return m
}

func (m *managerNoop) Enabled() bool {
	return true
}

func (m *managerNoop) Tracer() trace.Tracer {
	return m.tracer
}

func (m *managerNoop) NewTracer(name string, opts ...trace.TracerOption) trace.Tracer {
	return NewNoopTracer(name)
}

func (m *managerNoop) ForceFlush(ctx context.Context) error {
	return nil
}

func (m *managerNoop) Shutdown(ctx context.Context) error {
	return nil
}

func NewNoop() Manager {
	return &managerNoop{
		tracer: NewNoopTracer("noop"),
	}
}
