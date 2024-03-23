package command_query

import (
	"context"
	"officetime-api/internal/common/decorator"
	"officetime-api/internal/model/user/domain/user"
)

type CreateUser struct {
	User *user.User
}

type CreateUserHandler decorator.CommandQueryHandler[CreateUser, User]

type createUserHandler struct {
	userRepository user.Repository
}

func NewCreateUserHandler(userRepository user.Repository) CreateUserHandler {
	if userRepository == nil {
		panic("nil userRepository")
	}

	return decorator.ApplyCommandQueryDecorators[CreateUser, User](
		createUserHandler{userRepository: userRepository},
	)
}

func (h createUserHandler) Handle(ctx context.Context, cmdQr CreateUser) (User, error) {
	password, userNew, err := h.userRepository.CreateUser(ctx, cmdQr.User)
	if err != nil {
		return User{Password: "", User: nil}, err
	}

	return User{Password: password, User: userNew}, nil
}
