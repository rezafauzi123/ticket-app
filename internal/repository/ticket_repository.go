package repository

import (
	"ticket-app/config"
	"ticket-app/internal/entity"
	"ticket-app/pkg/utils"
)

type TicketRepository interface {
	Create(ticket entity.Ticket) (*entity.Ticket, error)
	Update(ticketID string, ticket entity.Ticket) (*entity.Ticket, error)
	FindByUserID(userID string) ([]*entity.Ticket, error)
	FindByID(id string) (*entity.Ticket, error)
	FindAll(eventID, userID, status string, sortBy, sortOrder string, pagination utils.Pagination) ([]*entity.Ticket, error)
}

type ticketRepository struct {
	app config.AppConfig
}

func NewTicketRepository(app config.AppConfig) TicketRepository {
	return &ticketRepository{app: app}
}

func (r *ticketRepository) Create(ticket entity.Ticket) (*entity.Ticket, error) {
	err := ticket.BeforeCreate()
	if err != nil {
		return nil, err
	}

	query := `INSERT INTO tickets (id, user_id, event_id, status)
              VALUES ($1, $2, $3, $4)
						RETURNING id, user_id, event_id, status`

	err = r.app.Db.QueryRow(query, ticket.ID, ticket.UserID, ticket.EventID, ticket.Status).
		Scan(&ticket.ID, &ticket.UserID, &ticket.EventID, &ticket.Status)
	return &ticket, err
}

func (r *ticketRepository) Update(ticketID string, ticket entity.Ticket) (*entity.Ticket, error) {
	updatedTicket := entity.Ticket{}
	query := `
			UPDATE tickets
			SET user_id = $1, event_id = $2, status = $3, updated_at = NOW()
			WHERE id = $4 AND deleted_at IS NULL
			RETURNING id, user_id, event_id, status, updated_at`

	err := r.app.Db.QueryRow(query, ticket.UserID, ticket.EventID, ticket.Status, ticketID).
		Scan(&updatedTicket.ID, &updatedTicket.UserID, &updatedTicket.EventID, &updatedTicket.Status, &updatedTicket.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &updatedTicket, nil
}

func (r *ticketRepository) FindByUserID(userID string) ([]*entity.Ticket, error) {
	var tickets []*entity.Ticket
	query := `SELECT * FROM tickets WHERE user_id = $1 AND deleted_at IS NULL`
	err := r.app.Db.Select(&tickets, query, userID)
	if err != nil {
		return nil, err
	}
	return tickets, nil
}

func (r *ticketRepository) FindByID(id string) (*entity.Ticket, error) {
	var ticket entity.Ticket
	query := `SELECT * FROM tickets WHERE id = $1 AND deleted_at IS NULL`
	err := r.app.Db.Get(&ticket, query, id)
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *ticketRepository) FindAll(eventID, userID, status string, sortBy, sortOrder string, pagination utils.Pagination) ([]*entity.Ticket, error) {
	var tickets []*entity.Ticket

	query := `SELECT * FROM tickets WHERE deleted_at IS NULL`
	args := []interface{}{}
	idx := 1

	query, args = utils.ApplyFilters(query, args, map[string]interface{}{
		"event_id": eventID,
		"user_id":  userID,
		"status":   status,
	}, &idx)

	// Sorting
	query = utils.ApplySorting(query, sortBy, sortOrder)

	// Pagination
	var paginationArgs []interface{}
	query, paginationArgs = utils.ApplyPagination(query, pagination, &idx)

	args = append(args, paginationArgs...)

	err := r.app.Db.Select(&tickets, query, args...)
	if err != nil {
		return nil, err
	}

	return tickets, nil
}
