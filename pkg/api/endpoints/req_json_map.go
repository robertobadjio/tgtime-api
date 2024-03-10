package endpoints

import "officetime-api/app/model"

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken         string `json:"access_token"`
	RefreshToken        string `json:"refresh_token"`
	AccessTokenExpires  int64  `json:"access_token_expires"`
	RefreshTokenExpires int64  `json:"refresh_token_expires"`
}

type ServiceStatusRequest struct{}

type ServiceStatusResponse struct {
	Code int    `json:"status"`
	Err  string `json:"err,omitempty"`
}

type GetRoutersRequest struct{}

type GetRoutersResponse struct {
	Routers []*model.Router `json:"routers"`
}
