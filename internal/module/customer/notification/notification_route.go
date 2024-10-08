package notification

import (
	"ticket-app/config"
	"ticket-app/internal/repository"
	middleware "ticket-app/pkg/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterNotificationRoutes(r chi.Router, app config.AppConfig, role string) {
	notificationRepo := repository.NewNotificationRepository(app)
	notificationService := NewNotificationService(notificationRepo, app)
	notificationHandler := NewNotificationHandler(notificationService, app)

	r.With(middleware.JWTMiddleware(role)).Route("/notification", func(r chi.Router) {
		r.Get("/", notificationHandler.SendNotification)
	})
}
