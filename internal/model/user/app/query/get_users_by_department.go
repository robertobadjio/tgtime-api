package query

import (
	"context"
	"officetime-api/internal/common/decorator"
	"officetime-api/internal/model/user/domain/user"
)

type GetUsersByDepartment struct {
	DepartmentId int
}

type GetUsersByDepartmentHandler decorator.QueryHandler[GetUsersByDepartment, []*user.User]

type getUsersByDepartmentHandler struct {
	userRepo user.Repository
}

func NewGetUsersByDepartmentHandler(userRepo user.Repository) GetUsersByDepartmentHandler {
	return decorator.ApplyQueryDecorators[GetUsersByDepartment, []*user.User](
		getUsersByDepartmentHandler{userRepo: userRepo},
	)
}

func (h getUsersByDepartmentHandler) Handle(ctx context.Context, qr GetUsersByDepartment) ([]*user.User, error) {
	return h.userRepo.GetUsersByDepartment(ctx, qr.DepartmentId)
}
