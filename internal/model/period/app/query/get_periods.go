package query

import (
	"context"
	"officetime-api/internal/common/decorator"
	"officetime-api/internal/model/period/domain/period"
)

type GetPeriods struct {
}

type GetPeriodsHandler decorator.QueryHandler[GetPeriods, []*period.Period]

type getPeriodsHandler struct {
	periodRepository period.Repository
}

func NewGetPeriodsHandler(r period.Repository) GetPeriodsHandler {
	return decorator.ApplyQueryDecorators[GetPeriods, []*period.Period](
		getPeriodsHandler{periodRepository: r},
	)
}

func (h getPeriodsHandler) Handle(ctx context.Context, _ GetPeriods) ([]*period.Period, error) {
	return h.periodRepository.GetPeriods(ctx)
}
