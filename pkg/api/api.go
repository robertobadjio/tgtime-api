package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"officetime-api/app/config"
	"officetime-api/app/model"
	"officetime-api/app/service"
	"officetime-api/internal/model/router/app"
	"officetime-api/internal/model/router/app/command"
	"officetime-api/internal/model/router/app/command_query"
	"officetime-api/internal/model/router/app/query"
	"officetime-api/internal/model/router/domain/router"
	"time"
)

type apiService struct {
	app app.Application
}

func NewService(app app.Application) Service {
	return &apiService{app: app}
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

func (s *apiService) GetRouters(ctx context.Context) ([]*router.Router, error) {
	qr := query.GetRouters{}
	routers, err := s.app.Queries.GetRouters.Handle(ctx, qr)
	if err != nil {
		return nil, err
	}

	return routers, nil
}

func (s *apiService) GetRouter(ctx context.Context, routerId int) (*router.Router, error) {
	qr := query.GetRouter{RouterId: routerId}
	routerUpdated, err := s.app.Queries.GetRouter.Handle(ctx, qr)
	if err != nil {
		return nil, err
	}

	return routerUpdated, nil
}

func (s *apiService) CreateRouter(ctx context.Context, router *router.Router) (*router.Router, error) {
	cmd := command_query.CreateRouter{Router: router}
	routerNew, err := s.app.Commands.CreateRouter.Handle(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return routerNew, nil
}

func (s *apiService) UpdateRouter(ctx context.Context, routerId int, router *router.Router) (*router.Router, error) {
	if routerId != router.Id {
		return nil, errors.New("error update ids not equals")
	}

	cmd := command.UpdateRouter{Router: router}
	err := s.app.Commands.UpdateRouter.Handle(ctx, cmd)
	if err != nil {
		return nil, err
	}

	qr := query.GetRouter{RouterId: routerId}
	routerNew, err := s.app.Queries.GetRouter.Handle(ctx, qr)
	if err != nil {
		return nil, err
	}

	return routerNew, nil
}

func (s *apiService) DeleteRouter(ctx context.Context, routerId int) error {
	cmd := command.DeleteRouter{RouterId: routerId}
	err := s.app.Commands.DeleteRouter.Handle(ctx, cmd)
	if err != nil {
		return err
	}

	return nil
}
