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
	req := grpcRequest.(*api.GetUserByMacAddressRequest)
	return endpoints.GetUserByMacAddressRequest{MacAddress: req.MacAddress}, nil
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

func (g *grpcServer) GetUserByTelegramId(
	ctx context.Context,
	request *api.GetUserByTelegramIdRequest,
) (*api.GetUserByTelegramIdResponse, error) {
	_, response, err := g.getUserByTelegramId.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.(*api.GetUserByTelegramIdResponse), nil
}

func decodeGRPCGetUserByTelegramIdRequest(_ context.Context, grpcRequest interface{}) (interface{}, error) {
	req := grpcRequest.(*api.GetUserByTelegramIdRequest)
	return endpoints.GetUserByTelegramIdRequest{TelegramId: req.TelegramId}, nil
}

func encodeGRPCGetUserByTelegramIdResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	resp, ok := grpcResponse.(endpoints.GetUserByTelegramIdResponse)
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

	return &api.GetUserByTelegramIdResponse{User: &user}, nil
}
