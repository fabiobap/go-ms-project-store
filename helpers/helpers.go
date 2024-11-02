package helpers

import (
	"fmt"
	"net/http"
	"time"
)

func DatetimeToString(field time.Time) string {
	return field.Format("2006-01-02 15:04:05")
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
