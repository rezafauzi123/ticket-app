package router

import (
	"net/http"
	"ticket-app/config"
	admin "ticket-app/internal/router/admin"
	customer "ticket-app/internal/router/customer"
	"ticket-app/pkg/log"
	middleware "ticket-app/pkg/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func Router(app config.AppConfig) {
	logger := log.GetLogger()
	router := chi.NewRouter()

	if logger == nil {
		panic("log error")
	}

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	router.Use(middleware.JSONContentTypeMiddleware)

	customer.SetupCustomerRoutes(router, app)
	admin.SetupAdminRoutes(router, app)

	logger.Infof("%s", "Starting server on port: 8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		logger.Error("Failed to start server:", err)
	}
}
