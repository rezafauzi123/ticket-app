package request

type UpdateProfileRequest struct {
	Name          string `json:"name" validate:"required"`
	Email         string `json:"email" validate:"required,email"`
	Address       string `json:"address"`
	Gender        string `json:"gender" validate:"required"`
	MaritalStatus string `json:"marital_status" validate:"required"`
}
