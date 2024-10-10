package transport

import (
	"context"
	"errors"
	"github.com/robertobadjio/tgtime-api/api/v1/pb/api"
	"github.com/robertobadjio/tgtime-api/pkg/api/endpoints"
)

func (g *grpcServer) GetUserByMacAddress(
	ctx context.Context,
	request *api.GetUserByMacAddressRequest,
) (*api.GetUserByMacAddressResponse, error) {
	_, response, err := g.getUserByMacAddress.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.(*api.GetUserByMacAddressResponse), nil
}

func decodeGRPCGetUserByMacAddressRequest(_ context.Context, grpcRequest interface{}) (interface{}, error) {
	_ = grpcRequest.(*api.GetUserByMacAddressRequest)
	return endpoints.GetUserByMacAddressRequest{}, nil
}

func encodeGRPCGetUserByMacAddressResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	resp, ok := grpcResponse.(endpoints.GetUserByMacAddressResponse)
	if !ok {
		return nil, errors.New("invalid response body")
	}

	user := api.User{
		Id:         int64(resp.User.Id),
		Name:       resp.User.Name,
		Surname:    resp.User.Surname,
		Lastname:   resp.User.Lastname,
		BirthDate:  resp.User.BirthDate,
		Email:      resp.User.Email,
		MacAddress: resp.User.MacAddress,
		TelegramId: resp.User.TelegramId,
		Role:       resp.User.Role,
		Department: resp.User.Department,
		Position:   resp.User.Position,
	}

	return &api.GetUserByMacAddressResponse{User: &user}, nil
}
