package auth

import (
	"encoding/json"
	"net/http"
	"ticket-app/config"
	"ticket-app/internal/module/admin/ticket/request"
	"ticket-app/internal/module/admin/ticket/response"
	"ticket-app/pkg/constant"
	"ticket-app/pkg/helpers"
	"ticket-app/pkg/rabbitmq"
	baseResponse "ticket-app/pkg/response"
	"ticket-app/pkg/utils"

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

func (h *TicketHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	eventID := r.URL.Query().Get("event_id")
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

	filters := request.GetTicketRequest{
		EventID:   eventID,
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

	datas, err := h.ticketService.FindAll(filters)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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

func (h *TicketHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req request.CreateTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(baseResponse.ErrorResponse(constant.INVALID_REQUEST))
		return
	}

	if err := h.validate.Struct(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(baseResponse.ErrorResponse(helpers.FormatValidationError(err)))
		return
	}

	data, err := h.ticketService.Create(req)
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

	go rabbitmq.PublishMessage((*rabbitmq.RabbitMQConnection)(&h.config.RabbitMQConn), rabbitmqResponse, "admin_create_ticket")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(baseResponse.SuccessResponse(constant.SUCCESS_CREATE, ticketResponse))
}

func (h *TicketHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req request.UpdateTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(baseResponse.ErrorResponse(constant.INVALID_REQUEST))
		return
	}

	if err := h.validate.Struct(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(baseResponse.ErrorResponse(helpers.FormatValidationError(err)))
		return
	}

	data, err := h.ticketService.Update(req, id)
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

	go rabbitmq.PublishMessage((*rabbitmq.RabbitMQConnection)(&h.config.RabbitMQConn), rabbitmqResponse, "admin_update_ticket")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(baseResponse.SuccessResponse(constant.SUCCESS_UPDATE, ticketResponse))
}
