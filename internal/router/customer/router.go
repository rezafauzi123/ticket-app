package customer

import (
	"strconv"
	"ticket-app/config"
	"ticket-app/pkg/constant"

	auth "ticket-app/internal/module/customer/auth"
	event "ticket-app/internal/module/customer/event"
	notif "ticket-app/internal/module/customer/notification"
	payment "ticket-app/internal/module/customer/payment"
	profile "ticket-app/internal/module/customer/profile"
	ticket "ticket-app/internal/module/customer/ticket"

	"github.com/go-chi/chi/v5"
)

func SetupCustomerRoutes(r *chi.Mux, app config.AppConfig) {
	role := strconv.Itoa(constant.CUSTOMER)
	r.Route("/api/v1", func(r chi.Router) {
		profile.RegisterProfileRoutes(r, app, role)
		auth.RegisterAuthRoutes(r, app, role)
		ticket.RegisterTicketRoutes(r, app, role)
		event.RegisterEventRoutes(r, app, role)
		payment.RegisterPaymentRoutes(r, app, role)
		notif.RegisterNotificationRoutes(r, app, role)
	})
}
