package api

import (
	"context"
	"errors"
	"github.com/robertobadjio/tgtime-api/internal/model/period/app/command"
	"github.com/robertobadjio/tgtime-api/internal/model/period/app/command_query"
	"github.com/robertobadjio/tgtime-api/internal/model/period/app/query"
	"github.com/robertobadjio/tgtime-api/internal/model/period/domain/period"
)

func (s *apiService) GetPeriod(ctx context.Context, periodId int) (*period.Period, error) {
	qr := query.GetPeriod{PeriodId: periodId}
	p, err := s.periodApp.Queries.GetPeriod.Handle(ctx, qr)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (s *apiService) GetPeriodCurrent(ctx context.Context) (*period.Period, error) {
	qr := query.GetPeriodCurrent{}
	p, err := s.periodApp.Queries.GetPeriodCurrent.Handle(ctx, qr)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (s *apiService) GetPeriods(ctx context.Context) ([]*period.Period, error) {
	qr := query.GetPeriods{}
	periods, err := s.periodApp.Queries.GetPeriods.Handle(ctx, qr)
	if err != nil {
		return nil, err
	}

	return periods, nil
}

func (s *apiService) CreatePeriod(ctx context.Context, period *period.Period) (*period.Period, error) {
	cmd := command_query.CreatePeriod{Period: period}
	periodNew, err := s.periodApp.CommandsQueries.CreatePeriod.Handle(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return periodNew, nil
}

func (s *apiService) UpdatePeriod(ctx context.Context, periodId int, period *period.Period) (*period.Period, error) {
	if periodId != period.Id {
		return nil, errors.New("error update ids not equals")
	}

	cmd := command.UpdatePeriod{Period: period}
	err := s.periodApp.Commands.UpdatePeriod.Handle(ctx, cmd)
	if err != nil {
		return nil, err
	}

	qr := query.GetPeriod{PeriodId: periodId}
	periodNew, err := s.periodApp.Queries.GetPeriod.Handle(ctx, qr)
	if err != nil {
		return nil, err
	}

	return periodNew, nil
}

func (s *apiService) DeletePeriod(ctx context.Context, periodId int) error {
	cmd := command.DeletePeriod{PeriodId: periodId}
	err := s.periodApp.Commands.DeletePeriod.Handle(ctx, cmd)
	if err != nil {
		return err
	}

	return nil
}
