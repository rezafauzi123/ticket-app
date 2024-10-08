package request

type CreateTicketRequest struct {
	EventID string `json:"event_id" validate:"required"`
	UserID  string `json:"user_id" validate:"required"`
}
