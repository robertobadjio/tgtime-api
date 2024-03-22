package command_query

import (
	"context"
	"officetime-api/internal/common/decorator"
	"officetime-api/internal/model/period/domain/period"
)

type CreatePeriod struct {
	Period *period.Period
}

type CreatePeriodHandler decorator.CommandQueryHandler[CreatePeriod, *period.Period]

type createPeriodHandler struct {
	periodRepository period.Repository
}

func NewCreatePeriodHandler(periodRepository period.Repository) CreatePeriodHandler {
	if periodRepository == nil {
		panic("nil periodRepository")
	}

	return decorator.ApplyCommandQueryDecorators[CreatePeriod, *period.Period](
		createPeriodHandler{periodRepository: periodRepository},
	)
}

func (h createPeriodHandler) Handle(ctx context.Context, cmdQr CreatePeriod) (*period.Period, error) {
	periodNew, err := h.periodRepository.CreatePeriod(ctx, cmdQr.Period)
	if err != nil {
		return nil, err
	}

	return periodNew, nil
}
