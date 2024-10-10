package query

import (
	"context"
	"github.com/robertobadjio/tgtime-api/internal/common/decorator"
	"github.com/robertobadjio/tgtime-api/internal/model/period/domain/period"
	"math"
	"time"
)

type GetPeriod struct {
	PeriodId int
}

type GetPeriodHandler decorator.QueryHandler[GetPeriod, *period.Period]

type getPeriodHandler struct {
	periodRepository period.Repository
}

func NewGetPeriodHandler(r period.Repository) GetPeriodHandler {
	return decorator.ApplyQueryDecorators[GetPeriod, *period.Period](
		getPeriodHandler{periodRepository: r},
	)
}

func (h getPeriodHandler) Handle(ctx context.Context, qr GetPeriod) (*period.Period, error) {
	p, err := h.periodRepository.GetPeriod(ctx, qr.PeriodId)
	if err != nil {
		return nil, err
	}

	start, err := time.Parse(time.RFC3339, p.BeginDate)
	if err != nil {
		return nil, err
	}

	end, err := time.Parse(time.RFC3339, p.EndDate)
	if err != nil {
		panic(err)
	}

	p.WorkingDays = getWorkHoursBetween(start, end)

	return p, err
}

// getWeekdaysBetween
// https://switch-case.ru/61590709
func getWeekdaysBetween(start, end time.Time) int {
	offset := -int(start.Weekday())
	start = start.AddDate(0, 0, -int(start.Weekday()))

	offset += int(end.Weekday())
	if end.Weekday() == time.Sunday {
		offset++
	}
	end = end.AddDate(0, 0, -int(end.Weekday()))

	dif := end.Sub(start).Truncate(time.Hour * 24)
	weeks := (dif.Hours() / 24) / 7

	return int(math.Round(weeks)*5) + offset
}

func getWorkHoursBetween(start, end time.Time) int {
	return getWeekdaysBetween(start, end) * 8
}
