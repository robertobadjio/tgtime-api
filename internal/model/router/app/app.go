package app

import (
	"officetime-api/internal/model/router/app/command"
	"officetime-api/internal/model/router/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateRouter command.CreateRouterHandler
	UpdateRouter command.UpdateRouterHandler
	DeleteRouter command.DeleteRouterHandler
}

type Queries struct {
	GetRouter  query.GetRouterHandler
	GetRouters query.GetRoutersHandler
}
