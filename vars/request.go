/**
 * @Author: steven
 * @Description:
 * @File: request
 * @Version: 1.0.0
 * @Date: 08/06/23 18.23
 */

package vars

import (
	"fmt"
	"net/http"
)

const (
	ClientIdUnknown string = "<client_unknown>"
)

var (
	HttpString = func(httpCode int) string {
		return fmt.Sprintf("%d: %s", httpCode, http.StatusText(httpCode))
	}
)

type RequestContext string

const (
	RequestContextClientId RequestContext = "client_id"
	RequestContextId       RequestContext = "req_id"
)

func (t RequestContext) String() string {
	return string(t)
}
