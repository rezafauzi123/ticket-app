package auth

import (
	"errors"
	"ticket-app/config"
	"ticket-app/internal/entity"
	"ticket-app/internal/module/admin/payment/request"
	"ticket-app/internal/repository"
	"ticket-app/pkg/constant"
	"time"
)

type PaymentService interface {
	ConfirmPayment(id string, req request.PaymentRequest) (*entity.Payment, error)
	FindAll(req request.GetPaymentRequest) ([]*entity.Payment, error)
	FindByID(id string) (*entity.Payment, error)
}

type paymentService struct {
	paymentRepo repository.PaymentRepository
	ticketRepo  repository.TicketRepository
	config      config.AppConfig
}

func NewPaymentService(paymentRepo repository.PaymentRepository, ticketRepository repository.TicketRepository, config config.AppConfig) PaymentService {
	return &paymentService{
		paymentRepo: paymentRepo,
		ticketRepo:  ticketRepository,
		config:      config,
	}
}

func (u *paymentService) ConfirmPayment(id string, req request.PaymentRequest) (*entity.Payment, error) {
	existingData, err := u.paymentRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if existingData == nil {
		return nil, errors.New(constant.DATA_NOT_FOUND)
	}

	if existingData.Status == constant.STATUS_CANCELED ||
		existingData.Status == constant.STATUS_FAILED {
		return nil, errors.New(constant.DATA_NOT_FOUND)
	}

	existingTicket, err := u.ticketRepo.FindByID(existingData.TicketID)
	if err != nil {
		return nil, err
	}

	if existingTicket == nil {
		return nil, errors.New(constant.DATA_NOT_FOUND)
	}

	currentDate := time.Now()
	existingData.Status = req.Status
	existingData.UpdatedAt = &currentDate
	updatedPayment, err := u.paymentRepo.Update(id, *existingData)
	if err != nil {
		return nil, err
	}

	existingTicket.Status = req.Status
	existingTicket.UpdatedAt = &currentDate
	_, err = u.ticketRepo.Update(existingData.TicketID, *existingTicket)
	if err != nil {
		return nil, err
	}

	return updatedPayment, nil
}

func (u *paymentService) FindAll(req request.GetPaymentRequest) ([]*entity.Payment, error) {
	return u.paymentRepo.FindAll(req.TicketID, req.UserID, req.Status, req.SortBy, req.SortOrder, req.Pagination)
}

func (u *paymentService) FindByID(id string) (*entity.Payment, error) {
	data, err := u.paymentRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return data, nil
}
