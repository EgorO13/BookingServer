package Handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type dummyLoginRequest struct {
	Role string `json:"role"`
}

type tokenResponse struct {
	Token string `json:"token"`
}

type CustomClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func DummyLoginHandler(secret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dummyLoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			SendErrorResponse(w, "INVALID_REQUEST", "invalid request body", http.StatusBadRequest)
			return
		}
		if req.Role != "admin" && req.Role != "user" {
			SendErrorResponse(w, "INVALID_REQUEST", "role must be admin or user", http.StatusBadRequest)
			return
		}

		var userID string
		if req.Role == "admin" {
			userID = "00000000-0000-0000-0000-000000000001"
		} else {
			userID = "00000000-0000-0000-0000-000000000002"
		}

		claims := CustomClaims{
			Role: req.Role,
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:   userID,
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(secret))
		if err != nil {
			SendErrorResponse(w, "INTERNAL_ERROR", "failed to generate token", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tokenResponse{Token: tokenString})
	}
}
