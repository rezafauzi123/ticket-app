package request

type PaymentRequest struct {
	TicketID string  `json:"ticket_id" validate:"required"`
	Amount   float64 `json:"amount" validate:"required"`
}
