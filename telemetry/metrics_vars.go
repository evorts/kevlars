/**
 * @Author: steven
 * @Description:
 * @File: metrics_vars
 * @Date: 18/12/23 00.55
 */

package telemetry

type MetricResponseCode string

const (
	MetricResponseCodeBase    MetricResponseCode = "response_code"
	MetricResponseCodeSuccess MetricResponseCode = MetricResponseCodeBase + ":200"
	MetricResponseCodeFail    MetricResponseCode = MetricResponseCodeBase + ":500"
)

func (v MetricResponseCode) String() string {
	return string(v)
}

type MetricProvider string

const (
	MetricDatadog MetricProvider = "datadog"
	MetricOTLP    MetricProvider = "otlp"
)

func (t MetricProvider) String() string {
	return string(t)
}
