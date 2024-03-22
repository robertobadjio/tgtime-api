package app

import (
	"officetime-api/internal/model/router/app/command"
	"officetime-api/internal/model/router/app/command_query"
	"officetime-api/internal/model/router/app/query"
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
