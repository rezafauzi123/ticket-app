package auth

import (
	"errors"
	"ticket-app/config"
	"ticket-app/internal/entity"
	"ticket-app/internal/module/customer/ticket/request"
	"ticket-app/internal/repository"
	"ticket-app/pkg/constant"
)

type TicketService interface {
	BuyTicket(req request.BuyTicketRequest, userID string) (*entity.Ticket, error)
	CancelTicket(id string) (*entity.Ticket, error)
	FindByUserID(userID string) ([]*entity.Ticket, error)
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

func (u *ticketService) BuyTicket(req request.BuyTicketRequest, userID string) (*entity.Ticket, error) {
	event, err := u.eventRepo.FindByID(req.EventID)
	if err != nil {
		return nil, err
	}

	if event == nil {
		return nil, errors.New(constant.DATA_NOT_FOUND)
	}

	var ticket = &entity.Ticket{
		UserID:  userID,
		EventID: req.EventID,
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

func (u *ticketService) FindByUserID(userID string) ([]*entity.Ticket, error) {
	return u.ticketRepo.FindByUserID(userID)
}

func (u *ticketService) FindByID(id string) (*entity.Ticket, error) {
	data, err := u.ticketRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (u *ticketService) CancelTicket(id string) (*entity.Ticket, error) {
	existingData, err := u.ticketRepo.FindByID(id)
	if err != nil {
		return nil, errors.New(constant.DATA_NOT_FOUND)
	}

	if existingData.Status == constant.STATUS_CANCELED {
		return nil, errors.New(constant.ALREADY_CANCELED)
	}

	event, err := u.eventRepo.FindByID(existingData.EventID)
	if err != nil {
		return nil, err
	}

	if event == nil {
		return nil, errors.New(constant.DATA_NOT_FOUND)
	}

	existingData.Status = constant.STATUS_CANCELED
	updatedData, err := u.ticketRepo.Update(id, *existingData)
	if err != nil {
		return nil, err
	}

	event.AvailableTickets += 1
	_, err = u.eventRepo.Update(existingData.EventID, *event)
	if err != nil {
		return nil, err
	}

	return updatedData, nil
}
