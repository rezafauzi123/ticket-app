package repository

import (
	"ticket-app/config"
	"ticket-app/internal/entity"
	"ticket-app/pkg/utils"
	"time"
)

type EventRepository interface {
	Create(event entity.Event) (*entity.Event, error)
	FindByID(id string) (*entity.Event, error)
	FindAll(name, location, date string, availableTickets *int, sortBy, sortOrder string, pagination utils.Pagination) ([]*entity.Event, error)
	Update(eventID string, event entity.Event) (*entity.Event, error)
	Delete(eventID string, deletedBy string) error
}

type eventRepository struct {
	app config.AppConfig
}

func NewEventRepository(app config.AppConfig) EventRepository {
	return &eventRepository{app: app}
}

func (r *eventRepository) FindByID(id string) (*entity.Event, error) {
	var event entity.Event
	query := `SELECT * FROM events WHERE id = $1 AND deleted_at IS NULL`
	err := r.app.Db.Get(&event, query, id)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *eventRepository) Create(event entity.Event) (*entity.Event, error) {
	err := event.BeforeCreate()
	if err != nil {
		return nil, err
	}

	query := `INSERT INTO events (id, name, date, location, available_tickets, description)
              VALUES ($1, $2, $3, $4, $5, $6)
						RETURNING id, name, date, location, available_tickets, description`

	err = r.app.Db.QueryRow(query, event.ID, event.Name, event.Date, event.Location, event.AvailableTickets, event.Description).
		Scan(&event.ID, &event.Name, &event.Date, &event.Location, &event.AvailableTickets, &event.Description)
	return &event, err
}

func (r *eventRepository) FindAll(name, location, date string, availableTickets *int, sortBy, sortOrder string, pagination utils.Pagination) ([]*entity.Event, error) {
	var events []*entity.Event

	query := `SELECT * FROM events WHERE deleted_at IS NULL`
	args := []interface{}{}
	idx := 1

	query, args = utils.ApplyFilters(query, args, map[string]interface{}{
		"name":              name,
		"location":          location,
		"date":              date,
		"available_tickets": availableTickets,
	}, &idx)

	// Sorting
	query = utils.ApplySorting(query, sortBy, sortOrder)

	// Pagination
	var paginationArgs []interface{}
	query, paginationArgs = utils.ApplyPagination(query, pagination, &idx)

	args = append(args, paginationArgs...)

	err := r.app.Db.Select(&events, query, args...)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (r *eventRepository) Update(eventID string, event entity.Event) (*entity.Event, error) {
	updatedEvent := entity.Event{}
	query := `
			UPDATE events
			SET name = $1, date = $2, location = $3, available_tickets = $4, description = $5, updated_at = NOW()
			WHERE id = $6 AND deleted_at IS NULL
			RETURNING id, name, date, location, available_tickets, description, updated_at
	`
	err := r.app.Db.QueryRow(query, event.Name, event.Date, event.Location, event.AvailableTickets, event.Description, eventID).
		Scan(&updatedEvent.ID, &updatedEvent.Name, &updatedEvent.Date, &updatedEvent.Location, &updatedEvent.AvailableTickets, &updatedEvent.Description, &updatedEvent.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &updatedEvent, nil
}

func (r *eventRepository) Delete(eventID string, deletedBy string) error {
	query := `
			UPDATE events
			SET deleted_at = $1, deleted_by = $2
			WHERE id = $3 AND deleted_at IS NULL
	`
	_, err := r.app.Db.Exec(query, time.Now(), deletedBy, eventID)
	if err != nil {
		return err
	}

	return nil
}
