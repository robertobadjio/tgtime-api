package api

import (
	"context"
	"officetime-api/app/model"
	"officetime-api/app/service"
)

type Service interface {
	Login(ctx context.Context, email, password string) (*service.TokenDetails, error)
	ServiceStatus(ctx context.Context) (int, error)
	GetRouters(ctx context.Context) ([]*model.Router, error)
	/*TokenRefresh(ctx context.Context, ticketID string) (internal.Status, error)
	Logout(ctx context.Context, ticketID, mark string) (int, error)

	GetTimeDayAll(ctx context.Context) (int, error)
	TimeDGetTimeByPeriodayAll(ctx context.Context) (int, error)
	CreateTime(ctx context.Context) (int, error)

	GetAllPeriods(ctx context.Context) (int, error)
	GetPeriod(ctx context.Context) (int, error)
	CreatePeriod(ctx context.Context) (int, error)
	UpdatePeriod(ctx context.Context) (int, error)
	DeletePeriod(ctx context.Context) (int, error)

	GetAllUsers(ctx context.Context) (int, error)
	GetUser(ctx context.Context) (int, error)
	UpdateUser(ctx context.Context) (int, error)
	CreateUser(ctx context.Context) (int, error)
	CreatDeleteUsereUser(ctx context.Context) (int, error)

	GetDepartment(ctx context.Context) (int, error)
	GetAllDepartments(ctx context.Context) (int, error)
	CreateDepartment(ctx context.Context) (int, error)
	UpdateDepartment(ctx context.Context) (int, error)
	DeleteDepartment(ctx context.Context) (int, error)

	GetRouter(ctx context.Context) (int, error)
	GetAllRouters(ctx context.Context) (int, error)
	CreateRouter(ctx context.Context) (int, error)
	UpdateRouter(ctx context.Context) (int, error)
	DeleteRouter(ctx context.Context) (int, error)

	GetStatByPeriodsAndRouters(ctx context.Context) (int, error)
	GetAllTimesDepartmentsByDate(ctx context.Context) (int, error)
	GetStatWorkingPeriod(ctx context.Context) (int, error)

	GetWeekend(ctx context.Context) (int, error)*/
}
