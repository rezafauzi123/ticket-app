package request

import "ticket-app/pkg/utils"

type GetTicketRequest struct {
	EventID    string           `json:"event_id" validate:"omitempty,min=1,max=100"`
	UserID     string           `json:"user_id" validate:"omitempty,min=1,max=100"`
	Status     string           `json:"status" validate:"omitempty,min=1,max=100"`
	SortBy     string           `json:"sort_by" validate:"omitempty"`
	SortOrder  string           `json:"sort_order" validate:"omitempty,oneof=asc desc"`
	Pagination utils.Pagination `json:"pagination" validate:"required"`
}
