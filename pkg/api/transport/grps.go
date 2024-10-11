package transport

import (
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/robertobadjio/tgtime-api/api/v1/pb/api"
	"github.com/robertobadjio/tgtime-api/pkg/api/endpoints"
)

type grpcServer struct {
	getRouters              grpctransport.Handler
	getUserByMacAddress     grpctransport.Handler
	getGetTimesByTelegramId grpctransport.Handler
	getUserByTelegramId     grpctransport.Handler
	api.UnimplementedApiServer
}

func NewGRPCServer(endpoints endpoints.Set) api.ApiServer {
	return &grpcServer{
		getRouters: grpctransport.NewServer(
			endpoints.GetRoutersEndpoint,
			decodeGRPCGetRoutersRequest,
			encodeGRPCGetRoutersResponse,
		),
		getUserByMacAddress: grpctransport.NewServer(
			endpoints.GetUserByMacAddressEndpoint,
			decodeGRPCGetUserByMacAddressRequest,
			encodeGRPCGetUserByMacAddressResponse,
		),
		getUserByTelegramId: grpctransport.NewServer(
			endpoints.GetUserByTelegramIdEndpoint,
			decodeGRPCGetUserByTelegramIdRequest,
			encodeGRPCGetUserByTelegramIdResponse,
		),
	}
}
