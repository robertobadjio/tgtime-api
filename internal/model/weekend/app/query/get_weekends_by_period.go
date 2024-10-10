package query

import (
	"context"
	"github.com/robertobadjio/tgtime-api/internal/common/decorator"
	"github.com/robertobadjio/tgtime-api/internal/model/weekend/domain/weekend"
	"time"
)

type GetWeekendsByPeriod struct {
	start, end time.Time
}

type GetWeekendsByPeriodHandler decorator.QueryHandler[GetWeekendsByPeriod, map[string]bool]

/*type GetWeekendsReadModel interface {
	GetWeekends(ctx context.Context) ([]*weekend.Weekend, error)
}*/

type getWeekendsByPeriodHandler struct {
	weekendRepo weekend.Repository
}

func NewGetWeekendsByPeriodHandler(weekendRepo weekend.Repository) GetWeekendsByPeriodHandler {
	return decorator.ApplyQueryDecorators[GetWeekendsByPeriod, map[string]bool](
		getWeekendsByPeriodHandler{weekendRepo: weekendRepo},
	)
}

func (h getWeekendsByPeriodHandler) Handle(ctx context.Context, gw GetWeekendsByPeriod) (map[string]bool, error) {
	return h.weekendRepo.GetWeekendsByPeriod(ctx, gw.start, gw.end)
}
