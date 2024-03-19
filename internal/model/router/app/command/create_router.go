package command

import (
	"context"
	"officetime-api/internal/common/decorator"
	"officetime-api/internal/model/router/domain/router"
)

type CreateRouter struct {
	Router *router.Router
}

type CreateRouterHandler decorator.CommandHandler[CreateRouter]

type createRouterHandler struct {
	routerRepository router.Repository
}

func NewCreateRouterHandler(routerRepository router.Repository) CreateRouterHandler {
	if routerRepository == nil {
		panic("nil routerRepository")
	}

	return decorator.ApplyCommandDecorators[CreateRouter](
		createRouterHandler{routerRepository: routerRepository},
	)
}

func (h createRouterHandler) Handle(ctx context.Context, cmd CreateRouter) error {
	_, err := h.routerRepository.CreateRouter(ctx, cmd.Router) // TODO: !
	if err != nil {
		return err
	}

	return nil
}
