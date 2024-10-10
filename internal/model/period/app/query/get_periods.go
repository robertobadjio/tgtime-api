package query

import (
	"context"
	"github.com/robertobadjio/tgtime-api/internal/common/decorator"
	"github.com/robertobadjio/tgtime-api/internal/model/period/domain/period"
	"time"
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
	ps, err := h.periodRepository.GetPeriods(ctx)
	if err != nil {
		return nil, err
	}

	for i, p := range ps {
		start, err := time.Parse(time.RFC3339, p.BeginDate)
		if err != nil {
			return nil, err
		}

		end, err := time.Parse(time.RFC3339, p.EndDate)
		if err != nil {
			panic(err)
		}
		ps[i].WorkingDays = getWorkHoursBetween(start, end)
	}

	return ps, nil
}
