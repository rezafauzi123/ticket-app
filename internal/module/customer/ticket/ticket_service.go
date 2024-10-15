package auth

import (
	"errors"
	"ticket-app/config"
	"ticket-app/internal/entity"
	event "ticket-app/internal/module/customer/event"
	"ticket-app/internal/module/customer/ticket/request"
	"ticket-app/internal/repository"
	"ticket-app/pkg/constant"
)

type TicketService interface {
	BuyTicket(req request.BuyTicketRequest, userID string) (*entity.Ticket, error)
	CancelTicket(id string) (*entity.Ticket, error)
	FindByUserID(userID string) ([]*entity.Ticket, error)
	FindByID(id string) (*entity.Ticket, error)
	UpdateStatus(ticketID string, status string) error
}

type ticketService struct {
	ticketRepo   repository.TicketRepository
	eventService event.EventService
	config       config.AppConfig
}

func NewTicketService(ticketRepo repository.TicketRepository, eventService event.EventService, config config.AppConfig) TicketService {
	return &ticketService{
		ticketRepo:   ticketRepo,
		eventService: eventService,
		config:       config,
	}
}

func (u *ticketService) BuyTicket(req request.BuyTicketRequest, userID string) (*entity.Ticket, error) {
	var ticket = &entity.Ticket{
		UserID:  userID,
		EventID: req.EventID,
		Status:  constant.STATUS_PENDING,
	}

	newTicket, err := u.ticketRepo.Create(*ticket)
	if err != nil {
		return nil, err
	}

	_, err = u.eventService.DecreaseTicket(req.EventID)
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

	existingData.Status = constant.STATUS_CANCELED
	updatedData, err := u.ticketRepo.Update(id, *existingData)
	if err != nil {
		return nil, err
	}

	_, err = u.eventService.IncreaseTicket(existingData.EventID)
	if err != nil {
		return nil, err
	}
	return updatedData, nil
}

func (u *ticketService) UpdateStatus(ticketID string, status string) error {
	ticket, err := u.ticketRepo.FindByID(ticketID)
	if err != nil || ticket == nil {
		return errors.New(constant.DATA_NOT_FOUND)
	}

	ticket.Status = status
	_, err = u.ticketRepo.Update(ticketID, *ticket)
	return err
}
