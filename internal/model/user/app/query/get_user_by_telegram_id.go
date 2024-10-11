package query

import (
	"context"
	"github.com/robertobadjio/tgtime-api/internal/common/decorator"
	"github.com/robertobadjio/tgtime-api/internal/model/user/domain/user"
)

type GetUserByTelegramId struct {
	TelegramId int64
}

type GetUserByTelegramIdHandler decorator.QueryHandler[GetUserByTelegramId, *user.User]

type getUserByTelegramIdHandler struct {
	userRepo user.Repository
}

func NewGetUserByTelegramIdHandler(userRepo user.Repository) GetUserByTelegramIdHandler {
	return decorator.ApplyQueryDecorators[GetUserByTelegramId, *user.User](
		getUserByTelegramIdHandler{userRepo: userRepo},
	)
}

func (h getUserByTelegramIdHandler) Handle(ctx context.Context, qr GetUserByTelegramId) (*user.User, error) {
	return h.userRepo.GetUserByTelegramId(ctx, qr.TelegramId)
}
