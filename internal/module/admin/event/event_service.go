package auth

import (
	"errors"
	"ticket-app/config"
	"ticket-app/internal/entity"
	"ticket-app/internal/module/admin/event/request"
	"ticket-app/internal/repository"
	"ticket-app/pkg/constant"
)

type EventService interface {
	Create(req request.CreateEventRequest) (*entity.Event, error)
	Update(req request.UpdateEventRequest, id string) (*entity.Event, error)
	FindAll(req request.GetEventRequest) ([]*entity.Event, error)
	FindByID(id string) (*entity.Event, error)
	Delete(id, userID string) error
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

func (u *eventService) Create(req request.CreateEventRequest) (*entity.Event, error) {
	var event = &entity.Event{
		Name:             req.Name,
		Date:             req.Date,
		Location:         req.Location,
		AvailableTickets: req.AvailableTickets,
		Description:      req.Description,
	}

	newEvent, err := u.eventRepo.Create(*event)
	if err != nil {
		return nil, err
	}

	return newEvent, nil
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

func (u *eventService) Update(req request.UpdateEventRequest, id string) (*entity.Event, error) {
	existingData, err := u.eventRepo.FindByID(id)
	if err != nil {
		return nil, errors.New(constant.DATA_NOT_FOUND)
	}

	existingData.Name = req.Name
	existingData.Date = req.Date
	existingData.Location = req.Location
	existingData.AvailableTickets = req.AvailableTickets
	existingData.Description = req.Description

	updatedData, err := u.eventRepo.Update(id, *existingData)
	if err != nil {
		return nil, err
	}

	return updatedData, nil
}

func (u *eventService) Delete(id, userID string) error {
	err := u.eventRepo.Delete(id, userID)
	if err != nil {
		return errors.New(constant.DATA_NOT_FOUND)
	}

	return nil
}
