package query

import (
	"context"
	"github.com/robertobadjio/tgtime-api/internal/common/decorator"
	"github.com/robertobadjio/tgtime-api/internal/model/user/domain/user"
)

type GetUserPasswordHashByEmail struct {
	Email string
}

type GetUserPasswordHashByEmailHandler decorator.QueryHandler[GetUserPasswordHashByEmail, string]

type getUserPasswordHashByEmailHandler struct {
	userRepo user.Repository
}

func NewGetUserPasswordHashByEmailHandler(userRepo user.Repository) GetUserPasswordHashByEmailHandler {
	return decorator.ApplyQueryDecorators[GetUserPasswordHashByEmail, string](
		getUserPasswordHashByEmailHandler{userRepo: userRepo},
	)
}

func (h getUserPasswordHashByEmailHandler) Handle(ctx context.Context, qr GetUserPasswordHashByEmail) (string, error) {
	return h.userRepo.GetUserPasswordHashByEmail(ctx, qr.Email)
}
