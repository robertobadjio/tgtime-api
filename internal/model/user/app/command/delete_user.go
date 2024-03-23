package command

import (
	"context"
	"officetime-api/internal/common/decorator"
	"officetime-api/internal/model/user/domain/user"
)

type DeleteUser struct {
	UserId int
}

type DeleteUserHandler decorator.CommandHandler[DeleteUser]

type deleteUserHandler struct {
	userRepository user.Repository
}

func NewDeleteUserHandler(userRepository user.Repository) DeleteUserHandler {
	if userRepository == nil {
		panic("nil userRepository")
	}

	return decorator.ApplyCommandDecorators[DeleteUser](
		deleteUserHandler{userRepository: userRepository},
	)
}

func (h deleteUserHandler) Handle(ctx context.Context, cmd DeleteUser) error {
	err := h.userRepository.DeleteUser(ctx, cmd.UserId)
	if err != nil {
		return err
	}

	return nil
}
