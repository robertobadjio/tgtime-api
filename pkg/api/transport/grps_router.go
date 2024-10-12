package transport

import (
	"context"
	"errors"
	apisvc "github.com/robertobadjio/tgtime-api/api/v1/pb/api"
	"github.com/robertobadjio/tgtime-api/pkg/api/endpoints"
)

func (g *grpcServer) GetRouters(
	ctx context.Context,
	request *apisvc.GetRoutersRequest,
) (*apisvc.GetRoutersResponse, error) {
	_, response, err := g.getRouters.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.(*apisvc.GetRoutersResponse), nil
}

func decodeGRPCGetRoutersRequest(_ context.Context, grpcRequest interface{}) (interface{}, error) {
	_ = grpcRequest.(*apisvc.GetRoutersRequest)
	return endpoints.GetRoutersRequest{}, nil
}

func encodeGRPCGetRoutersResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	resp, ok := grpcResponse.(endpoints.GetRoutersResponse)
	if !ok {
		return nil, errors.New("invalid response body")
	}

	var routers []*apisvc.Router
	for _, r := range resp.Routers {
		router := apisvc.Router{
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

	return &apisvc.GetRoutersResponse{Routers: routers}, nil
}
