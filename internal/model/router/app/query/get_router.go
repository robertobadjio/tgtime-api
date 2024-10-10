package query

import (
	"context"
	"github.com/robertobadjio/tgtime-api/internal/common/decorator"
	"github.com/robertobadjio/tgtime-api/internal/model/router/domain/router"
)

type GetRouter struct {
	RouterId int
}

type GetRouterHandler decorator.QueryHandler[GetRouter, *router.Router]

/*type GetRouterReadModel interface {
	GetRouter(ctx context.Context, routerId int) (*router.Router, error)
}*/

type getRouterHandler struct {
	routerRepo router.Repository
}

func NewGetRouterHandler(routerRepo router.Repository) GetRouterHandler {
	return decorator.ApplyQueryDecorators[GetRouter, *router.Router](
		getRouterHandler{routerRepo: routerRepo},
	)
}

func (h getRouterHandler) Handle(ctx context.Context, qr GetRouter) (*router.Router, error) {
	return h.routerRepo.GetRouter(ctx, qr.RouterId)
}
