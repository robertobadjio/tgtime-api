package endpoints

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/robertobadjio/tgtime-api/pkg/api"
)

type Set struct {
	LoginEndpoint endpoint.Endpoint

	ServiceStatusEndpoint endpoint.Endpoint

	GetRoutersEndpoint   endpoint.Endpoint
	GetRouterEndpoint    endpoint.Endpoint
	CreateRouterEndpoint endpoint.Endpoint
	UpdateRouterEndpoint endpoint.Endpoint
	DeleteRouterEndpoint endpoint.Endpoint

	GetPeriodEndpoint        endpoint.Endpoint
	GetPeriodCurrentEndpoint endpoint.Endpoint
	GetPeriodsEndpoint       endpoint.Endpoint
	CreatePeriodEndpoint     endpoint.Endpoint
	UpdatePeriodEndpoint     endpoint.Endpoint
	DeletePeriodEndpoint     endpoint.Endpoint

	GetDepartmentEndpoint    endpoint.Endpoint
	GetDepartmentsEndpoint   endpoint.Endpoint
	CreateDepartmentEndpoint endpoint.Endpoint
	UpdateDepartmentEndpoint endpoint.Endpoint
	DeleteDepartmentEndpoint endpoint.Endpoint

	GetWeekendsEndpoint endpoint.Endpoint

	GetUserEndpoint             endpoint.Endpoint
	GetUserByMacAddressEndpoint endpoint.Endpoint
	GetUsersEndpoint            endpoint.Endpoint
	CreateUserEndpoint          endpoint.Endpoint
	UpdateUserEndpoint          endpoint.Endpoint
	DeleteUserEndpoint          endpoint.Endpoint
}

func NewEndpointSet(svc api.Service) Set {
	return Set{
		//LoginEndpoint: MakeLoginEndpoint(svc),

		ServiceStatusEndpoint: MakeServiceStatusEndpoint(svc),

		GetRoutersEndpoint:   MakeGetRoutersEndpoint(svc),
		GetRouterEndpoint:    MakeGetRouterEndpoint(svc),
		CreateRouterEndpoint: MakeCreateRouterEndpoint(svc),
		UpdateRouterEndpoint: MakeUpdateRouterEndpoint(svc),
		DeleteRouterEndpoint: MakeDeleteRouterEndpoint(svc),

		GetPeriodEndpoint:        MakeGetPeriodEndpoint(svc),
		GetPeriodCurrentEndpoint: MakeGetPeriodCurrentEndpoint(svc),
		GetPeriodsEndpoint:       MakeGetPeriodsEndpoint(svc),
		CreatePeriodEndpoint:     MakeCreatePeriodEndpoint(svc),
		UpdatePeriodEndpoint:     MakeUpdatePeriodEndpoint(svc),
		DeletePeriodEndpoint:     MakeDeletePeriodEndpoint(svc),

		GetDepartmentEndpoint:    MakeGetDepartmentEndpoint(svc),
		GetDepartmentsEndpoint:   MakeGetDepartmentsEndpoint(svc),
		CreateDepartmentEndpoint: MakeCreateDepartmentEndpoint(svc),
		UpdateDepartmentEndpoint: MakeUpdateDepartmentEndpoint(svc),
		DeleteDepartmentEndpoint: MakeDeleteDepartmentEndpoint(svc),

		GetWeekendsEndpoint: MakeGetWeekendsEndpoint(svc),

		GetUserEndpoint:             MakeGetUserEndpoint(svc),
		GetUserByMacAddressEndpoint: MakeGetUserByMacAddressEndpoint(svc),
		GetUsersEndpoint:            MakeGetUsersEndpoint(svc),
		CreateUserEndpoint:          MakeCreateUserEndpoint(svc),
		UpdateUserEndpoint:          MakeUpdateUserEndpoint(svc),
		DeleteUserEndpoint:          MakeDeleteUserEndpoint(svc),
	}
}

/*func MakeLoginEndpoint(svc api.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(LoginRequest)
		authData, err := svc.Login(ctx, req.Email, req.Password)
		if err != nil {
			return LoginResponse{"", "", 0, 0}, err
		}
		return LoginResponse{
				authData.AccessToken,
				authData.RefreshToken,
				authData.AccessTokenExpires,
				authData.RefreshTokenExpires,
			},
			nil
	}
}*/

func MakeServiceStatusEndpoint(svc api.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(ServiceStatusRequest)
		code, err := svc.ServiceStatus(ctx)
		if err != nil {
			return ServiceStatusResponse{Code: code, Err: err.Error()}, nil
		}
		return ServiceStatusResponse{Code: code, Err: ""}, nil
	}
}

func MakeGetRoutersEndpoint(svc api.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(GetRoutersRequest)
		routers, err := svc.GetRouters(ctx)
		if err != nil {
			return GetRoutersResponse{Routers: routers}, nil
		}
		return GetRoutersResponse{Routers: routers}, nil
	}
}

func MakeGetRouterEndpoint(svc api.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetRouterRequest)
		router, err := svc.GetRouter(ctx, req.RouterId)
		if err != nil {
			return GetRouterResponse{}, err
		}
		return GetRouterResponse{Router: router}, nil
	}
}

func MakeCreateRouterEndpoint(svc api.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateRouterRequest)
		router, err := svc.CreateRouter(ctx, req.Router)
		if err != nil {
			return CreateRouterResponse{}, err
		}
		return CreateRouterResponse{Router: router}, nil
	}
}

func MakeUpdateRouterEndpoint(svc api.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateRouterRequest)
		router, err := svc.UpdateRouter(ctx, req.RouterId, req.Router)
		if err != nil {
			return UpdateRouterResponse{}, err
		}
		return UpdateRouterResponse{Router: router}, nil
	}
}

func MakeDeleteRouterEndpoint(svc api.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteRouterRequest)
		err := svc.DeleteRouter(ctx, req.RouterId)
		if err != nil {
			return DeleteRouterResponse{}, err
		}
		return DeleteRouterResponse{}, nil
	}
}

func MakeGetPeriodEndpoint(svc api.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetPeriodRequest)
		period, err := svc.GetPeriod(ctx, req.PeriodId)
		if err != nil {
			return GetPeriodResponse{}, err
		}
		return GetPeriodResponse{Period: period}, nil
	}
}

func MakeGetPeriodCurrentEndpoint(svc api.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(GetPeriodCurrentRequest)
		period, err := svc.GetPeriodCurrent(ctx)
		if err != nil {
			return GetPeriodResponse{}, err
		}
		return GetPeriodResponse{Period: period}, nil
	}
}

func MakeGetPeriodsEndpoint(svc api.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(GetPeriodsRequest)
		periods, err := svc.GetPeriods(ctx)
		if err != nil {
			return GetPeriodsResponse{Periods: periods}, nil
		}
		return GetPeriodsResponse{Periods: periods}, nil
	}
}

func MakeCreatePeriodEndpoint(svc api.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreatePeriodRequest)
		period, err := svc.CreatePeriod(ctx, req.Period)
		if err != nil {
			return CreatePeriodResponse{}, err
		}
		return CreatePeriodResponse{Period: period}, nil
	}
}

func MakeUpdatePeriodEndpoint(svc api.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdatePeriodRequest)
		period, err := svc.UpdatePeriod(ctx, req.PeriodId, req.Period)
		if err != nil {
			return UpdatePeriodResponse{}, err
		}
		return UpdatePeriodResponse{Period: period}, nil
	}
}

func MakeDeletePeriodEndpoint(svc api.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeletePeriodRequest)
		err := svc.DeletePeriod(ctx, req.PeriodId)
		if err != nil {
			return DeletePeriodResponse{}, err
		}
		return DeletePeriodResponse{}, nil
	}
}

func MakeGetDepartmentEndpoint(svc api.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetDepartmentRequest)
		department, err := svc.GetDepartment(ctx, req.DepartmentId)
		if err != nil {
			return GetDepartmentResponse{}, err
		}
		return GetDepartmentResponse{Department: department}, nil
	}
}

func MakeGetDepartmentsEndpoint(svc api.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(GetDepartmentsRequest)
		departments, err := svc.GetDepartments(ctx)
		if err != nil {
			return GetDepartmentsResponse{Departments: departments}, nil
		}
		return GetDepartmentsResponse{Departments: departments}, nil
	}
}

func MakeCreateDepartmentEndpoint(svc api.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateDepartmentRequest)
		department, err := svc.CreateDepartment(ctx, req.Department)
		if err != nil {
			return CreateDepartmentResponse{}, err
		}
		return CreateDepartmentResponse{Department: department}, nil
	}
}

func MakeUpdateDepartmentEndpoint(svc api.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateDepartmentRequest)
		department, err := svc.UpdateDepartment(ctx, req.DepartmentId, req.Department)
		if err != nil {
			return UpdateDepartmentResponse{}, err
		}
		return UpdateDepartmentResponse{Department: department}, nil
	}
}

func MakeDeleteDepartmentEndpoint(svc api.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteDepartmentRequest)
		err := svc.DeleteDepartment(ctx, req.DepartmentId)
		if err != nil {
			return DeleteDepartmentResponse{}, err
		}
		return DeleteDepartmentResponse{}, nil
	}
}

func MakeGetWeekendsEndpoint(svc api.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(GetWeekendsRequest)
		weekends, err := svc.GetWeekends(ctx)
		if err != nil {
			return GetWeekendsResponse{Weekends: weekends}, nil
		}
		return GetWeekendsResponse{Weekends: weekends}, nil
	}
}

func MakeGetUserEndpoint(svc api.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetUserRequest)
		user, err := svc.GetUser(ctx, req.UserId)
		if err != nil {
			return GetUserResponse{}, err
		}
		return GetUserResponse{User: user}, nil
	}
}

func MakeGetUserByMacAddressEndpoint(svc api.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetUserByMacAddressRequest)
		user, err := svc.GetUserByMacAddress(ctx, req.MacAddress)
		if err != nil {
			return GetUserByMacAddressResponse{}, err
		}
		return GetUserByMacAddressResponse{User: user}, nil
	}
}

func MakeGetUsersEndpoint(svc api.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetUsersRequest)
		users, err := svc.GetUsers(ctx, req.Offset, req.Limit)
		if err != nil {
			return GetUsersResponse{Users: users}, nil
		}
		return GetUsersResponse{Users: users}, nil
	}
}

func MakeCreateUserEndpoint(svc api.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateUserRequest)
		user, err := svc.CreateUser(ctx, req.User)
		if err != nil {
			return CreateUserResponse{}, err
		}
		return CreateUserResponse{User: user}, nil
	}
}

func MakeUpdateUserEndpoint(svc api.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateUserRequest)
		if req.UserId != req.User.Id {
			// TODO: Error
		}
		user, err := svc.UpdateUser(ctx, req.User)
		if err != nil {
			return UpdateUserResponse{}, err
		}
		return UpdateUserResponse{User: user}, nil
	}
}

func MakeDeleteUserEndpoint(svc api.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteUserRequest)
		err := svc.DeleteUser(ctx, req.UserId)
		if err != nil {
			return DeleteUserResponse{}, err
		}
		return DeleteUserResponse{}, nil
	}
}
