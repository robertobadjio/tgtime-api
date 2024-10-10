package app

import (
	"github.com/robertobadjio/tgtime-api/internal/model/period/app/command"
	"github.com/robertobadjio/tgtime-api/internal/model/period/app/command_query"
	"github.com/robertobadjio/tgtime-api/internal/model/period/app/query"
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
	GetPeriod        query.GetPeriodHandler
	GetPeriodCurrent query.GetPeriodCurrentHandler
	GetPeriods       query.GetPeriodsHandler
}

type CommandsQueries struct {
	CreatePeriod command_query.CreatePeriodHandler
}
