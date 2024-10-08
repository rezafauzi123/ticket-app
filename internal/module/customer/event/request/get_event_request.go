package request

import "ticket-app/pkg/utils"

type GetEventRequest struct {
	Name             string           `json:"name" validate:"omitempty,min=1,max=100"`
	Location         string           `json:"location" validate:"omitempty,min=1,max=100"`
	Date             string           `json:"date" validate:"omitempty,datetime=2006-01-02"`
	AvailableTickets *int             `json:"available_tickets" validate:"omitempty,min=0"`
	SortBy           string           `json:"sort_by" validate:"omitempty"`
	SortOrder        string           `json:"sort_order" validate:"omitempty,oneof=asc desc"`
	Pagination       utils.Pagination `json:"pagination" validate:"required"`
}
