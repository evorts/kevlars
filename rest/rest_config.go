/**
 * @Author: steven
 * @Description:
 * @File: rest_cfg
 * @Version: 1.0.0
 * @Date: 07/06/23 20.39
 */

package rest

import (
	"github.com/evorts/kevlars/telemetry"
	"github.com/go-resty/resty/v2"
	"github.com/sony/gobreaker"
	"net/http"
	"time"
)

type ManagerConfig interface {
	SetName(name string) Manager
	BaseUrl() string
	ChangeBaseUrl(url string) Manager
	Transport() http.RoundTripper

	// SetDefaultHeaders will set headers that will exist on every request
	SetDefaultHeaders(values map[string]string) Manager
	// WithDefaultHeaders add additional headers on top of default headers
	WithDefaultHeaders(additionalHeaders map[string]string) Manager
	WithTransport(transport http.RoundTripper) Manager
	WithTelemetry(tm telemetry.Manager) Manager
	WithCircuitBreaker(
		maxRequest uint32, interval time.Duration,
		timeout time.Duration,
	) Manager
	// WithToken will set authentication token for every request (Authorization: Bearer)
	WithToken(token string) Manager
	DisableCircuitBreaker() Manager
}

var defaultHeaders = map[string]string{
	"Accept":       "application/json",
	"Content-Type": "application/json",
}

func (m *manager) BaseUrl() string {
	return m.baseUrl
}

func (m *manager) SetName(name string) Manager {
	m.name = name
	return m
}

func (m *manager) DisableCircuitBreaker() Manager {
	m.circuitBreakerEnabled = false
	return m
}

func (m *manager) ChangeBaseUrl(url string) Manager {
	m.baseUrl = url
	return m
}

func (m *manager) Transport() http.RoundTripper {
	return m.transport
}

func (m *manager) WithToken(token string) Manager {
	m.token = token
	return m
}

func (m *manager) useTokenIfExist(r *resty.Request) Manager {
	if len(m.token) > 0 {
		r.SetAuthToken(m.token)
	}
	return m
}

func (m *manager) clearHeaders() Manager {
	// Loop over header names
	for name := range m.client.Header {
		m.client.Header.Del(name)
	}
	return m
}

func (m *manager) SetDefaultHeaders(values map[string]string) Manager {
	m.defaultHeaders = values
	return m
}

func (m *manager) WithDefaultHeaders(additionalHeaders map[string]string) Manager {
	if len(m.defaultHeaders) < 1 {
		m.defaultHeaders = defaultHeaders
	}
	if additionalHeaders != nil && len(additionalHeaders) > 0 {
		for k, v := range additionalHeaders {
			m.defaultHeaders[k] = v
		}
	}
	return m
}

func (m *manager) WithTransport(transport http.RoundTripper) Manager {
	m.transport = transport
	m.client.SetTransport(transport)
	return m
}

func (m *manager) WithTelemetry(tm telemetry.Manager) Manager {
	m.tm = tm
	return m
}

func (m *manager) WithCircuitBreaker(
	maxRequest uint32, interval time.Duration,
	timeout time.Duration,
) Manager {
	m.cb = gobreaker.NewCircuitBreaker(gobreaker.Settings{
		MaxRequests: maxRequest, //maximum number of requests allowed to pass through
		Interval:    interval,   //cyclic period of the closed state
		Timeout:     timeout,    //period of the open state
	})
	m.circuitBreakerEnabled = true
	return m
}
