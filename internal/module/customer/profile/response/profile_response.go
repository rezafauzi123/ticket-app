package response

import "ticket-app/internal/entity"

type UserResponse struct {
	Name          string `json:"name"`
	Email         string `json:"email"`
	Address       string `json:"address"`
	Gender        string `json:"gender"`
	MaritalStatus string `json:"marital_status"`
}

func (r *UserResponse) UserentityToUserResponse(data *entity.User) *UserResponse {
	r.Name = data.Name
	r.Email = data.Email
	r.Address = data.Address
	r.Gender = data.Gender
	r.MaritalStatus = data.MaritalStatus

	return r
}
