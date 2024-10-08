package auth

import (
	"errors"
	"ticket-app/config"
	"ticket-app/internal/entity"
	"ticket-app/internal/module/customer/payment/request"
	"ticket-app/internal/repository"
	"ticket-app/pkg/constant"
	"time"
)

type PaymentService interface {
	ProcessPayment(req request.PaymentRequest, userID string) (*entity.Payment, error)
	CancelPayment(id, userID string) (*entity.Payment, error)
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

func (u *paymentService) ProcessPayment(req request.PaymentRequest, userID string) (*entity.Payment, error) {
	ticket, err := u.ticketRepo.FindByID(req.TicketID)
	if err != nil {
		return nil, err
	}

	if ticket == nil {
		return nil, errors.New(constant.DATA_NOT_FOUND)
	}

	if ticket.Status != constant.STATUS_PENDING {
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

	ticket.Status = constant.STATUS_ON_PROCESS
	_, err = u.ticketRepo.Update(req.TicketID, *ticket)
	if err != nil {
		return nil, err
	}

	return newPayment, nil
}

func (u *paymentService) CancelPayment(id, userID string) (*entity.Payment, error) {
	payment, err := u.paymentRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if payment == nil {
		return nil, errors.New(constant.DATA_NOT_FOUND)
	}

	if payment.Status != constant.STATUS_PENDING {
		return nil, errors.New(constant.DATA_NOT_FOUND)
	}

	ticket, err := u.ticketRepo.FindByID(payment.TicketID)
	if err != nil {
		return nil, err
	}

	if ticket == nil {
		return nil, errors.New(constant.DATA_NOT_FOUND)
	}

	if ticket.Status != constant.STATUS_ON_PROCESS {
		return nil, errors.New(constant.DATA_NOT_FOUND)
	}

	payment.Status = constant.STATUS_CANCELED
	newPayment, err := u.paymentRepo.Update(id, *payment)
	if err != nil {
		return nil, err
	}

	ticket.Status = constant.STATUS_CANCELED
	_, err = u.ticketRepo.Update(payment.TicketID, *ticket)
	if err != nil {
		return nil, err
	}

	return newPayment, nil
}
