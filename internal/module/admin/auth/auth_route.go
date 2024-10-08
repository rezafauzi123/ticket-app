package auth

import (
	"ticket-app/config"
	"ticket-app/internal/repository"
	middleware "ticket-app/pkg/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterAuthRoutes(r chi.Router, app config.AppConfig, role string) {
	authRepo := repository.NewUserRepository(app)
	authService := NewAuthService(authRepo, app)
	authHandler := NewAuthHandler(authService, app)

	r.Post("/login", authHandler.Login)
	r.With(middleware.JWTMiddleware(role)).Post("/refresh-token", authHandler.RefreshTokenHandler)
}
