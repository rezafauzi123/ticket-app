package auth

import (
	"ticket-app/config"
	event "ticket-app/internal/module/customer/event"
	ticket "ticket-app/internal/module/customer/ticket"
	"ticket-app/internal/repository"
	middleware "ticket-app/pkg/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterPaymentRoutes(r chi.Router, app config.AppConfig, role string) {
	paymentRepo := repository.NewPaymentRepository(app)
	ticketRepo := repository.NewTicketRepository(app)
	eventRepo := repository.NewEventRepository(app)
	eventService := event.NewEventService(eventRepo, app)
	ticketService := ticket.NewTicketService(ticketRepo, eventService, app)
	paymentService := NewPaymentService(paymentRepo, ticketService, app)
	paymentHandler := NewPaymentHandler(paymentService, app)

	r.With(middleware.JWTMiddleware(role)).Route("/payment", func(r chi.Router) {
		r.Post("/process-payment", paymentHandler.ProcessPayment)
		r.Post("/cancel/{id}", paymentHandler.CancelPayment)
	})
}
