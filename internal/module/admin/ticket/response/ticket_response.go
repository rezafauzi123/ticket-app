package response

import "ticket-app/internal/entity"

type TicketResponse struct {
	ID      string `json:"id"`
	EventID string `json:"event_id"`
	UserID  string `json:"user_id"`
	Status  string `json:"status"`
}

func (r *TicketResponse) TicketEntityToTicketResponse(data *entity.Ticket) *TicketResponse {
	r.ID = data.ID
	r.EventID = data.EventID
	r.UserID = data.UserID
	r.Status = data.Status
	return r
}
