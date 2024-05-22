/**
 * @Author: steven
 * @Description:
 * @File: request
 * @Version: 1.0.0
 * @Date: 09/06/23 18.46
 */

package requests

import (
	"context"
	"github.com/evorts/kevlars/rules"
	"github.com/evorts/kevlars/rules/eval"
	"github.com/gofrs/uuid/v5"
	"github.com/labstack/echo/v4"
)

// IsHttp1xx check http code is between 100 and 200
func IsHttp1xx(httpCode int) bool {
	return httpCode >= 100 && httpCode < 200
}

// IsHttp2xx check http code is between 200 and 300
func IsHttp2xx(httpCode int) bool {
	return httpCode >= 200 && httpCode < 300
}

// IsHttp3xx check http code is between 300 and 400
func IsHttp3xx(httpCode int) bool {
	return httpCode >= 300 && httpCode < 400
}

// IsHttp4xx check http code is between 400 and 500
func IsHttp4xx(httpCode int) bool {
	return httpCode >= 400 && httpCode < 500
}

// IsHttp5xx check http code is above inclusive 500
func IsHttp5xx(httpCode int) bool {
	return httpCode >= 500
}

// IsHttpError treat error when http code above inclusive 300
func IsHttpError(httpCode int) bool {
	return !IsHttp2xx(httpCode) && !IsHttp1xx(httpCode)
}

// IdEcho get request id from echo context
func IdEcho(ec echo.Context) string {
	return rules.WhenTrueR1(eval.IsNil(ec.Get(ContextRequestId.String())), func() string {
		return ""
	}, func() string {
		return ec.Get(ContextRequestId.String()).(string)
	})
}

// IdEchoWithSuffix get request id from echo context and add suffix
func IdEchoWithSuffix(ec echo.Context, suffix string) string {
	return IdEcho(ec) + "-" + suffix
}

// IdEchoUid get request id from echo context and convert into uuid
func IdEchoUid(ec echo.Context) string {
	id := IdEcho(ec)
	if len(id) < 1 {
		return id
	}
	if uid, err := uuid.FromString(id); err == nil {
		return uid.String()
	}
	return id
}

// IdEchoUidWithSuffix get request id from echo context, convert into uuid and add suffix
func IdEchoUidWithSuffix(ec echo.Context, suffix string) string {
	return IdEchoUid(ec) + "-" + suffix
}

// ClientIdEcho get client id from echo context
func ClientIdEcho(ec echo.Context) string {
	return rules.WhenTrueR1(eval.IsNil(ec.Get(ContextClientId.String())), func() string {
		return ""
	}, func() string {
		return ec.Get(ContextClientId.String()).(string)
	})
}

// Id get request id from context
func Id(ctx context.Context) string {
	return rules.WhenTrueR1(eval.IsNil(ctx.Value(ContextRequestId)), func() string {
		return ""
	}, func() string {
		return ctx.Value(ContextRequestId).(string)
	})
}

// IdUid get request id from context and convert it into uuid
func IdUid(ctx context.Context) string {
	rs := Id(ctx)
	if len(rs) < 1 {
		return rs
	}
	if uid, err := uuid.FromString(rs); err == nil {
		return uid.String()
	}
	return rs
}

// IdUidWithSuffix get request id from context, convert it into uuid and add suffix
func IdUidWithSuffix(ctx context.Context, suffix string) string {
	return IdUid(ctx) + "-" + suffix
}

// ClientId get client id from context
func ClientId(ctx context.Context) string {
	return rules.WhenTrueR1(eval.IsNil(ctx.Value(ContextClientId)), func() string {
		return ""
	}, func() string {
		return ctx.Value(ContextClientId).(string)
	})
}

// GenerateUUIDV4 produce UUIDv4 string value
func GenerateUUIDV4() string {
	return uuid.Must(uuid.NewV4()).String()
}
