package app

import (
	"officetime-api/internal/model/period/app/command"
	"officetime-api/internal/model/period/app/command_query"
	"officetime-api/internal/model/period/app/query"
)

type Application struct {
	Commands        Commands
	Queries         Queries
	CommandsQueries CommandsQueries
}

type Commands struct {
	UpdatePeriod command.UpdatePeriodHandler
	DeletePeriod command.DeletePeriodHandler
}

type Queries struct {
	GetPeriod  query.GetPeriodHandler
	GetPeriods query.GetPeriodsHandler
}

type CommandsQueries struct {
	CreatePeriod command_query.CreatePeriodHandler
}
