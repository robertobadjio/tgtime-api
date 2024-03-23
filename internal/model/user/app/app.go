package app

import (
	"officetime-api/internal/model/user/app/command"
	"officetime-api/internal/model/user/app/command_query"
	"officetime-api/internal/model/user/app/query"
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
	GetUsers                   query.GetUsersHandler
	GetUsersByDepartment       query.GetUsersByDepartmentHandler
	GetUserPasswordHashByEmail query.GetUserPasswordHashByEmailHandler
}

type CommandsQueries struct {
	CreateUser command_query.CreateUserHandler
}
