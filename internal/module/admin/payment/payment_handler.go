package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"ticket-app/config"
	"ticket-app/internal/module/admin/payment/request"
	"ticket-app/internal/module/admin/payment/response"
	"ticket-app/pkg/constant"
	"ticket-app/pkg/helpers"
	"ticket-app/pkg/rabbitmq"
	baseResponse "ticket-app/pkg/response"
	"ticket-app/pkg/utils"

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

func (h *PaymentHandler) ConfirmPayment(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
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

	data, err := h.paymentService.ConfirmPayment(id, req)
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

	go rabbitmq.PublishMessage((*rabbitmq.RabbitMQConnection)(&h.config.RabbitMQConn), rabbitmqResponse, "payment_confirmed")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(baseResponse.SuccessResponse(constant.SUCCESS_CREATE, paymentResponse))
}

func (h *PaymentHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	data, err := h.paymentService.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(baseResponse.ErrorResponse(err.Error()))
		return
	}

	var paymentResponse response.PaymentResponse
	paymentResponse.PaymentEntityToPaymentResponse(data)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(baseResponse.SuccessResponse(constant.SUCCESS_GET, paymentResponse))
}

func (h *PaymentHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	ticketID := r.URL.Query().Get("ticket_id")
	userID := r.URL.Query().Get("user_id")
	status := r.URL.Query().Get("status")
	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")

	pagination, err := utils.ExtractPagination(r, 1, 10, 100)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(baseResponse.ErrorResponse(err.Error()))
		return
	}

	filters := request.GetPaymentRequest{
		TicketID:  ticketID,
		UserID:    userID,
		Status:    status,
		SortBy:    sortBy,
		SortOrder: sortOrder,
		Pagination: utils.Pagination{
			Page:    pagination.Page,
			PerPage: pagination.PerPage,
		},
	}

	if err := h.validate.Struct(filters); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(baseResponse.ErrorResponse(helpers.FormatValidationError(err)))
		return
	}

	datas, err := h.paymentService.FindAll(filters)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(baseResponse.ErrorResponse(err.Error()))
		return
	}

	var paymentResponses []*response.PaymentResponse
	for _, data := range datas {
		var paymentResponse response.PaymentResponse
		paymentResponse.PaymentEntityToPaymentResponse(data)
		paymentResponses = append(paymentResponses, &paymentResponse)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(baseResponse.SuccessResponse(constant.SUCCESS_GET, paymentResponses))
}
