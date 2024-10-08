package repository

import (
	"ticket-app/config"
	"ticket-app/internal/entity"
	"ticket-app/pkg/utils"
)

type PaymentRepository interface {
	Create(payment entity.Payment) (*entity.Payment, error)
	Update(id string, payment entity.Payment) (*entity.Payment, error)
	FindByID(id string) (*entity.Payment, error)
	FindAll(ticketID, userID, status string, sortBy, sortOrder string, pagination utils.Pagination) ([]*entity.Payment, error)
}

type paymentRepository struct {
	app config.AppConfig
}

func NewPaymentRepository(app config.AppConfig) PaymentRepository {
	return &paymentRepository{app: app}
}

func (r *paymentRepository) Create(payment entity.Payment) (*entity.Payment, error) {
	err := payment.BeforeCreate()
	if err != nil {
		return nil, err
	}

	query := `INSERT INTO payments (id, user_id, ticket_id, amount, status, payment_date)
              VALUES ($1, $2, $3, $4, $5, $6)
						RETURNING id, user_id, ticket_id, amount, status, payment_date`

	err = r.app.Db.QueryRow(query, payment.ID, payment.UserID, payment.TicketID, payment.Amount, payment.Status, payment.PaymentDate).
		Scan(&payment.ID, &payment.UserID, &payment.TicketID, &payment.Amount, &payment.Status, &payment.PaymentDate)
	return &payment, err
}

func (r *paymentRepository) Update(id string, payment entity.Payment) (*entity.Payment, error) {
	updatedPayment := entity.Payment{}
	query := `
			UPDATE payments
			SET user_id = $1, ticket_id = $2, amount = $3, status = $4, payment_date = $5, updated_at = NOW()
			WHERE id = $6 AND deleted_at IS NULL
			RETURNING id, user_id, ticket_id, amount, status, payment_date, updated_at`

	err := r.app.Db.QueryRow(query, payment.UserID, payment.TicketID, payment.Amount, payment.Status, payment.PaymentDate, id).
		Scan(&updatedPayment.ID, &updatedPayment.UserID, &updatedPayment.TicketID, &updatedPayment.Amount, &updatedPayment.Status, &updatedPayment.PaymentDate, &updatedPayment.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &updatedPayment, nil
}

func (r *paymentRepository) FindByID(id string) (*entity.Payment, error) {
	var payment entity.Payment
	query := `SELECT * FROM payments WHERE id = $1 AND deleted_at IS NULL`
	err := r.app.Db.Get(&payment, query, id)
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) FindAll(ticketID, userID, status string, sortBy, sortOrder string, pagination utils.Pagination) ([]*entity.Payment, error) {
	var payments []*entity.Payment

	query := `SELECT * FROM payments WHERE deleted_at IS NULL`
	args := []interface{}{}
	idx := 1

	query, args = utils.ApplyFilters(query, args, map[string]interface{}{
		"ticket_id": ticketID,
		"user_id":   userID,
		"status":    status,
	}, &idx)

	query = utils.ApplySorting(query, sortBy, sortOrder)

	var paginationArgs []interface{}
	query, paginationArgs = utils.ApplyPagination(query, pagination, &idx)

	args = append(args, paginationArgs...)

	err := r.app.Db.Select(&payments, query, args...)
	if err != nil {
		return nil, err
	}

	return payments, nil
}
