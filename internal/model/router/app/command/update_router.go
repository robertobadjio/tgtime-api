package command

import (
	"context"
	"github.com/robertobadjio/tgtime-api/internal/common/decorator"
	"github.com/robertobadjio/tgtime-api/internal/model/router/domain/router"
)

type UpdateRouter struct {
	Router *router.Router
}

type UpdateRouterHandler decorator.CommandHandler[UpdateRouter]

type updateRouterHandler struct {
	routerRepository router.Repository
}

func NewUpdateRouterHandler(routerRepository router.Repository) UpdateRouterHandler {
	if routerRepository == nil {
		panic("nil routerRepository")
	}

	return decorator.ApplyCommandDecorators[UpdateRouter](
		updateRouterHandler{routerRepository: routerRepository},
	)
}

func (h updateRouterHandler) Handle(ctx context.Context, cmd UpdateRouter) error {
	_, err := h.routerRepository.UpdateRouter(ctx, cmd.Router) // TODO: !
	if err != nil {
		return err
	}

	return nil
}
