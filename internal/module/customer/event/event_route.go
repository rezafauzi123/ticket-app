package auth

import (
	"ticket-app/config"
	"ticket-app/internal/repository"
	middleware "ticket-app/pkg/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterEventRoutes(r chi.Router, app config.AppConfig, role string) {
	eventRepo := repository.NewEventRepository(app)
	eventService := NewEventService(eventRepo, app)
	eventHandler := NewEventHandler(eventService, app)

	r.With(middleware.JWTMiddleware(role)).Route("/event", func(r chi.Router) {
		r.Get("/", eventHandler.FindAll)
		r.Get("/detail/{id}", eventHandler.FindByID)
	})
}
