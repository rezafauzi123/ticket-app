package request

type BuyTicketRequest struct {
	EventID string `json:"event_id" validate:"required"`
}
