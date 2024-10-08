package profile

import (
	"ticket-app/config"
	"ticket-app/internal/repository"
	middleware "ticket-app/pkg/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterProfileRoutes(r chi.Router, app config.AppConfig, role string) {
	profileRepo := repository.NewUserRepository(app)
	profileUsecase := NewProfileService(profileRepo)
	profileHandler := NewProfileHandler(profileUsecase)

	r.With(middleware.JWTMiddleware(role)).Route("/profile", func(r chi.Router) {
		r.Get("/me", profileHandler.GetMe)
		r.Post("/update", profileHandler.UpdateUser)
		r.Post("/delete", profileHandler.DeleteUser)
	})
}
