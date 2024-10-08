package auth

import (
	"ticket-app/config"
	"ticket-app/internal/entity"
	"ticket-app/internal/module/customer/event/request"
	"ticket-app/internal/repository"
)

type EventService interface {
	FindAll(req request.GetEventRequest) ([]*entity.Event, error)
	FindByID(id string) (*entity.Event, error)
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
		return nil, err
	}

	return data, nil
}
