package endpoints

import (
	"github.com/robertobadjio/tgtime-api/internal/model/department/domain/department"
	"github.com/robertobadjio/tgtime-api/internal/model/period/domain/period"
	"github.com/robertobadjio/tgtime-api/internal/model/router/domain/router"
	"github.com/robertobadjio/tgtime-api/internal/model/user/domain/user"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginResponse struct {
	AccessToken         string `json:"access_token"`
	RefreshToken        string `json:"refresh_token"`
	AccessTokenExpires  int64  `json:"access_token_expires"`
	RefreshTokenExpires int64  `json:"refresh_token_expires"`
}

type ServiceStatusRequest struct{}
type ServiceStatusResponse struct {
	Code int    `json:"status"`
	Err  string `json:"err,omitempty"`
}

type GetRoutersRequest struct{}
type GetRoutersResponse struct {
	Routers []*router.Router `json:"routers"`
}

type GetRouterRequest struct {
	RouterId int `json:"router_id"`
}
type GetRouterResponse struct {
	Router *router.Router `json:"router"`
}

type CreateRouterRequest struct {
	Router *router.Router `json:"router"`
}
type CreateRouterResponse struct {
	Router *router.Router `json:"router"`
}

type UpdateRouterRequest struct {
	RouterId int            `json:"router_id"`
	Router   *router.Router `json:"router"`
}
type UpdateRouterResponse struct {
	Router *router.Router `json:"router"`
}

type DeleteRouterRequest struct {
	RouterId int `json:"router_id"`
}
type DeleteRouterResponse struct{}

type GetPeriodsRequest struct{}
type GetPeriodsResponse struct {
	Periods []*period.Period `json:"periods"`
}

type GetPeriodCurrentRequest struct{}

type GetPeriodRequest struct {
	PeriodId int `json:"period_id"`
}
type GetPeriodResponse struct {
	Period *period.Period `json:"period"`
}

type CreatePeriodRequest struct {
	Period *period.Period `json:"period"`
}
type CreatePeriodResponse struct {
	Period *period.Period `json:"period"`
}

type UpdatePeriodRequest struct {
	PeriodId int            `json:"period_id"`
	Period   *period.Period `json:"period"`
}
type UpdatePeriodResponse struct {
	Period *period.Period `json:"period"`
}

type DeletePeriodRequest struct {
	PeriodId int `json:"period_id"`
}
type DeletePeriodResponse struct{}

type GetDepartmentsRequest struct{}
type GetDepartmentsResponse struct {
	Departments []*department.Department `json:"departments"`
}

type GetDepartmentRequest struct {
	DepartmentId int `json:"department_id"`
}
type GetDepartmentResponse struct {
	Department *department.Department `json:"department"`
}

type CreateDepartmentRequest struct {
	Department *department.Department `json:"department"`
}
type CreateDepartmentResponse struct {
	Department *department.Department `json:"department"`
}

type UpdateDepartmentRequest struct {
	DepartmentId int                    `json:"department_id"`
	Department   *department.Department `json:"department"`
}
type UpdateDepartmentResponse struct {
	Department *department.Department `json:"department"`
}

type DeleteDepartmentRequest struct {
	DepartmentId int `json:"department_id"`
}
type DeleteDepartmentResponse struct{}

type GetWeekendsRequest struct{}
type GetWeekendsResponse struct {
	Weekends []string `json:"weekends"`
}

type GetUsersRequest struct {
	Offset, Limit int
}
type GetUsersResponse struct {
	Users []*user.User `json:"users"`
}

type GetUserRequest struct {
	UserId int `json:"user_id"`
}
type GetUserResponse struct {
	User *user.User `json:"user"`
}

type GetUserByMacAddressRequest struct {
	MacAddress string `json:"mac_address"`
}
type GetUserByMacAddressResponse struct {
	User *user.User `json:"user"`
}

type CreateUserRequest struct {
	User *user.User `json:"user"`
}
type CreateUserResponse struct {
	User *user.User `json:"user"`
}

type UpdateUserRequest struct {
	UserId int        `json:"user_id"`
	User   *user.User `json:"user"`
}
type UpdateUserResponse struct {
	User *user.User `json:"user"`
}

type DeleteUserRequest struct {
	UserId int `json:"user_id"`
}
type DeleteUserResponse struct{}
