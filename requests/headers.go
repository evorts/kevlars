/**
 * @Author: steven
 * @Description:
 * @File: headers
 * @Date: 22/12/23 11.05
 */

package requests

type Header string

const (
	HeaderRestApiKey         Header = "X-REST-API-KEY"
	HeaderRestClientId       Header = "X-REST-CLIENT-ID"
	HeaderRestSignature      Header = "X-REST-SIGNATURE"
	HeaderRestIdempotencyKey Header = "X-REST-IDEMPOTENCY-KEY"
)
