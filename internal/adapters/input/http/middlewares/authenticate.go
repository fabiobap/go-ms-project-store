package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-ms-project-store/internal/pkg/errs"
	"github.com/go-ms-project-store/internal/pkg/helpers"
)

const USER_ID_CONTEXT_KEY = "user_id"

type TokenValidator interface {
	GetTokenAbilities(fullToken string) ([]string, *errs.AppError)
	ValidateToken(token string) (uint64, *errs.AppError)
}

type AuthMiddleware struct {
	tokenValidator TokenValidator
}

func NewAuthMiddleware(validator TokenValidator) *AuthMiddleware {
	return &AuthMiddleware{
		tokenValidator: validator,
	}
}

func (am *AuthMiddleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := am.validateBearerToken(r)
		if err != nil {
			helpers.WriteResponse(w, http.StatusOK, err.AsMessage())
			return
		}

		ctx := context.WithValue(r.Context(), USER_ID_CONTEXT_KEY, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (am *AuthMiddleware) validateBearerToken(r *http.Request) (uint64, *errs.AppError) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return 0, errs.NewUnauthorizedError("unauthorized: no token provided")
	}

	// Check Bearer scheme
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return 0, errs.NewUnauthorizedError("unauthorized: invalid token format")
	}

	// Extract token
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		return 0, errs.NewUnauthorizedError("unauthorized: token is empty")
	}

	// Validate token and get user ID
	userID, appErr := am.tokenValidator.ValidateToken(token)
	if appErr != nil {
		return 0, errs.NewUnauthorizedError("unauthorized: " + appErr.Message)
	}

	return userID, nil
}

func GetUserID(ctx context.Context) (uint64, bool) {
	userID, ok := ctx.Value(USER_ID_CONTEXT_KEY).(uint64)
	return userID, ok
}
