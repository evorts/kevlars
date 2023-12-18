/**
 * @Author: steven
 * @Description:
 * @File: vars
 * @Date: 22/12/23 11.01
 */

package requests

import (
	"fmt"
	"net/http"
)

type RequestContext string

const (
	ContextClientId RequestContext = "client_id"
	ContextId       RequestContext = "req_id"
)

func (t RequestContext) String() string {
	return string(t)
}

var (
	HttpString = func(httpCode int) string {
		return fmt.Sprintf("%d: %s", httpCode, http.StatusText(httpCode))
	}
)
