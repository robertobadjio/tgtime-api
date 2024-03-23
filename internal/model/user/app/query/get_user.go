package query

import (
	"context"
	"officetime-api/internal/common/decorator"
	"officetime-api/internal/model/user/domain/user"
)

type GetUser struct {
	UserId int
}

type GetUserHandler decorator.QueryHandler[GetUser, *user.User]

/*type GetUserReadModel interface {
	GetUser(ctx context.Context, userId int) (*user.User, error)
}*/

type getUserHandler struct {
	userRepo user.Repository
}

func NewGetUserHandler(userRepo user.Repository) GetUserHandler {
	return decorator.ApplyQueryDecorators[GetUser, *user.User](
		getUserHandler{userRepo: userRepo},
	)
}

func (h getUserHandler) Handle(ctx context.Context, qr GetUser) (*user.User, error) {
	return h.userRepo.GetUser(ctx, qr.UserId)
}
