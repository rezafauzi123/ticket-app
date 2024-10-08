package request

import "time"

type UpdateEventRequest struct {
	Name             string    `json:"name" validate:"required"`
	Location         string    `json:"location" validate:"required"`
	Date             time.Time `json:"date" validate:"required"`
	Description      string    `json:"description"`
	AvailableTickets int       `json:"available_tickets" validate:"required"`
}
