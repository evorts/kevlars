/**
 * @Author: steven
 * @Description:
 * @File: echo_header
 * @Version: 1.0.0
 * @Date: 17/04/23 16.09
 */

package midware

import (
	"github.com/evorts/kevlars/logger"
	"github.com/evorts/kevlars/requests"
	"github.com/labstack/echo/v4"
)

var securityHeaders = map[string]string{
	"Strict-Transport-Security":         "max-age=63072000; includeSubDomains; preload",
	"X-Frame-Options":                   "DENY",
	"X-Content-Type-Options":            "nosniff",
	"Content-Security-Policy":           "default-src 'self'; object-src 'none'; frame-ancestors 'none'; upgrade-insecure-requests; block-all-mixed-content",
	"X-Permitted-Cross-Domain-Policies": "none",
	"Referrer-Policy":                   "no-referrer",
	"Clear-Site-Data":                   "\"cache\", \"cookies\", \"storage\"",
	"Cross-Origin-Embedder-Policy":      "require-corp",
	"Cross-Origin-Opener-Policy":        "same-origin",
	"Cross-Origin-Resource-Policy":      "same-origin",
	"Permissions-Policy":                "accelerometer=(),ambient-light-sensor=(),autoplay=(),battery=(),camera=(),display-capture=(),document-domain=(),encrypted-media=(),fullscreen=(),gamepad=(),geolocation=(),gyroscope=(),layout-animations=(self),legacy-image-formats=(self),magnetometer=(),microphone=(),midi=(),oversized-images=(self),payment=(),picture-in-picture=(),publickey-credentials-get=(),speaker-selection=(),sync-xhr=(self),unoptimized-images=(self),unsized-media=(self),usb=(),screen-wake-lock=(),web-share=(),xr-spatial-tracking=()",
	"Cache-Control":                     "no-store, max-age=0",
	"Pragma":                            "no-cache",
}

func EchoWithSecurityHeader(overrides map[string]string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if len(securityHeaders) < 1 {
				return next(c)
			}
			resp := c.Response()
			for k, v := range securityHeaders {
				if overrides != nil {
					if overrideValue, ok := overrides[k]; ok {
						v = overrideValue
					}
				}
				if len(v) < 1 {
					continue
				}
				resp.Header().Add(k, v)
			}
			return next(c)
		}
	}
}

func EchoResponseWithProperRequestId(log logger.Manager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			resp := c.Response()
			req := c.Request()
			reqIdEcho := requests.IdEC(c)
			reqIdContext := requests.Id(req.Context())
			log.InfoWithProps(map[string]interface{}{
				"req_id_echo":    reqIdEcho,
				"req_id_context": reqIdContext,
			}, "response request id")
			resp.Header().Set(echo.HeaderXRequestID, reqIdEcho)
			return next(c)
		}
	}
}
