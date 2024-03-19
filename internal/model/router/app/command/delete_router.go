package command

import (
	"context"
	"officetime-api/internal/common/decorator"
	"officetime-api/internal/model/router/domain/router"
)

type DeleteRouter struct {
	RouterId int
}

type DeleteRouterHandler decorator.CommandHandler[DeleteRouter]

type deleteRouterHandler struct {
	routerRepository router.Repository
}

func NewDeleteRouterHandler(routerRepository router.Repository) DeleteRouterHandler {
	if routerRepository == nil {
		panic("nil routerRepository")
	}

	return decorator.ApplyCommandDecorators[DeleteRouter](
		deleteRouterHandler{routerRepository: routerRepository},
	)
}

func (h deleteRouterHandler) Handle(ctx context.Context, cmd DeleteRouter) error {
	err := h.routerRepository.DeleteRouter(ctx, cmd.RouterId)
	if err != nil {
		return err
	}

	return nil
}
