package api

import (
	"context"
	"fmt"
	"net/http"
	"officetime-api/app/service"
	"officetime-api/internal/config"
	"officetime-api/internal/db"
	departmentApp "officetime-api/internal/model/department/app"
	periodApp "officetime-api/internal/model/period/app"
	routerApp "officetime-api/internal/model/router/app"
	"officetime-api/internal/model/user/adapter"
	userApp "officetime-api/internal/model/user/app"
	"officetime-api/internal/model/user/app/query"
	weekendApp "officetime-api/internal/model/weekend/app"
	"time"
)

type apiService struct {
	routerApp     routerApp.Application
	periodApp     periodApp.Application
	departmentApp departmentApp.Application
	weekendApp    weekendApp.Application
	userApp       userApp.Application
}

func NewService(
	routerApp routerApp.Application,
	periodApp periodApp.Application,
	departmentApp departmentApp.Application,
	weekendApp weekendApp.Application,
	userApp userApp.Application,
) Service {
	return &apiService{
		routerApp:     routerApp,
		periodApp:     periodApp,
		departmentApp: departmentApp,
		weekendApp:    weekendApp,
		userApp:       userApp,
	}
}

func (s *apiService) Login(_ context.Context, email, password string) (*service.TokenDetails, error) {
	cfg := config.New()

	td := &service.TokenDetails{}
	td.AccessTokenExpires = time.Now().Add(time.Minute * time.Duration(cfg.AuthAccessTokenExpires)).Unix()
	td.RefreshTokenExpires = time.Now().Add(time.Hour * time.Duration(cfg.AuthRefreshTokenExpires)).Unix()

	uApp := buildUserApp()
	qr := query.GetUserByEmail{Email: email}
	ctx := context.TODO()
	user, err := uApp.Queries.GetUserByEmail.Handle(ctx, qr)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	qr2 := query.GetUserPasswordHashByEmail{Email: email}
	// TODO: Убрать
	userPasswordHash, _ := uApp.Queries.GetUserPasswordHashByEmail.Handle(ctx, qr2) // TODO: Handle error
	if !service.CheckAuth(userPasswordHash, password) {
		return nil, fmt.Errorf("wrong password")
	}

	return service.CreateTokenPair(user), nil
}

func buildUserApp() userApp.Application {
	userRepository := adapter.NewPgUserRepository(db.GetDB())
	return userApp.Application{
		Queries: userApp.Queries{
			GetUser:                    query.NewGetUserHandler(userRepository),
			GetUserByEmail:             query.NewGetUserByEmailHandler(userRepository),
			GetUserPasswordHashByEmail: query.NewGetUserPasswordHashByEmailHandler(userRepository),
		},
	}
}

func (s *apiService) ServiceStatus(_ context.Context) (int, error) {
	fmt.Println("Checking the Service health...")
	//logger.Log("Checking the Service health...")
	return http.StatusOK, nil
}
