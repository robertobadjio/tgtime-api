package app

import (
	"github.com/robertobadjio/tgtime-api/internal/model/user/app/command"
	"github.com/robertobadjio/tgtime-api/internal/model/user/app/command_query"
	"github.com/robertobadjio/tgtime-api/internal/model/user/app/query"
)

type Application struct {
	Commands        Commands
	Queries         Queries
	CommandsQueries CommandsQueries
}

type Commands struct {
	UpdateUser command.UpdateUserHandler
	DeleteUser command.DeleteUserHandler
}

type Queries struct {
	GetUser                    query.GetUserHandler
	GetUserByEmail             query.GetUserByEmailHandler
	GetUserByMacAddress        query.GetUserByMacAddressHandler
	GetUsers                   query.GetUsersHandler
	GetUsersByDepartment       query.GetUsersByDepartmentHandler
	GetUserPasswordHashByEmail query.GetUserPasswordHashByEmailHandler
}

type CommandsQueries struct {
	CreateUser command_query.CreateUserHandler
}
