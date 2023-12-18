package midware

import (
	"errors"
	"github.com/evorts/kevlars/contracts"
	"github.com/evorts/kevlars/logger"
	"github.com/evorts/kevlars/telemetry"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/codes"
	otelTrace "go.opentelemetry.io/otel/trace"
	"net/http"
	"time"
)

type EchoHttpError interface {
	Middleware() echo.MiddlewareFunc
	Handler(err error, c echo.Context)
}

type echoHttpError struct {
	l  logger.Manager
	tm telemetry.Manager
}

func (eh *echoHttpError) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			timeStarted := time.Now()
			err := next(c)
			status := c.Response().Status
			httpErr := new(echo.HTTPError)
			if errors.As(err, &httpErr) {
				status = httpErr.Code
			}
			fields := map[string]interface{}{
				"latency": int64(time.Since(timeStarted) / time.Millisecond),
				"method":  c.Request().Method,
				"path":    c.Request().URL.Path,
				"query":   c.Request().URL.RawQuery,
				"status":  status,
			}
			eh.l.ErrorWithPropsWhen(err != nil, fields, func(messages func(...interface{})) {
				messages(err.Error())
			})
			return err
		}
	}
}

func (eh *echoHttpError) Handler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}
	req := c.Request()
	ctx := req.Context()
	var he *echo.HTTPError
	ok := errors.As(err, &he)
	if ok {
		if he.Internal != nil {
			var herr *echo.HTTPError
			if errors.As(he.Internal, &herr) {
				he = herr
			}
		}
	} else {
		he = &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
		}
	}

	// Issue #1426
	code := he.Code
	errorMessage := ""
	if m, okm := he.Message.(string); okm {
		errorMessage = m + " " + err.Error()
	}

	eh.l.ErrorWithProps(map[string]interface{}{"code": code}, "http error:", he)

	newCtx, span := eh.tm.Tracer().Start(ctx, "http.error", otelTrace.WithSpanKind(otelTrace.SpanKindServer))
	// Record the error.
	span.RecordError(he)
	// Also mark span as failed.
	span.SetStatus(codes.Error, he.Error())
	defer span.End()

	// Send response
	if c.Request().Method == http.MethodHead { // Issue #608
		err = c.NoContent(he.Code)
	} else {
		httpCode, rs := contracts.NewResponseFail(code, errorMessage, contracts.ErrorDetail{})
		err = c.JSON(httpCode, rs)
	}
	req = req.WithContext(newCtx)
	c.SetRequest(req)
}

func NewEchoHttpError(tm telemetry.Manager, l logger.Manager) EchoHttpError {
	return &echoHttpError{tm: tm, l: l}
}
