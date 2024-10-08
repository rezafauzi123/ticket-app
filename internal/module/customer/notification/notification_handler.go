package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"ticket-app/config"
	"ticket-app/pkg/constant"
	middleware "ticket-app/pkg/middleware"
	"ticket-app/pkg/rabbitmq"
	baseResponse "ticket-app/pkg/response"

	validator "github.com/go-playground/validator/v10"
)

type NotificationHandler struct {
	notificationService NotificationService
	validate            *validator.Validate
	config              config.AppConfig
}

func NewNotificationHandler(notificationService NotificationService, config config.AppConfig) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
		validate:            validator.New(),
		config:              config,
	}
}

func (h *NotificationHandler) SendNotification(w http.ResponseWriter, r *http.Request) {
	go rabbitmq.ReceiveMessages(r.Context(), (*rabbitmq.RabbitMQConnection)(&h.config.RabbitMQConn), "payment_created", h.handleMessage)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(baseResponse.SuccessResponse(constant.SUCCESS_CREATE, nil))
}

func (h *NotificationHandler) handleMessage(ctx context.Context, msg []byte) {
	userID := middleware.GetUserID(ctx)
	_, err := h.notificationService.SendNotification(string(msg), userID)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
}
