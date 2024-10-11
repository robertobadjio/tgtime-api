package api

import (
	"context"
	"github.com/robertobadjio/tgtime-api/internal/model/department/domain/department"
	"github.com/robertobadjio/tgtime-api/internal/model/period/domain/period"
	"github.com/robertobadjio/tgtime-api/internal/model/router/domain/router"
	"github.com/robertobadjio/tgtime-api/internal/model/user/domain/user"
)

type Service interface {
	//Login(ctx context.Context, email, password string) (*service.TokenDetails, error)
	//TokenRefresh(ctx context.Context, ticketID string) (internal.Status, error)
	//Logout(ctx context.Context, ticketID, mark string) (int, error)

	ServiceStatus(ctx context.Context) (int, error)

	GetRouters(ctx context.Context) ([]*router.Router, error)
	GetRouter(ctx context.Context, routerId int) (*router.Router, error)
	CreateRouter(ctx context.Context, router *router.Router) (*router.Router, error)
	UpdateRouter(ctx context.Context, routerId int, router *router.Router) (*router.Router, error)
	DeleteRouter(ctx context.Context, routerId int) error

	GetPeriod(ctx context.Context, periodId int) (*period.Period, error)
	GetPeriodCurrent(ctx context.Context) (*period.Period, error)
	GetPeriods(ctx context.Context) ([]*period.Period, error)
	CreatePeriod(ctx context.Context, period *period.Period) (*period.Period, error)
	UpdatePeriod(ctx context.Context, periodId int, period *period.Period) (*period.Period, error)
	DeletePeriod(ctx context.Context, periodId int) error

	GetDepartment(ctx context.Context, departmentId int) (*department.Department, error)
	GetDepartments(ctx context.Context) ([]*department.Department, error)
	CreateDepartment(ctx context.Context, department *department.Department) (*department.Department, error)
	UpdateDepartment(ctx context.Context, departmentId int, department *department.Department) (*department.Department, error)
	DeleteDepartment(ctx context.Context, departmentId int) error

	GetWeekends(ctx context.Context) ([]string, error)

	GetUsers(ctx context.Context, offset, limit int) ([]*user.User, error)
	GetUser(ctx context.Context, userId int) (*user.User, error)
	GetUserByMacAddress(ctx context.Context, macAddress string) (*user.User, error)
	GetUserByTelegramId(ctx context.Context, telegramId int64) (*user.User, error)
	UpdateUser(ctx context.Context, user *user.User) (*user.User, error)
	CreateUser(ctx context.Context, user *user.User) (*user.User, error)
	DeleteUser(ctx context.Context, userId int) error

	/*GetTimeDayAll(ctx context.Context) (int, error)
	TimeDGetTimeByPeriodayAll(ctx context.Context) (int, error)
	CreateTime(ctx context.Context) (int, error)

	GetStatByPeriodsAndRouters(ctx context.Context) (int, error)
	GetAllTimesDepartmentsByDate(ctx context.Context) (int, error)
	GetStatWorkingPeriod(ctx context.Context) (int, error)
	*/
}
