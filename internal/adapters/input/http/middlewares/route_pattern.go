package middlewares

import (
	"context"
	"net/http"

	"github.com/go-ms-project-store/internal/core/enums"
)

func StoreRoutePattern(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fullPath := r.URL.Path
		ctx := context.WithValue(r.Context(), enums.RoutePatternKey, fullPath)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
