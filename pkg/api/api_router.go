package api

import (
	"context"
	"errors"
	"github.com/robertobadjio/tgtime-api/internal/model/router/app/command"
	"github.com/robertobadjio/tgtime-api/internal/model/router/app/command_query"
	"github.com/robertobadjio/tgtime-api/internal/model/router/app/query"
	"github.com/robertobadjio/tgtime-api/internal/model/router/domain/router"
)

func (s *apiService) GetRouters(ctx context.Context) ([]*router.Router, error) {
	qr := query.GetRouters{}
	routers, err := s.routerApp.Queries.GetRouters.Handle(ctx, qr)
	if err != nil {
		return nil, err
	}

	return routers, nil
}

func (s *apiService) GetRouter(ctx context.Context, routerId int) (*router.Router, error) {
	qr := query.GetRouter{RouterId: routerId}
	routerUpdated, err := s.routerApp.Queries.GetRouter.Handle(ctx, qr)
	if err != nil {
		return nil, err
	}

	return routerUpdated, nil
}

func (s *apiService) CreateRouter(ctx context.Context, router *router.Router) (*router.Router, error) {
	cmd := command_query.CreateRouter{Router: router}
	routerNew, err := s.routerApp.CommandsQueries.CreateRouter.Handle(ctx, cmd)
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
	err := s.routerApp.Commands.UpdateRouter.Handle(ctx, cmd)
	if err != nil {
		return nil, err
	}

	qr := query.GetRouter{RouterId: routerId}
	routerNew, err := s.routerApp.Queries.GetRouter.Handle(ctx, qr)
	if err != nil {
		return nil, err
	}

	return routerNew, nil
}

func (s *apiService) DeleteRouter(ctx context.Context, routerId int) error {
	cmd := command.DeleteRouter{RouterId: routerId}
	err := s.routerApp.Commands.DeleteRouter.Handle(ctx, cmd)
	if err != nil {
		return err
	}

	return nil
}
