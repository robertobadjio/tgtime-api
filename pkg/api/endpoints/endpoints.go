package endpoints

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"officetime-api/pkg/api"
)

type Set struct {
	LoginEndpoint         endpoint.Endpoint
	ServiceStatusEndpoint endpoint.Endpoint
	GetRoutersEndpoint    endpoint.Endpoint
}

func NewEndpointSet(svc api.Service) Set {
	return Set{
		LoginEndpoint:         MakeLoginEndpoint(svc),
		ServiceStatusEndpoint: MakeServiceStatusEndpoint(svc),
		GetRoutersEndpoint:    MakeGetRoutersEndpoint(svc),
	}
}

func MakeLoginEndpoint(svc api.Service) endpoint.Endpoint {
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
}

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
