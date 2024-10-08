package auth

import (
	"ticket-app/config"
	"ticket-app/internal/repository"
	middleware "ticket-app/pkg/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterTicketRoutes(r chi.Router, app config.AppConfig, role string) {
	ticketRepo := repository.NewTicketRepository(app)
	eventRepo := repository.NewEventRepository(app)
	ticketService := NewTicketService(ticketRepo, eventRepo, app)
	ticketHandler := NewTicketHandler(ticketService, app)

	r.With(middleware.JWTMiddleware(role)).Route("/ticket", func(r chi.Router) {
		r.Get("/", ticketHandler.FindAll)
		r.Get("/detail/{id}", ticketHandler.FindByID)
		r.Post("/create", ticketHandler.Create)
		r.Put("/update/{id}", ticketHandler.Update)
	})
}
