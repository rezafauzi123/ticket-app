package auth

import (
	"errors"
	"ticket-app/config"
	"ticket-app/internal/entity"
	"ticket-app/internal/module/customer/payment/request"
	"ticket-app/internal/repository"
	"ticket-app/pkg/constant"
	"time"

	ticket "ticket-app/internal/module/customer/ticket"
)

type PaymentService interface {
	ProcessPayment(req request.PaymentRequest, userID string) (*entity.Payment, error)
	CancelPayment(id, userID string) (*entity.Payment, error)
}

type paymentService struct {
	paymentRepo   repository.PaymentRepository
	ticketService ticket.TicketService
	config        config.AppConfig
}

func NewPaymentService(paymentRepo repository.PaymentRepository, ticketService ticket.TicketService, config config.AppConfig) PaymentService {
	return &paymentService{
		paymentRepo:   paymentRepo,
		ticketService: ticketService,
		config:        config,
	}
}

func (u *paymentService) ProcessPayment(req request.PaymentRequest, userID string) (*entity.Payment, error) {
	ticket, err := u.ticketService.FindByID(req.TicketID)
	if err != nil {
		return nil, err
	}

	if ticket == nil || ticket.Status != constant.STATUS_PENDING {
		return nil, errors.New(constant.DATA_NOT_FOUND)
	}

	currentDate := time.Now()
	var payment = &entity.Payment{
		UserID:      userID,
		TicketID:    req.TicketID,
		Amount:      req.Amount,
		Status:      constant.STATUS_PENDING,
		PaymentDate: currentDate,
	}

	newPayment, err := u.paymentRepo.Create(*payment)
	if err != nil {
		return nil, err
	}

	err = u.ticketService.UpdateStatus(req.TicketID, constant.STATUS_ON_PROCESS)
	if err != nil {
		return nil, err
	}

	return newPayment, nil
}

func (u *paymentService) CancelPayment(id, userID string) (*entity.Payment, error) {
	payment, err := u.paymentRepo.FindByID(id)
	if err != nil || payment == nil || payment.Status != constant.STATUS_PENDING {
		return nil, errors.New(constant.DATA_NOT_FOUND)
	}

	payment.Status = constant.STATUS_CANCELED
	newPayment, err := u.paymentRepo.Update(id, *payment)
	if err != nil {
		return nil, err
	}

	err = u.ticketService.UpdateStatus(payment.TicketID, constant.STATUS_CANCELED)
	if err != nil {
		return nil, err
	}

	return newPayment, nil
}
