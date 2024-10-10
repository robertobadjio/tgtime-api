package query

import (
	"context"
	"github.com/robertobadjio/tgtime-api/internal/common/decorator"
	"github.com/robertobadjio/tgtime-api/internal/model/period/domain/period"
)

type GetPeriodCurrent struct{}

type GetPeriodCurrentHandler decorator.QueryHandler[GetPeriodCurrent, *period.Period]

type getPeriodCurrentHandler struct {
	periodRepository period.Repository
}

func NewGetPeriodCurrentHandler(r period.Repository) GetPeriodCurrentHandler {
	return decorator.ApplyQueryDecorators[GetPeriodCurrent, *period.Period](
		getPeriodCurrentHandler{periodRepository: r},
	)
}

func (h getPeriodCurrentHandler) Handle(ctx context.Context, qr GetPeriodCurrent) (*period.Period, error) {
	return h.periodRepository.GetPeriodCurrent(ctx)
}
