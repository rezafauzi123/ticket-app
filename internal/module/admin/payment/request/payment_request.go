package request

type PaymentRequest struct {
	Status string `json:"status" validate:"required"`
}
