package transport

import (
	"context"
	"errors"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"officetime-api/api/v1/pb/api"
	"officetime-api/pkg/api/endpoints"
)

type grpcServer struct {
	getRouters grpctransport.Handler
	api.UnimplementedApiServer
}

func NewGRPCServer(endpoints endpoints.Set) api.ApiServer {
	return &grpcServer{
		getRouters: grpctransport.NewServer(
			endpoints.GetRoutersEndpoint,
			decodeGRPCGetRoutersRequest,
			encodeGRPCGetRoutersResponse,
		),
	}
}

func (g *grpcServer) GetRouters(ctx context.Context, request *api.GetRoutersRequest) (*api.GetRoutersResponse, error) {
	_, response, err := g.getRouters.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.(*api.GetRoutersResponse), nil
}

func decodeGRPCGetRoutersRequest(_ context.Context, grpcRequest interface{}) (interface{}, error) {
	_ = grpcRequest.(*api.GetRoutersRequest)
	return endpoints.GetRoutersRequest{}, nil
}

func encodeGRPCGetRoutersResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	resp, ok := grpcResponse.(endpoints.GetRoutersResponse)
	if !ok {
		return nil, errors.New("invalid response body")
	}

	var routers []*api.Router
	for _, r := range resp.Routers {
		router := api.Router{
			Id:          int64(r.Id),
			Name:        r.Name,
			Description: r.Description,
			Address:     r.Address,
			Login:       r.Login,
			Password:    r.Password,
			Status:      r.Status,
			WorkTime:    r.WorkTime,
		}
		routers = append(routers, &router)
	}

	return &api.GetRoutersResponse{Routers: routers}, nil
}
