package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/hexzedels/gosdlworkshop/internal/auth"
)

type contextKey string

const ClaimsKey contextKey = "claims"

// RequireJWT is middleware that validates the Authorization header.
func RequireJWT(secret string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			http.Error(w, `{"error":"missing authorization header"}`, http.StatusUnauthorized)
			return
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")

		claims, err := auth.ValidateJWT(secret, tokenStr)
		if err != nil {
			http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ClaimsKey, claims)
		next(w, r.WithContext(ctx))
	}
}

// RequireAdmin wraps RequireJWT and additionally checks that the user has admin role.
func RequireAdmin(secret string, next http.HandlerFunc) http.HandlerFunc {
	return RequireJWT(secret, func(w http.ResponseWriter, r *http.Request) {
		claims := r.Context().Value(ClaimsKey).(*auth.Claims)
		if claims.Role != "admin" {
			http.Error(w, `{"error":"admin access required"}`, http.StatusForbidden)
			return
		}
		next(w, r)
	})
}

// GetClaims extracts JWT claims from the request context.
func GetClaims(r *http.Request) *auth.Claims {
	claims, _ := r.Context().Value(ClaimsKey).(*auth.Claims)
	return claims
}
