package midware

import (
	"net/http"
	"strings"
)

func apiKeyFromHeader(r *http.Request) string {
	return r.Header.Get("X-API-KEY")
}

func tokenFromHeader(r *http.Request) string {
	return strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
}
