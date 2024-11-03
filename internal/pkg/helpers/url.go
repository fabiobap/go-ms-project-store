package helpers

import (
	"fmt"
	"net/http"

	"github.com/go-ms-project-store/internal/adapters/input/http/middlewares"
)

func GetCurrentUri(r *http.Request) string {
	return r.Context().Value(middlewares.RoutePatternKey).(string)
}

func GetBaseURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}

	host := r.Host
	if forwardedHost := r.Header.Get("X-Forwarded-Host"); forwardedHost != "" {
		host = forwardedHost
	}

	return fmt.Sprintf("%s://%s", scheme, host)
}

func GetFullRouteUrl(r *http.Request) string {
	return GetBaseURL(r) + GetCurrentUri(r)
}
