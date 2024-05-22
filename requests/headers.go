/**
 * @Author: steven
 * @Description:
 * @File: headers
 * @Date: 22/12/23 11.05
 */

package requests

type Header string

const (
	HeaderAuthorization       Header = "Authorization"
	HeaderAuthorizationBearer Header = "Bearer"

	HeaderApiKey         Header = "X-API-KEY"
	HeaderClientId       Header = "X-CLIENT-ID"
	HeaderSignature      Header = "X-SIGNATURE"
	HeaderIdempotencyKey Header = "X-IDEMPOTENCY-KEY"
)

func (h Header) String() string { return string(h) }
