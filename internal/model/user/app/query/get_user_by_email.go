package query

import (
	"context"
	"github.com/robertobadjio/tgtime-api/internal/common/decorator"
	"github.com/robertobadjio/tgtime-api/internal/model/user/domain/user"
)

type GetUserByEmail struct {
	Email string
}

type GetUserByEmailHandler decorator.QueryHandler[GetUserByEmail, *user.User]

type getUserByEmailHandler struct {
	userRepo user.Repository
}

func NewGetUserByEmailHandler(userRepo user.Repository) GetUserByEmailHandler {
	return decorator.ApplyQueryDecorators[GetUserByEmail, *user.User](
		getUserByEmailHandler{userRepo: userRepo},
	)
}

func (h getUserByEmailHandler) Handle(ctx context.Context, qr GetUserByEmail) (*user.User, error) {
	return h.userRepo.GetUserByEmail(ctx, qr.Email)
}
