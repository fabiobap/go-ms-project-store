package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-ms-project-store/internal/pkg/errs"
	"github.com/go-ms-project-store/internal/pkg/helpers"
	"github.com/go-ms-project-store/internal/pkg/logger"
)

type AbilityMiddleware struct {
	authRepo TokenValidator
}

func NewAbilityMiddleware(authRepo TokenValidator) *AbilityMiddleware {
	return &AbilityMiddleware{
		authRepo: authRepo,
	}
}

// RequireAbilities checks if the token has all the required abilities
func (am *AbilityMiddleware) RequireAbilities(abilities ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token == "" {
				return
			}

			token = strings.TrimPrefix(token, "Bearer ")

			tokenAbilities, err := am.authRepo.GetTokenAbilities(token)
			if err != nil {
				helpers.WriteResponse(w, http.StatusUnauthorized, err.AsMessage())
				return
			}

			// Check for all required abilities
			for _, requiredAbility := range abilities {
				found := false
				for _, tokenAbility := range tokenAbilities {
					if tokenAbility == requiredAbility {
						found = true
						break
					}
				}
				if !found {
					logger.Error(fmt.Sprintf("missing required ability: %s", requiredAbility))
					helpers.WriteResponse(w, http.StatusForbidden, errs.NewUnauthorizedError("Invalid Token"))
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
