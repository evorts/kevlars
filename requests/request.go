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
	"github.com/evorts/kevlars/utils"
	"github.com/gofrs/uuid/v5"
	"github.com/labstack/echo/v4"
)

func IsHttp1xx(httpCode int) bool {
	return httpCode >= 100 && httpCode < 200
}

func IsHttp2xx(httpCode int) bool {
	return httpCode >= 200 && httpCode < 300
}

func IsHttp3xx(httpCode int) bool {
	return httpCode >= 300 && httpCode < 400
}

func IsHttp4xx(httpCode int) bool {
	return httpCode >= 400 && httpCode < 500
}

func IsHttp5xx(httpCode int) bool {
	return httpCode >= 500
}

func IsHttpError(httpCode int) bool {
	return !IsHttp2xx(httpCode) && !IsHttp1xx(httpCode)
}

func IdEC(ec echo.Context) string {
	return utils.IfNil(ec.Get(ContextId.String()), "")
}

func IdECWithFlag(ec echo.Context, parse bool) string {
	id := utils.IfNil(ec.Get(ContextId.String()), "")
	if len(id) < 1 || !parse {
		return id
	}
	if uid, err := uuid.FromString(id); err == nil {
		return uid.String()
	}
	return id
}

func IdECWithFlagAndSuffix(ec echo.Context, parse bool, suffix string) string {
	return IdECWithFlag(ec, parse) + "-" + suffix
}

func ClientIdEC(ec echo.Context) string {
	return utils.IfNil(ec.Get(ContextClientId.String()), "")
}

func Id(ctx context.Context) string {
	rs := utils.IfNil(ctx.Value(ContextId), "")
	return rs
}

func IdWithFlag(ctx context.Context, parse bool) string {
	rs := utils.IfNil(ctx.Value(ContextId), "")
	if len(rs) < 1 || !parse {
		return rs
	}
	if uid, err := uuid.FromString(rs); err == nil {
		return uid.String()
	}
	return rs
}

func IdWithFlagAndSuffix(ctx context.Context, parse bool, suffix string) string {
	return IdWithFlag(ctx, parse) + "-" + suffix
}

func ClientId(ctx context.Context) string {
	return utils.IfNil(ctx.Value(ContextClientId), "")
}

func GenerateUUIDV4() string {
	return uuid.Must(uuid.NewV4()).String()
}
