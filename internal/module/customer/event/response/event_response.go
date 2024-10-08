package response

import "ticket-app/internal/entity"

type EventResponse struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Date             string `json:"date"`
	Location         string `json:"location"`
	AvailableTickets int    `json:"available_tickets"`
	Description      string `json:"description"`
}

func (r *EventResponse) EventEntityToEventResponse(data *entity.Event) *EventResponse {
	r.ID = data.ID
	r.Name = data.Name
	r.Date = data.Date.String()
	r.Location = data.Location
	r.AvailableTickets = data.AvailableTickets
	r.Description = data.Description
	return r
}
