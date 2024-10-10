package app

import (
	"github.com/robertobadjio/tgtime-api/internal/model/router/app/command"
	"github.com/robertobadjio/tgtime-api/internal/model/router/app/command_query"
	"github.com/robertobadjio/tgtime-api/internal/model/router/app/query"
)

type Application struct {
	Commands        Commands
	Queries         Queries
	CommandsQueries CommandsQueries
}

type Commands struct {
	UpdateRouter command.UpdateRouterHandler
	DeleteRouter command.DeleteRouterHandler
}

type Queries struct {
	GetRouter  query.GetRouterHandler
	GetRouters query.GetRoutersHandler
}

type CommandsQueries struct {
	CreateRouter command_query.CreateRouterHandler
}
