package rest

import (
	"context"
	"net/http"
	"strings"

	"financial_system/pkg/jwt"
)

type contextKey string

const (
	ContextKeyUserID contextKey = "user_id"
	ContextKeyRole   contextKey = "role"
)

// authMiddleware validates JWT and role; on success sets user_id and role in context.
func (h *Handler) authMiddleware(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				respondError(w, http.StatusUnauthorized, "требуется авторизация")
				return
			}

			tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
			claims, err := jwt.ParseToken(h.jwtSecret, tokenString)
			if err != nil {
				if err == jwt.ErrInvalidToken {
					respondError(w, http.StatusUnauthorized, "неверный или просроченный токен")
					return
				}
				respondError(w, http.StatusUnauthorized, "неверный или просроченный токен")
				return
			}

			if claims.Role != requiredRole {
				respondError(w, http.StatusForbidden, "недостаточно прав")
				return
			}

			ctx := context.WithValue(r.Context(), ContextKeyUserID, claims.UserID)
			ctx = context.WithValue(ctx, ContextKeyRole, claims.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// userIDFromRequest returns user_id from context; 0 if not set.
func userIDFromRequest(r *http.Request) int {
	v, _ := r.Context().Value(ContextKeyUserID).(int)
	return v
}
