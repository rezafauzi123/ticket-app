package request

type UpdateTicketRequest struct {
	Status string `json:"status" validate:"required"`
}
