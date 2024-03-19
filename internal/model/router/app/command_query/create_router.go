package command_query

import (
	"context"
	"officetime-api/internal/common/decorator"
	"officetime-api/internal/model/router/domain/router"
)

type CreateRouter struct {
	Router *router.Router
}

type CreateRouterHandler decorator.CommandQueryHandler[CreateRouter, *router.Router]

type createRouterHandler struct {
	routerRepository router.Repository
}

func NewCreateRouterHandler(routerRepository router.Repository) CreateRouterHandler {
	if routerRepository == nil {
		panic("nil routerRepository")
	}

	return decorator.ApplyCommandQueryDecorators[CreateRouter, *router.Router](
		createRouterHandler{routerRepository: routerRepository},
	)
}

func (h createRouterHandler) Handle(ctx context.Context, cmdQr CreateRouter) (*router.Router, error) {
	routerNew, err := h.routerRepository.CreateRouter(ctx, cmdQr.Router)
	if err != nil {
		return nil, err
	}

	return routerNew, nil
}
