/**
 * @Author: steven
 * @Description:
 * @File: exporters
 * @Date: 18/12/23 00.49
 */

package telemetry

type ExporterProvider string

const (
	ExporterStandard ExporterProvider = "standard"
	ExporterZipkin   ExporterProvider = "zipkin"
	ExporterDatadog  ExporterProvider = "datadog"
)

func (t ExporterProvider) String() string {
	return string(t)
}
