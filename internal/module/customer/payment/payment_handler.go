package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"ticket-app/config"
	"ticket-app/internal/module/customer/payment/request"
	"ticket-app/internal/module/customer/payment/response"
	"ticket-app/pkg/constant"
	"ticket-app/pkg/helpers"
	middleware "ticket-app/pkg/middleware"
	"ticket-app/pkg/rabbitmq"
	baseResponse "ticket-app/pkg/response"

	validator "github.com/go-playground/validator/v10"
)

type PaymentHandler struct {
	paymentService PaymentService
	validate       *validator.Validate
	config         config.AppConfig
}

func NewPaymentHandler(paymentService PaymentService, config config.AppConfig) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
		validate:       validator.New(),
		config:         config,
	}
}

func (h *PaymentHandler) ProcessPayment(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserContextKey("userID")).(string)
	var req request.PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(baseResponse.ErrorResponse(constant.INVALID_REQUEST))
		return
	}

	if err := h.validate.Struct(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(baseResponse.ErrorResponse(helpers.FormatValidationError(err)))
		return
	}

	data, err := h.paymentService.ProcessPayment(req, userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(baseResponse.ErrorResponse(err.Error()))
		return
	}

	var paymentResponse response.PaymentResponse
	paymentResponse.PaymentEntityToPaymentResponse(data)

	rabbitmqResponse, err := rabbitmq.MappingJsonToRabbitMQMessage(paymentResponse)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(baseResponse.ErrorResponse(err.Error()))
		return
	}

	go rabbitmq.PublishMessage((*rabbitmq.RabbitMQConnection)(&h.config.RabbitMQConn), rabbitmqResponse, "payment_created")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(baseResponse.SuccessResponse(constant.SUCCESS_CREATE, paymentResponse))
}

func (h *PaymentHandler) CancelPayment(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserContextKey("userID")).(string)
	id := r.PathValue("id")

	data, err := h.paymentService.CancelPayment(id, userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(baseResponse.ErrorResponse(err.Error()))
		return
	}

	var paymentResponse response.PaymentResponse
	paymentResponse.PaymentEntityToPaymentResponse(data)

	rabbitmqResponse, err := rabbitmq.MappingJsonToRabbitMQMessage(paymentResponse)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(baseResponse.ErrorResponse(err.Error()))
		return
	}

	go rabbitmq.PublishMessage((*rabbitmq.RabbitMQConnection)(&h.config.RabbitMQConn), rabbitmqResponse, "payment_canceled")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(baseResponse.SuccessResponse(constant.SUCCESS_CREATE, paymentResponse))
}
