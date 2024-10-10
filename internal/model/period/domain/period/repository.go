package period

import (
	"context"
)

type Repository interface {
	GetPeriod(ctx context.Context, periodId int) (*Period, error)
	GetPeriodCurrent(ctx context.Context) (*Period, error)
	GetPeriods(ctx context.Context) ([]*Period, error)
	CreatePeriod(ctx context.Context, period *Period) (*Period, error)
	UpdatePeriod(ctx context.Context, period *Period) (*Period, error)
	DeletePeriod(ctx context.Context, periodId int) error
}
