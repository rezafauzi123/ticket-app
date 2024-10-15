package auth

import (
	"ticket-app/config"
	event "ticket-app/internal/module/customer/event"
	"ticket-app/internal/repository"
	middleware "ticket-app/pkg/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterTicketRoutes(r chi.Router, app config.AppConfig, role string) {
	ticketRepo := repository.NewTicketRepository(app)
	eventRepo := repository.NewEventRepository(app)
	eventService := event.NewEventService(eventRepo, app)
	ticketService := NewTicketService(ticketRepo, eventService, app)
	ticketHandler := NewTicketHandler(ticketService, app)

	r.With(middleware.JWTMiddleware(role)).Route("/ticket", func(r chi.Router) {
		r.Get("/", ticketHandler.FindByUserID)
		r.Get("/detail/{id}", ticketHandler.FindByID)
		r.Post("/buy-ticket", ticketHandler.BuyTicket)
		r.Post("/cancel-ticket/{id}", ticketHandler.CancelTicket)
	})
}
