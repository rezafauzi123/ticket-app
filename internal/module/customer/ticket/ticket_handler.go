package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"ticket-app/config"
	"ticket-app/internal/module/customer/ticket/request"
	"ticket-app/internal/module/customer/ticket/response"
	"ticket-app/pkg/constant"
	"ticket-app/pkg/helpers"
	middleware "ticket-app/pkg/middleware"
	"ticket-app/pkg/rabbitmq"
	baseResponse "ticket-app/pkg/response"

	validator "github.com/go-playground/validator/v10"
)

type TicketHandler struct {
	ticketService TicketService
	validate      *validator.Validate
	config        config.AppConfig
}

func NewTicketHandler(ticketService TicketService, config config.AppConfig) *TicketHandler {
	return &TicketHandler{
		ticketService: ticketService,
		validate:      validator.New(),
		config:        config,
	}
}

func (h *TicketHandler) BuyTicket(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserContextKey("userID")).(string)
	var req request.BuyTicketRequest
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

	data, err := h.ticketService.BuyTicket(req, userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(baseResponse.ErrorResponse(err.Error()))
		return
	}

	var ticketResponse response.TicketResponse
	ticketResponse.TicketEntityToTicketResponse(data)

	rabbitmqResponse, err := rabbitmq.MappingJsonToRabbitMQMessage(ticketResponse)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(baseResponse.ErrorResponse(err.Error()))
		return
	}

	go rabbitmq.PublishMessage((*rabbitmq.RabbitMQConnection)(&h.config.RabbitMQConn), rabbitmqResponse, "ticket_created")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(baseResponse.SuccessResponse(constant.SUCCESS_CREATE, ticketResponse))
}

func (h *TicketHandler) CancelTicket(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	data, err := h.ticketService.CancelTicket(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(baseResponse.ErrorResponse(err.Error()))
		return
	}

	var ticketResponse response.TicketResponse
	ticketResponse.TicketEntityToTicketResponse(data)

	rabbitmqResponse, err := rabbitmq.MappingJsonToRabbitMQMessage(ticketResponse)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(baseResponse.ErrorResponse(err.Error()))
		return
	}

	go rabbitmq.PublishMessage((*rabbitmq.RabbitMQConnection)(&h.config.RabbitMQConn), rabbitmqResponse, "ticket_canceled")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(baseResponse.SuccessResponse(constant.SUCCESS_UPDATE, ticketResponse))
}

func (h *TicketHandler) FindByUserID(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserContextKey("userID")).(string)
	datas, err := h.ticketService.FindByUserID(userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(baseResponse.ErrorResponse(err.Error()))
		return
	}

	var ticketResponses []*response.TicketResponse
	for _, data := range datas {
		var ticketResponse response.TicketResponse
		ticketResponse.TicketEntityToTicketResponse(data)

		ticketResponses = append(ticketResponses, &ticketResponse)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(baseResponse.SuccessResponse(constant.SUCCESS_GET, ticketResponses))
}

func (h *TicketHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	data, err := h.ticketService.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(baseResponse.ErrorResponse(constant.DATA_NOT_FOUND))
		return
	}

	var ticketResponse response.TicketResponse
	ticketResponse.TicketEntityToTicketResponse(data)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(baseResponse.SuccessResponse(constant.SUCCESS_GET, ticketResponse))
}
