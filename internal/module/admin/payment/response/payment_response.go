package response

import "ticket-app/internal/entity"

type PaymentResponse struct {
	ID          string  `json:"id"`
	TicketID    string  `json:"ticket_id"`
	Amount      float64 `json:"amount"`
	Status      string  `json:"status"`
	PaymentDate string  `json:"payment_date"`
}

func (r *PaymentResponse) PaymentEntityToPaymentResponse(data *entity.Payment) *PaymentResponse {
	r.ID = data.ID
	r.TicketID = data.TicketID
	r.Amount = data.Amount
	r.Status = data.Status
	r.PaymentDate = data.PaymentDate.String()
	return r
}
