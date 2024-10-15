package auth

import (
	"errors"
	"ticket-app/config"
	"ticket-app/internal/entity"
	"ticket-app/internal/module/customer/event/request"
	"ticket-app/internal/repository"
	"ticket-app/pkg/constant"
)

type EventService interface {
	FindAll(req request.GetEventRequest) ([]*entity.Event, error)
	FindByID(id string) (*entity.Event, error)
	DecreaseTicket(eventID string) (*entity.Event, error)
	IncreaseTicket(eventID string) (*entity.Event, error)
}

type eventService struct {
	eventRepo repository.EventRepository
	config    config.AppConfig
}

func NewEventService(eventRepo repository.EventRepository, config config.AppConfig) EventService {
	return &eventService{
		eventRepo: eventRepo,
		config:    config,
	}
}

func (u *eventService) FindAll(req request.GetEventRequest) ([]*entity.Event, error) {
	return u.eventRepo.FindAll(req.Name, req.Location, req.Date, req.AvailableTickets, req.SortBy, req.SortOrder, req.Pagination)
}

func (u *eventService) FindByID(id string) (*entity.Event, error) {
	data, err := u.eventRepo.FindByID(id)
	if err != nil {
		return nil, errors.New(constant.DATA_NOT_FOUND)
	}

	return data, nil
}

func (e *eventService) DecreaseTicket(eventID string) (*entity.Event, error) {
	event, err := e.eventRepo.FindByID(eventID)
	if err != nil || event == nil {
		return nil, errors.New(constant.DATA_NOT_FOUND)
	}

	if event.AvailableTickets <= 0 {
		return nil, errors.New(constant.DATA_NOT_FOUND)
	}

	event.AvailableTickets -= 1
	updatedEvent, err := e.eventRepo.Update(eventID, *event)
	if err != nil {
		return nil, err
	}

	return updatedEvent, nil
}

func (e *eventService) IncreaseTicket(eventID string) (*entity.Event, error) {
	event, err := e.eventRepo.FindByID(eventID)
	if err != nil || event == nil {
		return nil, errors.New(constant.DATA_NOT_FOUND)
	}

	event.AvailableTickets += 1
	updatedEvent, err := e.eventRepo.Update(eventID, *event)
	if err != nil {
		return nil, err
	}

	return updatedEvent, nil
}
