package api

import (
	"context"
	"fmt"
	"net/http"
	"officetime-api/app/config"
	"officetime-api/app/model"
	"officetime-api/app/service"
	"time"
)

type apiService struct{}

func NewService() Service {
	return &apiService{}
}

func (s *apiService) Login(_ context.Context, email, password string) (*service.TokenDetails, error) {
	cfg := config.New()

	td := &service.TokenDetails{}
	td.AccessTokenExpires = time.Now().Add(time.Minute * time.Duration(cfg.AuthAccessTokenExpires)).Unix()
	td.RefreshTokenExpires = time.Now().Add(time.Hour * time.Duration(cfg.AuthRefreshTokenExpires)).Unix()

	user, err := model.GetUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	userPasswordHash := model.GetUserPasswordHashByEmail(user.Email) // TODO: Убрать
	if !service.CheckAuth(userPasswordHash, password) {
		return nil, fmt.Errorf("wrong password")
	}

	return service.CreateTokenPair(user), nil
}

func (s *apiService) ServiceStatus(_ context.Context) (int, error) {
	fmt.Println("Checking the Service health...")
	//logger.Log("Checking the Service health...")
	return http.StatusOK, nil
}

func (s *apiService) GetRouters(_ context.Context) ([]*model.Router, error) {
	return model.GetAllRouters(), nil
}
