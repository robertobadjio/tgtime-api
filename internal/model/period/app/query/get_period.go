package query

import (
	"context"
	"officetime-api/internal/common/decorator"
	"officetime-api/internal/model/period/domain/period"
)

type GetPeriod struct {
	PeriodId int
}

type GetPeriodHandler decorator.QueryHandler[GetPeriod, *period.Period]

type getPeriodHandler struct {
	routerRepository period.Repository
}

func NewGetPeriodHandler(r period.Repository) GetPeriodHandler {
	return decorator.ApplyQueryDecorators[GetPeriod, *period.Period](
		getPeriodHandler{routerRepository: r},
	)
}

func (h getPeriodHandler) Handle(ctx context.Context, qr GetPeriod) (*period.Period, error) {
	return h.routerRepository.GetPeriod(ctx, qr.PeriodId)
}
