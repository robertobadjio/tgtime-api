package app

import (
	"github.com/robertobadjio/tgtime-api/internal/model/department/app/command"
	"github.com/robertobadjio/tgtime-api/internal/model/department/app/command_query"
	"github.com/robertobadjio/tgtime-api/internal/model/department/app/query"
)

type Application struct {
	Commands        Commands
	Queries         Queries
	CommandsQueries CommandsQueries
}

type Commands struct {
	UpdateDepartment command.UpdateDepartmentHandler
	DeleteDepartment command.DeleteDepartmentHandler
}

type Queries struct {
	GetDepartment  query.GetDepartmentHandler
	GetDepartments query.GetDepartmentsHandler
}

type CommandsQueries struct {
	CreateDepartment command_query.CreateDepartmentHandler
}
