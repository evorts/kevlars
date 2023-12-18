/**
 * @Author: steven
 * @Description:
 * @File: soap
 * @Date: 29/09/23 10.44
 */

package soap

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"github.com/evorts/kevlars/logger"
	"github.com/evorts/kevlars/telemetry"
	"github.com/evorts/kevlars/utils"
	"github.com/sony/gobreaker"
	"io"
	"net/http"
	"time"
)

type Manager interface {
	ServiceUrl() string

	Post(ctx context.Context, action, url string, body interface{}, bindTo interface{}) (int, error)
	PostRaw(ctx context.Context, action, url string, payload []byte, bindTo interface{}) (int, error)
	PostWithRawResult(ctx context.Context, action, url string, payload []byte, bindTo interface{}) (int, []byte, error)
	PostWithLargeResult(ctx context.Context, action, url string, payload []byte) (int, io.Reader, error)

	WithCircuitBreaker(
		maxRequest uint32, interval time.Duration,
		timeout time.Duration,
	) Manager

	metricIsEnabled() bool
}

type manager struct {
	client                *http.Client
	largeFileClient       *http.Client
	cb                    *gobreaker.CircuitBreaker
	log                   logger.Manager
	metric                telemetry.MetricsManager
	metricEnabled         bool
	circuitBreakerEnabled bool
	serviceUrl            string
	user                  string
	pass                  string
	name                  string
}

func (m *manager) metricIsEnabled() bool {
	return m.metricEnabled && m.metric != nil
}

func (m *manager) ServiceUrl() string {
	return m.serviceUrl
}

type Option interface {
	apply(m *manager)
}

type optionFunc func(*manager)

func (fn optionFunc) apply(m *manager) {
	fn(m)
}

func WithBasicAuth(user, pass string) Option {
	return optionFunc(func(m *manager) {
		m.user = user
		m.pass = pass
	})
}

func WithServiceUrl(url string) Option {
	return optionFunc(func(m *manager) {
		m.serviceUrl = url
	})
}

func WithTransport(t ...TransportOption) Option {
	return optionFunc(func(m *manager) {
		m.client.Transport = NewTransport(t...)
	})
}

func WithLargeContentTimeout(v time.Duration) Option {
	return optionFunc(func(m *manager) {
		m.largeFileClient.Timeout = v
	})
}

func WithLogger(l logger.Manager) Option {
	return optionFunc(func(m *manager) {
		m.log = l
	})
}

func WithMetrics(enabled bool, metric telemetry.MetricsManager) Option {
	return optionFunc(func(m *manager) {
		m.metricEnabled = enabled
		m.metric = metric
	})
}

func WithName(v string) Option {
	return optionFunc(func(m *manager) {
		m.name = v
	})
}

// wrap the request with circuit breaker if needed
func (m *manager) wrap(ctx context.Context, f func(newCtx context.Context) (httpCode int, raw []byte, err error)) (httpCode int, raw []byte, err error) {
	var newCtx = ctx
	if !m.circuitBreakerEnabled {
		code, raws, errF := f(newCtx)
		utils.IfTrueThen(m.metricIsEnabled(), func() {
			tags := []string{
				"response_code:" + utils.IntToString(code),
			}
			utils.IfErrorThen(errF, func() {
				tags = append(tags, "error:"+errF.Error())
			})
			m.metric.StartDefault("soap.call." + m.name).Push(tags...)
		})
		return code, raws, errF
	}
	rs, err := m.cb.Execute(func() (interface{}, error) {
		code, raws, errF := f(newCtx)
		rs := []interface{}{code, raws}
		utils.IfTrueThen(m.metricIsEnabled(), func() {
			tags := []string{
				"response_code:" + utils.IntToString(code),
			}
			utils.IfErrorThen(errF, func() {
				tags = append(tags, "error:"+errF.Error())
			})
			m.metric.StartDefault("soap.call.cb." + m.name).Push(tags...)
		})
		return rs, errF
	})
	if err != nil {
		return http.StatusExpectationFailed, nil, err
	}
	response, ok := rs.([]interface{})
	if ok {
		return response[0].(int), response[1].([]byte), nil
	}
	return http.StatusInternalServerError, nil, errors.New("not a valid response instance")
}

func (m *manager) doRaw(ctx context.Context, action, method, url string, payload []byte, bindTo interface{}) (httpCode int, raw []byte, err error) {
	var (
		req  *http.Request
		resp *http.Response
		buf  bytes.Buffer
	)
	url = utils.IfEmpty(url, m.serviceUrl)
	m.log.InfoWithProps(map[string]interface{}{
		"method": method,
		"url":    url,
		"action": action,
	}, "calling remote soap api")
	req, err = http.NewRequest(method, url, bytes.NewReader(payload))
	if err != nil {
		return http.StatusBadRequest, nil, err
	}
	ctx = context.WithValue(ctx, "soap.action", action)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "text/xml")
	req.Header.Set("SoapAction", action)
	req.SetBasicAuth(m.user, m.pass)
	resp, err = m.client.Do(req)
	if err != nil {
		return http.StatusExpectationFailed, nil, err
	}
	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		return resp.StatusCode, nil, err
	}
	raw = buf.Bytes()
	err = xml.Unmarshal(raw, bindTo)
	if err != nil {
		return resp.StatusCode, raw, err
	}
	return resp.StatusCode, raw, nil
}

func (m *manager) do(ctx context.Context, action, method, url string, props interface{}, bindTo interface{}) (httpCode int, err error) {
	var payload []byte
	payload, err = xml.Marshal(props)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	httpCode, _, err = m.wrap(ctx, func(newCtx context.Context) (httpCode int, raw []byte, err error) {
		return m.doRaw(newCtx, action, method, url, payload, bindTo)
	})
	return
}

func (m *manager) Post(ctx context.Context, action, url string, props interface{}, bindTo interface{}) (int, error) {
	rs, _, err := m.wrap(ctx, func(newCtx context.Context) (httpCode int, raw []byte, err error) {
		code, err := m.do(ctx, action, "POST", url, props, bindTo)
		return code, nil, err
	})
	return rs, err
}

func (m *manager) PostRaw(ctx context.Context, action, url string, payload []byte, bindTo interface{}) (int, error) {
	httpCode, _, err := m.wrap(ctx, func(newCtx context.Context) (httpCode int, raw []byte, err error) {
		return m.doRaw(ctx, action, "POST", url, payload, bindTo)
	})
	return httpCode, err
}

func (m *manager) PostWithRawResult(ctx context.Context, action, url string, payload []byte, bindTo interface{}) (int, []byte, error) {
	return m.wrap(ctx, func(newCtx context.Context) (httpCode int, raw []byte, err error) {
		return m.doRaw(ctx, action, "POST", url, payload, bindTo)
	})
}

func (m *manager) PostWithLargeResult(ctx context.Context, action, url string, payload []byte) (int, io.Reader, error) {
	var (
		req  *http.Request
		resp *http.Response
		//buf  bytes.Buffer
		err error
	)
	url = utils.IfEmpty(url, m.serviceUrl)
	m.log.WarnWithProps(map[string]interface{}{
		"method": "POST",
		"url":    url,
		"action": action,
	}, "calling remote soap api")
	req, err = http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		return http.StatusBadRequest, nil, err
	}
	ctx = context.WithValue(ctx, "soap.action", action)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "text/xml")
	req.Header.Set("SoapAction", action)
	req.SetBasicAuth(m.user, m.pass)
	resp, err = m.largeFileClient.Do(req)
	if err != nil {
		return http.StatusExpectationFailed, nil, err
	}
	// compress reader
	//reader := bzip2.NewReader()
	return resp.StatusCode, resp.Body, nil
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

func New(opts ...Option) Manager {
	m := &manager{
		client:          &http.Client{},
		largeFileClient: &http.Client{},
	}
	for _, opt := range opts {
		opt.apply(m)
	}
	return m
}
