package response

import "ticket-app/internal/entity"

type RegisterResponse struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (r *RegisterResponse) UserentityToRegisterResponse(data *entity.User, token *map[string]string) *RegisterResponse {
	var resulToken map[string]string

	if token != nil {
		resulToken = *token
	}

	r.Name = data.Name
	r.Email = data.Email
	r.AccessToken = resulToken["access_token"]
	r.RefreshToken = resulToken["refresh_token"]

	return r
}