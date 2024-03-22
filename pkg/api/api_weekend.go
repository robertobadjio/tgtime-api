package api

import (
	"context"
	"officetime-api/internal/model/weekend/app/query"
)

func (s *apiService) GetWeekends(ctx context.Context) ([]string, error) {
	qr := query.GetWeekends{}
	routers, err := s.weekendApp.Queries.GetWeekends.Handle(ctx, qr)
	if err != nil {
		return nil, err
	}

	return routers, nil
}
