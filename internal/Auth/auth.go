package Auth

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	UserIDKey contextKey = "user_id"
	RoleKey   contextKey = "role"
)

type CustomClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func AuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				sendError(w, "UNAUTHORIZED", "missing authorization header", http.StatusUnauthorized)
				return
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				sendError(w, "UNAUTHORIZED", "invalid authorization header", http.StatusUnauthorized)
				return
			}
			tokenString := parts[1]
			claims := &CustomClaims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})
			if err != nil || !token.Valid {
				sendError(w, "UNAUTHORIZED", "invalid token", http.StatusUnauthorized)
				return
			}
			userID := claims.Subject
			if userID == "" {
				sendError(w, "UNAUTHORIZED", "user_id not found in token", http.StatusUnauthorized)
				return
			}
			role := claims.Role
			if role == "" {
				sendError(w, "UNAUTHORIZED", "role not found in token", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			ctx = context.WithValue(ctx, RoleKey, role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AdminOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role := r.Context().Value(RoleKey).(string)
		if role != "admin" {
			sendError(w, "FORBIDDEN", "admin role required", http.StatusForbidden)
			return
		}
		next(w, r)
	}
}

func UserOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role := r.Context().Value(RoleKey).(string)
		if role != "user" {
			sendError(w, "FORBIDDEN", "user role required", http.StatusForbidden)
			return
		}
		next(w, r)
	}
}

func sendError(w http.ResponseWriter, code, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]string{
			"code":    code,
			"message": message,
		},
	})
}
