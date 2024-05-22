package midware

import (
	"github.com/evorts/kevlars/requests"
	"net/http"
	"strings"
)

func apiKeyFromHeader(r *http.Request) string {
	return r.Header.Get(requests.HeaderApiKey.String())
}

func tokenFromHeader(r *http.Request) string {
	return strings.TrimSpace(
		strings.TrimPrefix(
			r.Header.Get(requests.HeaderAuthorization.String()),
			requests.HeaderAuthorizationBearer.String(),
		),
	)
}
