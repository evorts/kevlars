package rest

import "go.opentelemetry.io/otel/metric"

type httpMetric struct {
	attemptsCounter         metric.Int64Counter
	noRequestCounter        metric.Int64Counter
	errorsCounter           metric.Int64Counter
	successesCounter        metric.Int64Counter
	failureCounter          metric.Int64Counter
	redirectCounter         metric.Int64Counter
	timeoutsCounter         metric.Int64Counter
	canceledCounter         metric.Int64Counter
	deadlineExceededCounter metric.Int64Counter
	totalDurationCounter    metric.Int64Histogram
	inFlightCounter         metric.Int64UpDownCounter
}
