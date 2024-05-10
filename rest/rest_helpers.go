/**
 * @Author: steven
 * @Description:
 * @File: rest_helpers
 * @Date: 29/09/23 10.50
 */

package rest

import (
	"context"
	"errors"
	"fmt"
	"github.com/evorts/kevlars/requests"
	"github.com/evorts/kevlars/rules"
	"github.com/evorts/kevlars/utils"
	"github.com/go-resty/resty/v2"
	otelAttr "go.opentelemetry.io/otel/attribute"
	semConv "go.opentelemetry.io/otel/semconv/v1.12.0"
	otelTrace "go.opentelemetry.io/otel/trace"
	"net/http"
	"strings"
)

func (m *manager) buildUrl(path string) string {
	format := "%s/%s"
	if strings.HasPrefix(path, "/") {
		format = "%s%s"
	}
	return fmt.Sprintf(format, m.baseUrl, path)
}

// parseUrl parse value given, when full url then there's no need to build with base url
//
//goland:noinspection HttpUrlsUsage
func (m *manager) parseUrl(v string) string {
	if strings.HasPrefix(v, "http://") || strings.HasPrefix(v, "https://") {
		return v
	}
	return m.buildUrl(v)
}

func (m *manager) buildRequest(ctx context.Context, params map[string]string, body interface{}) *resty.Request {
	return m.buildRequestWithOpts(ctx, params, body, nil, nil)
}

func (m *manager) buildRequestWithOpts(ctx context.Context, params map[string]string, body, bindSuccessTo, bindFailedTo interface{}, opts ...Props) *resty.Request {
	p := newProps()
	req := m.client.R()
	req.SetContext(ctx)
	if len(opts) > 0 {
		// populate options
		for _, opt := range opts {
			opt.apply(p)
		}
	}
	if m.retry != nil {
		// @todo: support automatic retry with step back interval
	}
	if params != nil {
		req.SetQueryParams(params)
	}
	if body != nil {
		req.SetBody(body)
	}
	if bindSuccessTo != nil {
		req.SetResult(bindSuccessTo)
	}
	if bindFailedTo != nil {
		req.SetError(bindFailedTo)
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

// wrapWithContext the request with circuit breaker if needed
func (m *manager) wrapWithContext(
	ctx context.Context, path string,
	f func(newCtx context.Context) (*resty.Response, error),
) (context.Context, *resty.Response, error) {
	var newCtx = ctx
	if m.tm != nil {
		opts := make([]otelTrace.SpanStartOption, 0)
		opts = append(
			opts,
			otelTrace.WithSpanKind(otelTrace.SpanKindInternal),
			otelTrace.WithAttributes(
				semConv.PeerServiceKey.String(m.name),
				otelAttr.String("remote.path", path),
			),
		)
		var span otelTrace.Span
		newCtx, span = m.tm.Tracer().Start(ctx, "rest.call", opts...)
		defer span.End()
	}
	m.log.InfoWithProps(map[string]interface{}{
		"path":   path,
		"base":   m.parseUrl(path),
		"req.id": requests.Id(ctx),
	}, "calling remote rest api")
	if !m.circuitBreakerEnabled {
		resp, errF := f(newCtx)
		rules.WhenTrue(m.metricIsEnabled(), func() {
			tags := []string{
				"response_code:" + utils.IntToString(resp.StatusCode()),
				"path:" + path,
				"base:" + m.parseUrl(path),
			}
			rules.WhenError(errF, func() {
				tags = append(tags, "error:"+errF.Error())
			})
			m.metric.StartDefault("rest.response." + m.name).Push(tags...)
		})
		return newCtx, resp, errF
	}
	rs, err := m.cb.Execute(func() (interface{}, error) {
		resp, errF := f(newCtx)
		rules.WhenTrue(m.metricIsEnabled(), func() {
			tags := []string{
				"response_code:" + utils.IntToString(resp.StatusCode()),
				"path:" + path,
				"base:" + m.parseUrl(path),
			}
			rules.WhenError(errF, func() {
				tags = append(tags, "error:"+errF.Error())
			})
			m.metric.StartDefault("rest.response.cb." + m.name).Push(tags...)
		})
		return resp, errF
	})
	var (
		response *resty.Response
		ok       bool
	)
	if rs != nil {
		response, ok = rs.(*resty.Response)
	}
	if err != nil {
		return newCtx, response, err
	}
	if ok {
		return newCtx, response, nil
	}
	return newCtx, nil, errors.New("not a valid response instance")
}

// wrap the request with circuit breaker if needed
func (m *manager) wrap(ctx context.Context, path string, f func(newCtx context.Context) (*resty.Response, error)) (*resty.Response, error) {
	var newCtx = ctx
	if m.tm != nil {
		opts := make([]otelTrace.SpanStartOption, 0)
		opts = append(
			opts,
			otelTrace.WithSpanKind(otelTrace.SpanKindInternal),
			otelTrace.WithAttributes(
				semConv.PeerServiceKey.String(m.name),
				otelAttr.String("remote.path", path),
			),
		)
		var span otelTrace.Span
		newCtx, span = m.tm.Tracer().Start(ctx, "rest.call", opts...)
		defer span.End()
	}
	m.log.InfoWithProps(map[string]interface{}{
		"path":   path,
		"base":   m.parseUrl(path),
		"req.id": requests.Id(ctx),
	}, "calling remote rest api")
	if !m.circuitBreakerEnabled {
		resp, errF := f(newCtx)
		rules.WhenTrue(m.metricIsEnabled(), func() {
			tags := []string{
				"response_code:" + utils.IntToString(resp.StatusCode()),
			}
			rules.WhenError(errF, func() {
				tags = append(tags, "error:"+errF.Error())
			})
			m.metric.StartDefault("rest.response." + m.name).Push(tags...)
		})
		return resp, errF
	}
	rs, err := m.cb.Execute(func() (interface{}, error) {
		resp, errF := f(newCtx)
		rules.WhenTrue(m.metricIsEnabled(), func() {
			tags := []string{
				"response_code:" + utils.IntToString(resp.StatusCode()),
			}
			rules.WhenError(errF, func() {
				tags = append(tags, "error:"+errF.Error())
			})
			m.metric.StartDefault("rest.response.cb." + m.name).Push(tags...)
		})
		return resp, errF
	})
	var (
		response *resty.Response
		ok       bool
	)
	if rs != nil {
		response, ok = rs.(*resty.Response)
	}
	if err != nil {
		return response, err
	}
	if ok {
		return response, nil
	}
	return nil, errors.New("not a valid response instance")
}

func (m *manager) parseReturnWithContext(ctx context.Context, rs *resty.Response, e error) (httpCode int, err error) {
	clientId := requests.ClientId(ctx)
	reqId := requests.Id(ctx)
	reqAccept := rules.WhenTrueR1(rs != nil && rs.Request != nil, func() string {
		return rs.Request.Header.Get("Accept")
	}, func() string { return "" })
	reqContentType := rules.WhenTrueR1(rs != nil && rs.Request != nil, func() string {
		return rs.Request.Header.Get("Content-Type")
	}, func() string { return "" })
	respAccept := rules.WhenTrueR1(rs != nil, func() string { return rs.Header().Get("Accept") }, func() string { return "" })
	respContentType := rules.WhenTrueR1(rs != nil, func() string { return rs.Header().Get("Content-Type") }, func() string { return "" })
	traceMap := map[string]interface{}{
		"http_code":         httpCode,
		"client_id":         clientId,
		"req_id":            reqId,
		"req_accept":        reqAccept,
		"req_content_type":  reqContentType,
		"resp_accept":       respAccept,
		"resp_content_type": respContentType,
	}
	m.log.WhenErrorWithProps(e, traceMap)
	m.log.InfoWithProps(traceMap, "parse return with context")
	method := rules.WhenTrueR1(rs != nil && rs.Request != nil, func() string {
		return rs.Request.Method
	}, func() string {
		return "unknown"
	})
	url := rules.WhenTrueR1(rs != nil && rs.Request != nil, func() string {
		return rs.Request.URL
	}, func() string {
		return ""
	})
	payload := rules.WhenTrueR1(rs != nil && rs.Request != nil, func() interface{} {
		return rs.Request.Body
	}, func() interface{} {
		return nil
	})
	errString := rules.WhenTrueR1(e != nil, func() string {
		return e.Error()
	}, func() string {
		return ""
	})
	rules.WhenTrue(m.logRequestPayload, func() {
		m.log.InfoWithProps(map[string]interface{}{
			"ctx":       "rest.request.payload",
			"http_code": httpCode,
			"client_id": clientId,
			"req_id":    reqId,
			"path":      url,
			"method":    method,
			"payload":   payload,
		}, errString)
	})
	rules.WhenTrue(m.logResponse, func() {
		m.log.InfoWithProps(map[string]interface{}{
			"ctx":       "rest.response",
			"http_code": httpCode,
			"client_id": clientId,
			"req_id":    reqId,
			"path":      rs.Request.URL,
			"method":    method,
			"body": rules.WhenTrueR1(rs != nil, func() string {
				return rs.String()
			}, func() string {
				return ""
			}),
		}, errString)
	})
	if e != nil {
		return http.StatusExpectationFailed, e
	}
	return rs.StatusCode(), nil
}

func (m *manager) parseReturn(rs *resty.Response, e error) (httpCode int, err error) {
	return m.parseReturnWithContext(rs.Request.Context(), rs, e)
}
