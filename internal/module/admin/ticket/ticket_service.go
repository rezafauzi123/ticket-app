package auth

import (
	"errors"
	"ticket-app/config"
	"ticket-app/internal/entity"
	"ticket-app/internal/module/admin/ticket/request"
	"ticket-app/internal/repository"
	"ticket-app/pkg/constant"
)

type TicketService interface {
	Create(req request.CreateTicketRequest) (*entity.Ticket, error)
	Update(req request.UpdateTicketRequest, id string) (*entity.Ticket, error)
	FindAll(req request.GetTicketRequest) ([]*entity.Ticket, error)
	FindByID(id string) (*entity.Ticket, error)
}

type ticketService struct {
	ticketRepo repository.TicketRepository
	eventRepo  repository.EventRepository
	config     config.AppConfig
}

func NewTicketService(ticketRepo repository.TicketRepository, eventRepo repository.EventRepository, config config.AppConfig) TicketService {
	return &ticketService{
		ticketRepo: ticketRepo,
		eventRepo:  eventRepo,
		config:     config,
	}
}

func (u *ticketService) Create(req request.CreateTicketRequest) (*entity.Ticket, error) {
	event, err := u.eventRepo.FindByID(req.EventID)
	if err != nil {
		return nil, err
	}

	if event == nil {
		return nil, errors.New(constant.DATA_NOT_FOUND)
	}

	var ticket = &entity.Ticket{
		EventID: req.EventID,
		UserID:  req.UserID,
		Status:  constant.STATUS_PENDING,
	}

	newTicket, err := u.ticketRepo.Create(*ticket)
	if err != nil {
		return nil, err
	}

	event.AvailableTickets -= 1
	_, err = u.eventRepo.Update(req.EventID, *event)
	if err != nil {
		return nil, err
	}

	return newTicket, nil
}

func (u *ticketService) Update(req request.UpdateTicketRequest, id string) (*entity.Ticket, error) {
	existingData, err := u.ticketRepo.FindByID(id)
	if err != nil {
		return nil, errors.New(constant.DATA_NOT_FOUND)
	}

	existingData.Status = req.Status

	updatedData, err := u.ticketRepo.Update(id, *existingData)
	if err != nil {
		return nil, err
	}

	return updatedData, nil
}

func (u *ticketService) FindAll(req request.GetTicketRequest) ([]*entity.Ticket, error) {
	return u.ticketRepo.FindAll(req.EventID, req.UserID, req.Status, req.SortBy, req.SortOrder, req.Pagination)
}

func (u *ticketService) FindByID(id string) (*entity.Ticket, error) {
	data, err := u.ticketRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return data, nil
}
