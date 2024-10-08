package response

import "ticket-app/internal/entity"

type TicketResponse struct {
	ID      string `json:"id"`
	EventID string `json:"event_id"`
	Status  string `json:"status"`
}

func (r *TicketResponse) TicketEntityToTicketResponse(data *entity.Ticket) *TicketResponse {
	r.ID = data.ID
	r.EventID = data.EventID
	r.Status = data.Status
	return r
}
