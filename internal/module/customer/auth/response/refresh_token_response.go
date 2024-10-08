package response

type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
}

func (r *RefreshTokenResponse) UserentityToRefreshTokenResponse(token string) *RefreshTokenResponse {
	r.AccessToken = token

	return r
}
