package auth

import (
	"context"
	"net/http"
	"strings"

	"D/Go/messenger/internal/platform/httpx"
)

type AuthMiddleware func(http.Handler) http.Handler

func Auth(verifier JWTVerifier) AuthMiddleware {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				httpx.WriteErrorResponse(w, httpx.HTTPError{
					Code:    401,
					Message: `{"error":"missing token"}`,
				}, nil)
				return
			}

			parts := strings.Split(authHeader, " ")

			if len(parts) != 2 || parts[0] != "Bearer" {
				httpx.WriteErrorResponse(w, httpx.HTTPError{
					Code:    401,
					Message: `{"error":"invalid token"}`,
				}, nil)
				return
			}

			userID, err := verifier.ParseAccess(parts[1])
			if err != nil {
				httpx.WriteErrorResponse(w, httpx.HTTPError{
					Code:    401,
					Message: `{"error":"invalid token"}`,
				}, nil)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
