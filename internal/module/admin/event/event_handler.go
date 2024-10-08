package auth

import (
	"encoding/json"
	"net/http"
	"strconv"
	"ticket-app/config"
	"ticket-app/internal/module/admin/event/request"
	"ticket-app/internal/module/admin/event/response"
	"ticket-app/pkg/constant"
	"ticket-app/pkg/helpers"
	middleware "ticket-app/pkg/middleware"
	"ticket-app/pkg/rabbitmq"
	baseResponse "ticket-app/pkg/response"
	"ticket-app/pkg/utils"

	validator "github.com/go-playground/validator/v10"
)

type EventHandler struct {
	eventService EventService
	validate     *validator.Validate
	config       config.AppConfig
}

func NewEventHandler(eventService EventService, config config.AppConfig) *EventHandler {
	return &EventHandler{
		eventService: eventService,
		validate:     validator.New(),
		config:       config,
	}
}

func (h *EventHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req request.CreateEventRequest
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

	data, err := h.eventService.Create(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(baseResponse.ErrorResponse(err.Error()))
		return
	}

	var eventResponse response.EventResponse
	eventResponse.EventEntityToEventResponse(data)

	rabbitmqResponse, err := rabbitmq.MappingJsonToRabbitMQMessage(eventResponse)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(baseResponse.ErrorResponse(err.Error()))
		return
	}

	go rabbitmq.PublishMessage((*rabbitmq.RabbitMQConnection)(&h.config.RabbitMQConn), rabbitmqResponse, "admin_create_event")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(baseResponse.SuccessResponse(constant.SUCCESS_CREATE, eventResponse))
}

func (h *EventHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	location := r.URL.Query().Get("location")
	date := r.URL.Query().Get("date")
	availableTickets := r.URL.Query().Get("available_tickets")
	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")

	var availableTicketsInt *int
	if availableTickets != "" {
		availableTicketsVal, err := strconv.Atoi(availableTickets)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(baseResponse.ErrorResponse("Invalid available_tickets value"))
			return
		}
		availableTicketsInt = &availableTicketsVal
	}

	pagination, err := utils.ExtractPagination(r, 1, 10, 100)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(baseResponse.ErrorResponse(err.Error()))
		return
	}

	filters := request.GetEventRequest{
		Name:             name,
		Location:         location,
		Date:             date,
		AvailableTickets: availableTicketsInt,
		SortBy:           sortBy,
		SortOrder:        sortOrder,
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

	datas, err := h.eventService.FindAll(filters)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(baseResponse.ErrorResponse(err.Error()))
		return
	}

	var eventResponses []*response.EventResponse
	for _, data := range datas {
		var eventResponse response.EventResponse
		eventResponse.EventEntityToEventResponse(data)
		eventResponses = append(eventResponses, &eventResponse)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(baseResponse.SuccessResponse(constant.SUCCESS_GET, eventResponses))
}

func (h *EventHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	data, err := h.eventService.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(baseResponse.ErrorResponse(constant.DATA_NOT_FOUND))
		return
	}

	var eventResponse response.EventResponse
	eventResponse.EventEntityToEventResponse(data)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(baseResponse.SuccessResponse(constant.SUCCESS_GET, eventResponse))
}

func (h *EventHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req request.UpdateEventRequest
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

	data, err := h.eventService.Update(req, id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(baseResponse.ErrorResponse(err.Error()))
		return
	}

	var eventResponse response.EventResponse
	eventResponse.EventEntityToEventResponse(data)

	rabbitmqResponse, err := rabbitmq.MappingJsonToRabbitMQMessage(eventResponse)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(baseResponse.ErrorResponse(err.Error()))
		return
	}

	go rabbitmq.PublishMessage((*rabbitmq.RabbitMQConnection)(&h.config.RabbitMQConn), rabbitmqResponse, "admin_update_event")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(baseResponse.SuccessResponse(constant.SUCCESS_UPDATE, eventResponse))
}

func (h *EventHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userID := r.Context().Value(middleware.UserContextKey("userID")).(string)
	err := h.eventService.Delete(id, userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(baseResponse.ErrorResponse(constant.DATA_NOT_FOUND))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(baseResponse.SuccessResponse(constant.SUCCESS_DELETED, nil))
}
