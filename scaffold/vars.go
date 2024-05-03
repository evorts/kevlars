/**
 * @Author: steven
 * @Description:
 * @File: vars
 * @Date: 27/03/24 10.46
 */

package scaffold

type ConfigKey string

const (
	AppSection ConfigKey = "app"

	AppName            = AppSection + ".name"
	AppVersion         = AppSection + ".version"
	AppEnv             = AppSection + ".env"
	AppPortRest        = AppSection + ".port.rest"
	AppPortGrpc        = AppSection + ".port.grpc"
	AppGracefulTimeout = AppSection + ".graceful_timeout"

	AppLogLevel             = AppSection + ".log.level"
	AppLogTimezone          = AppSection + ".log.tz"
	AppLogUseCustomTimezone = AppSection + ".log.use_custom_timezone"

	AppHealth       = AppSection + ".healthcheck.health"
	AppMetric       = AppSection + ".healthcheck.metrics"
	AppDependencies = AppSection + ".healthcheck.dependencies"
)

func (c ConfigKey) String() string {
	return string(c)
}

type EnvKey string

const (
	EAppEnv EnvKey = "APP_ENV"
)

func (e EnvKey) String() string {
	return string(e)
}

type Key string

const (
	DbKey     Key = "db"
	CacheKey  Key = "cache"
	StatusKey Key = "status"
)

func (k Key) String() string {
	return string(k)
}
