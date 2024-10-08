package router

import (
	"strconv"
	"ticket-app/config"
	"ticket-app/pkg/constant"

	auth "ticket-app/internal/module/admin/auth"
	event "ticket-app/internal/module/admin/event"
	payment "ticket-app/internal/module/admin/payment"
	ticket "ticket-app/internal/module/admin/ticket"

	"github.com/go-chi/chi/v5"
)

func SetupAdminRoutes(r *chi.Mux, app config.AppConfig) {
	role := strconv.Itoa(constant.ADMIN)
	r.Route("/api/admin/v1", func(r chi.Router) {
		auth.RegisterAuthRoutes(r, app, role)
		event.RegisterEventRoutes(r, app, role)
		ticket.RegisterTicketRoutes(r, app, role)
		payment.RegisterPaymentRoutes(r, app, role)
	})
}
