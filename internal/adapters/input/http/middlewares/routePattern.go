package middlewares

import (
	"context"
	"net/http"
)

type RouteContextKey string

const RoutePatternKey RouteContextKey = "routePattern"

func StoreRoutePattern(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fullPath := r.URL.Path
		ctx := context.WithValue(r.Context(), RoutePatternKey, fullPath)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
