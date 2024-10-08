package request

type LoginRequest struct {
	Email    string `form:"email" validate:"required,email"`
	Password string `form:"password,omitempty" validate:"required"`
}
