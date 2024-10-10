package query

import (
	"context"
	"github.com/robertobadjio/tgtime-api/internal/common/decorator"
	"github.com/robertobadjio/tgtime-api/internal/model/user/domain/user"
)

type GetUserByMacAddress struct {
	MacAddress string
}

type GetUserByMacAddressHandler decorator.QueryHandler[GetUserByMacAddress, *user.User]

type getUserByMacAddressHandler struct {
	userRepo user.Repository
}

func NewGetUserByMacAddressHandler(userRepo user.Repository) GetUserByMacAddressHandler {
	return decorator.ApplyQueryDecorators[GetUserByMacAddress, *user.User](
		getUserByMacAddressHandler{userRepo: userRepo},
	)
}

func (h getUserByMacAddressHandler) Handle(ctx context.Context, qr GetUserByMacAddress) (*user.User, error) {
	return h.userRepo.GetUserByMacAddress(ctx, qr.MacAddress)
}
