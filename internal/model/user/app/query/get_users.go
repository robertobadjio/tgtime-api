package query

import (
	"context"
	"officetime-api/internal/common/decorator"
	"officetime-api/internal/model/user/domain/user"
)

type GetUsers struct {
	Offset, Limit int
}

type GetUsersHandler decorator.QueryHandler[GetUsers, []*user.User]

/*type GetUsersReadModel interface {
	GetUsers(ctx context.Context) ([]*user.User, error)
}*/

type getUsersHandler struct {
	userRepo user.Repository
}

func NewGetUsersHandler(userRepo user.Repository) GetUsersHandler {
	return decorator.ApplyQueryDecorators[GetUsers, []*user.User](
		getUsersHandler{userRepo: userRepo},
	)
}

func (h getUsersHandler) Handle(ctx context.Context, gu GetUsers) ([]*user.User, error) {
	return h.userRepo.GetUsers(ctx, gu.Offset, gu.Limit)
}
