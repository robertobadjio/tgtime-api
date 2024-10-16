package app

import "github.com/robertobadjio/tgtime-api/internal/model/weekend/app/query"

type Application struct {
	Queries Queries
}

type Queries struct {
	GetWeekends         query.GetWeekendsHandler
	GetWeekendsByPeriod query.GetWeekendsByPeriodHandler
}
