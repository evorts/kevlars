package midware

import (
	"fmt"
	"github.com/evorts/kevlars/logger"
	"github.com/evorts/kevlars/rules"
	"github.com/evorts/kevlars/utils"
	"github.com/labstack/echo/v4"
	"math"
	"net/http"
	"os"
	"time"
)

// EchoLogger is the logrus logger handler
func EchoLogger(l logger.Manager, pathNotLogged ...string) echo.MiddlewareFunc {
	var echoSkipper = func(c echo.Context) bool {
		path := c.Request().URL.Path
		return utils.InArray(pathNotLogged, path)
	}
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if echoSkipper(c) {
				return next(c)
			}
			if err = next(c); err != nil {
				c.Error(err)
			}
			req := c.Request()
			resp := c.Response()
			path := req.URL.Path
			start := time.Now()
			stop := time.Since(start)
			latency := int(math.Ceil(float64(stop.Nanoseconds()) / 1000000.0))
			statusCode := resp.Status
			clientIP := c.RealIP()
			clientUserAgent := req.UserAgent()
			referer := req.Referer()
			dataLength := resp.Size
			if dataLength < 0 {
				dataLength = 0
			}
			clientID := c.Get("client.id")
			entry := map[string]interface{}{
				"hostname":    hostname,
				"statusCode":  statusCode,
				"latency":     latency, // time to process
				"clientIP":    clientIP,
				"method":      req.Method,
				"path":        path,
				"referer":     referer,
				"dataLength":  dataLength,
				"userAgent":   clientUserAgent,
				"accept":      req.Header.Get("Accept"),
				"contentType": req.Header.Get("Content-Type"),
				"clientID": rules.WhenTrueR1[string](clientID == nil, func() string { return "<undefined>" }, func() string {
					if v, ok := clientID.(string); ok {
						return v
					}
					return "<undefined>"
				}),
			}
			if err != nil {
				l.ErrorWithProps(entry, err.Error())
			} else {
				msg := fmt.Sprintf("%s - %s [%s] \"%s %s\" %d %d \"%s\" \"%s\" (%dms)",
					clientIP, hostname, time.Now().Format(timeFormat), c.Request().Method,
					path, statusCode, dataLength, referer, clientUserAgent, latency)
				if statusCode >= http.StatusInternalServerError {
					l.ErrorWithProps(entry, msg)
				} else if statusCode >= http.StatusBadRequest {
					l.WarnWithProps(entry, msg)
				} else {
					l.InfoWithProps(entry, msg)
				}
			}
			return
		}
	}
}
