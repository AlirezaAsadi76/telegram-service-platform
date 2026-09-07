package userparams

type LoginRequest struct {
	InitData string `json:"init_data"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
