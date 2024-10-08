package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"ticket-app/pkg/constant"
	"ticket-app/pkg/response"

	"github.com/dgrijalva/jwt-go"
)

type UserContextKey string

func JWTMiddleware(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(response.ErrorResponse(constant.UNAUTHORIZED))
				return
			}

			tokenString := strings.TrimSpace(strings.Replace(authHeader, "Bearer", "", 1))
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("JWT_SECRET")), nil
			})

			if err != nil || !token.Valid {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(response.ErrorResponse(constant.INVALID_TOKEN))
				return
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				userID := claims["user_id"].(string)

				userRole := claims["role_id"].(string)
				if userRole != requiredRole {
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(response.ErrorResponse(constant.UNAUTHORIZED))
					return
				}

				ctx := context.WithValue(r.Context(), UserContextKey("userID"), userID)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(response.UnauthorizedResponse(constant.INVALID_TOKEN))
			}
		})
	}
}

func JSONContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func GetUserID(ctx context.Context) string {
	return ctx.Value(UserContextKey("userID")).(string)
}
