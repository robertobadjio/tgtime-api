package query

import (
	"context"
	"officetime-api/internal/common/decorator"
	"officetime-api/internal/model/weekend/domain/weekend"
)

type GetWeekends struct{}

type GetWeekendsHandler decorator.QueryHandler[GetWeekends, []string]

/*type GetWeekendsReadModel interface {
	GetWeekends(ctx context.Context) ([]*weekend.Weekend, error)
}*/

type getWeekendsHandler struct {
	weekendRepo weekend.Repository
}

func NewGetWeekendsHandler(weekendRepo weekend.Repository) GetWeekendsHandler {
	return decorator.ApplyQueryDecorators[GetWeekends, []string](
		getWeekendsHandler{weekendRepo: weekendRepo},
	)
}

func (h getWeekendsHandler) Handle(ctx context.Context, _ GetWeekends) ([]string, error) {
	return h.weekendRepo.GetWeekends(ctx)
}
