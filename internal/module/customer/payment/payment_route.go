package auth

import (
	"ticket-app/config"
	"ticket-app/internal/repository"
	middleware "ticket-app/pkg/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterPaymentRoutes(r chi.Router, app config.AppConfig, role string) {
	paymentRepo := repository.NewPaymentRepository(app)
	ticketRepo := repository.NewTicketRepository(app)
	paymentService := NewPaymentService(paymentRepo, ticketRepo, app)
	paymentHandler := NewPaymentHandler(paymentService, app)

	r.With(middleware.JWTMiddleware(role)).Route("/payment", func(r chi.Router) {
		r.Post("/process-payment", paymentHandler.ProcessPayment)
		r.Post("/cancel/{id}", paymentHandler.CancelPayment)
	})
}
