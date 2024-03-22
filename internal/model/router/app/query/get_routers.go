package query

import (
	"context"
	"officetime-api/internal/common/decorator"
	"officetime-api/internal/model/router/domain/router"
)

type GetRouters struct {
}

type GetRoutersHandler decorator.QueryHandler[GetRouters, []*router.Router]

/*type GetRoutersReadModel interface {
	GetRouters(ctx context.Context) ([]*router.Router, error)
}*/

type getRoutersHandler struct {
	routerRepo router.Repository
}

func NewGetRoutersHandler(routerRepo router.Repository) GetRoutersHandler {
	return decorator.ApplyQueryDecorators[GetRouters, []*router.Router](
		getRoutersHandler{routerRepo: routerRepo},
	)
}

func (h getRoutersHandler) Handle(ctx context.Context, _ GetRouters) ([]*router.Router, error) {
	return h.routerRepo.GetRouters(ctx)
}
