/**
 * @Author: steven
 * @Description:
 * @File: rest
 * @Date: 29/09/23 10.49
 */

package rest

import (
	"context"
	"github.com/evorts/kevlars/logger"
	"github.com/evorts/kevlars/telemetry"
	"github.com/go-resty/resty/v2"
	"github.com/sony/gobreaker"
	"net/http"
)

type Manager interface {
	// Get use bindSuccessTo and bindFailedTo to bind remote response to struct, while props to replace manager config on the fly
	Get(ctx context.Context, path string, params map[string]string, bindSuccessTo, bindFailedTo interface{}, opts ...Props) (httpCode int, err error)
	// Post use bindSuccessTo and bindFailedTo to bind remote response to struct, while props to replace manager config on the fly
	Post(ctx context.Context, path string, body interface{}, bindSuccessTo, bindFailedTo interface{}, opts ...Props) (httpCode int, err error)
	// Put use bindSuccessTo and bindFailedTo to bind remote response to struct, while props to replace manager config on the fly
	Put(ctx context.Context, path string, body interface{}, bindSuccessTo, bindFailedTo interface{}, opts ...Props) (httpCode int, err error)
	// Delete use bindSuccessTo and bindFailedTo to bind remote response to struct, while props to replace manager config on the fly
	Delete(ctx context.Context, path string, params map[string]string, bindSuccessTo, bindFailedTo interface{}, opts ...Props) (httpCode int, err error)
	// Head use bindSuccessTo and bindFailedTo to bind remote response to struct, while props to replace manager config on the fly
	Head(ctx context.Context, path string, bindSuccessTo, bindFailedTo interface{}, opts ...Props) (httpCode int, err error)
	// Patch use bindSuccessTo and bindFailedTo to bind remote response to struct, while props to replace manager config on the fly
	Patch(ctx context.Context, path string, body interface{}, bindSuccessTo, bindFailedTo interface{}, opts ...Props) (httpCode int, err error)

	CustomRequest(ctx context.Context, qp map[string]string, body interface{}, opts ...Props) *resty.Request
	GetClient() *resty.Client

	metricIsEnabled() bool

	ManagerConfig
}

type retry struct {
}

type manager struct {
	client    *resty.Client
	transport http.RoundTripper
	cb        *gobreaker.CircuitBreaker
	tm        telemetry.Manager
	metric    telemetry.MetricsManager
	log       logger.Manager
	retry     *retry

	name    string
	baseUrl string

	traceEnabled          bool
	debugEnabled          bool
	metricEnabled         bool
	circuitBreakerEnabled bool

	token string

	logRequestPayload bool
	logResponse       bool

	defaultHeaders map[string]string // header that will always be included in request
}

func (m *manager) CustomRequest(ctx context.Context, qp map[string]string, body interface{}, opts ...Props) *resty.Request {
	p := newProps()
	req := m.client.R()
	req.SetContext(ctx)
	if len(opts) > 0 {
		// populate options
		for _, opt := range opts {
			opt.apply(p)
		}
	}
	if qp != nil {
		req.SetQueryParams(qp)
	}
	if body != nil {
		req.SetBody(body)
	}
	// when initial default headers not set, then use the default declared value
	if m.defaultHeaders != nil && len(m.defaultHeaders) > 0 {
		for k, v := range m.defaultHeaders {
			req.Header.Set(k, v)
		}
	}
	if p.headers != nil && len(p.headers) > 0 {
		for k, v := range p.headers {
			req.Header.Set(k, v)
		}
	}
	// inject token if exist
	if len(p.token) > 0 {
		req.SetAuthToken(p.token)
	}
	return req
}

func (m *manager) GetClient() *resty.Client {
	return m.client
}

func (m *manager) metricIsEnabled() bool {
	return m.metricEnabled && m.metric != nil
}

func (m *manager) Get(ctx context.Context, path string, params map[string]string, bindSuccessTo, bindFailedTo interface{}, opts ...Props) (httpCode int, err error) {
	return m.parseReturnWithContext(m.wrapWithContext(ctx, path, func(newCtx context.Context) (*resty.Response, error) {
		return m.buildRequestWithOpts(newCtx, params, nil, bindSuccessTo, bindFailedTo, opts...).
			SetResult(bindSuccessTo).
			SetError(bindFailedTo).
			Get(m.parseUrl(path))
	}))
}

func (m *manager) Post(ctx context.Context, path string, body interface{}, bindSuccessTo, bindFailedTo interface{}, opts ...Props) (httpCode int, err error) {
	httpCode, err = m.parseReturnWithContext(m.wrapWithContext(ctx, path, func(newCtx context.Context) (*resty.Response, error) {
		return m.buildRequestWithOpts(newCtx, nil, body, bindSuccessTo, bindFailedTo, opts...).
			SetResult(bindSuccessTo).
			SetError(bindFailedTo).
			Post(m.parseUrl(path))
	}))
	return
}

func (m *manager) Put(ctx context.Context, path string, body interface{}, bindSuccessTo, bindFailedTo interface{}, opts ...Props) (httpCode int, err error) {
	return m.parseReturnWithContext(m.wrapWithContext(ctx, path, func(newCtx context.Context) (*resty.Response, error) {
		return m.buildRequestWithOpts(newCtx, nil, body, bindSuccessTo, bindFailedTo, opts...).
			Put(m.parseUrl(path))
	}))
}

func (m *manager) Delete(ctx context.Context, path string, params map[string]string, bindSuccessTo, bindFailedTo interface{}, opts ...Props) (httpCode int, err error) {
	return m.parseReturnWithContext(m.wrapWithContext(ctx, path, func(newCtx context.Context) (*resty.Response, error) {
		return m.buildRequestWithOpts(newCtx, params, nil, bindSuccessTo, bindFailedTo, opts...).
			Delete(m.parseUrl(path))
	}))
}

func (m *manager) Head(ctx context.Context, path string, bindSuccessTo, bindFailedTo interface{}, opts ...Props) (httpCode int, err error) {
	return m.parseReturnWithContext(m.wrapWithContext(ctx, path, func(newCtx context.Context) (*resty.Response, error) {
		return m.buildRequestWithOpts(newCtx, nil, nil, bindSuccessTo, bindFailedTo, opts...).
			Head(m.parseUrl(path))
	}))
}

func (m *manager) Patch(ctx context.Context, path string, body, bindSuccessTo, bindFailedTo interface{}, opts ...Props) (httpCode int, err error) {
	return m.parseReturnWithContext(m.wrapWithContext(ctx, path, func(newCtx context.Context) (*resty.Response, error) {
		return m.buildRequestWithOpts(newCtx, nil, body, bindSuccessTo, bindFailedTo, opts...).
			Patch(m.parseUrl(path))
	}))
}

func New(opts ...Option) Manager {
	m := &manager{
		client:         resty.New(),
		defaultHeaders: map[string]string{},
	}
	for _, opt := range opts {
		opt.apply(m)
	}
	return m
}
