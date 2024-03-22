package command

import (
	"context"
	"officetime-api/internal/common/decorator"
	"officetime-api/internal/model/period/domain/period"
)

type DeletePeriod struct {
	PeriodId int
}

type DeletePeriodHandler decorator.CommandHandler[DeletePeriod]

type deletePeriodHandler struct {
	periodRepository period.Repository
}

func NewDeletePeriodHandler(periodRepository period.Repository) DeletePeriodHandler {
	if periodRepository == nil {
		panic("nil periodRepository")
	}

	return decorator.ApplyCommandDecorators[DeletePeriod](
		deletePeriodHandler{periodRepository: periodRepository},
	)
}

func (h deletePeriodHandler) Handle(ctx context.Context, cmd DeletePeriod) error {
	err := h.periodRepository.DeletePeriod(ctx, cmd.PeriodId)
	if err != nil {
		return err
	}

	return nil
}
