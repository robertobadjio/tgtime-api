package api

import (
	"context"
	"fmt"
	"net/http"
	"officetime-api/app/model"
	"officetime-api/app/service"
	"officetime-api/internal/config"
	departmentApp "officetime-api/internal/model/department/app"
	periodApp "officetime-api/internal/model/period/app"
	routerApp "officetime-api/internal/model/router/app"
	weekendApp "officetime-api/internal/model/weekend/app"
	"time"
)

type apiService struct {
	routerApp     routerApp.Application
	periodApp     periodApp.Application
	departmentApp departmentApp.Application
	weekendApp    weekendApp.Application
}

func NewService(
	routerApp routerApp.Application,
	periodApp periodApp.Application,
	departmentApp departmentApp.Application,
	weekendApp weekendApp.Application,
) Service {
	return &apiService{
		routerApp:     routerApp,
		periodApp:     periodApp,
		departmentApp: departmentApp,
		weekendApp:    weekendApp,
	}
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
