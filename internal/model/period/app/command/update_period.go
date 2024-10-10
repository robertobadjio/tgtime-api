package command

import (
	"context"
	"github.com/robertobadjio/tgtime-api/internal/common/decorator"
	"github.com/robertobadjio/tgtime-api/internal/model/period/domain/period"
)

type UpdatePeriod struct {
	Period *period.Period
}

type UpdatePeriodHandler decorator.CommandHandler[UpdatePeriod]

type updatePeriodHandler struct {
	periodRepository period.Repository
}

func NewUpdatePeriodHandler(periodRepository period.Repository) UpdatePeriodHandler {
	if periodRepository == nil {
		panic("nil periodRepository")
	}

	return decorator.ApplyCommandDecorators[UpdatePeriod](
		updatePeriodHandler{periodRepository: periodRepository},
	)
}

func (h updatePeriodHandler) Handle(ctx context.Context, cmd UpdatePeriod) error {
	_, err := h.periodRepository.UpdatePeriod(ctx, cmd.Period) // TODO: !
	if err != nil {
		return err
	}

	return nil
}
