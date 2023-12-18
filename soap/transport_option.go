/**
 * @Author: steven
 * @Description:
 * @File: transport_option
 * @Date: 29/09/23 10.45
 */

package soap

import "github.com/evorts/kevlars/telemetry"

type TransportOption interface {
	apply(tc *transport)
}

type transportOptionFunc func(t *transport)

func (fn transportOptionFunc) apply(t *transport) {
	fn(t)
}

func WithTelemetry(tm telemetry.Manager) TransportOption {
	return transportOptionFunc(func(t *transport) {
		t.tm = tm
	})
}

func WithTransportName(name string) TransportOption {
	return transportOptionFunc(func(t *transport) {
		t.name = name
	})
}
