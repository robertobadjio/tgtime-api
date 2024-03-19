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
	GetRouterEndpoint     endpoint.Endpoint
	CreateRouterEndpoint  endpoint.Endpoint
	UpdateRouterEndpoint  endpoint.Endpoint
	DeleteRouterEndpoint  endpoint.Endpoint
}

func NewEndpointSet(svc api.Service) Set {
	return Set{
		LoginEndpoint:         MakeLoginEndpoint(svc),
		ServiceStatusEndpoint: MakeServiceStatusEndpoint(svc),
		GetRoutersEndpoint:    MakeGetRoutersEndpoint(svc),
		GetRouterEndpoint:     MakeGetRouterEndpoint(svc),
		CreateRouterEndpoint:  MakeCreateRouterEndpoint(svc),
		UpdateRouterEndpoint:  MakeUpdateRouterEndpoint(svc),
		DeleteRouterEndpoint:  MakeDeleteRouterEndpoint(svc),
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
