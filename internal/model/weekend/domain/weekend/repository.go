package weekend

import (
	"context"
	"time"
)

type Repository interface {
	GetWeekends(ctx context.Context) ([]string, error)
	GetWeekendsByPeriod(ctx context.Context, start, end time.Time) (map[string]bool, error)
}
