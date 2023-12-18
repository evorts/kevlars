package midware

import (
	"fmt"
	"github.com/evorts/kevlars/telemetry"
	"github.com/evorts/kevlars/ts"
	"github.com/evorts/kevlars/utils"
	"github.com/labstack/echo/v4"
	"strings"
	"time"
)

type echoMetricConfig struct {
	name   string
	env    string
	tm     telemetry.Manager
	metric telemetry.MetricsManager
}

// EchoMetricOption represents an option that can be passed to Middleware.
type EchoMetricOption func(*echoMetricConfig)

func EchoMetricWithTelemetry(oTel telemetry.Manager) EchoMetricOption {
	return func(config *echoMetricConfig) {
		config.tm = oTel
	}
}
func EchoMetricWithMetrics(metric telemetry.MetricsManager) EchoMetricOption {
	return func(config *echoMetricConfig) {
		config.metric = metric
	}
}

func EchoMetricWithEnv(env string) EchoMetricOption {
	return func(config *echoMetricConfig) {
		config.env = env
	}
}

func EchoMetricWithName(name string) EchoMetricOption {
	return func(config *echoMetricConfig) {
		config.name = name
	}
}

func EchoMetrics(opts ...EchoMetricOption) echo.MiddlewareFunc {
	mc := new(echoMetricConfig)
	for _, fn := range opts {
		fn(mc)
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			path := req.URL.Path

			defer func(start time.Time) {
				elapsedTime := time.Since(start).Milliseconds()
				// send metrics to datadog
				method := strings.ToLower(req.Method)
				name := fmt.Sprintf("%s", mc.name)
				endpoint := fmt.Sprintf("%s:%s", method, path)
				// send metrics to datadog
				responseStatus := c.Response().Status
				tags := []string{
					"http_endpoint:" + endpoint,
					"http_method:" + method,
					"src_env:" + mc.env,
					"response_code:" + utils.CastToStringND(responseStatus),
				}
				mc.metric.Count(name, 1, tags)
				mc.metric.Histogram(name+".histogram", float64(elapsedTime), tags)
				mc.metric.Distribution(name+".distribution", float64(elapsedTime), tags)

				// @todo: send metric to open telemetry collector
				// here>
			}(ts.Now())

			return next(c)
		}
	}
}
