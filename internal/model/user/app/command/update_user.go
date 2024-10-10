package command

import (
	"context"
	"github.com/robertobadjio/tgtime-api/internal/common/decorator"
	"github.com/robertobadjio/tgtime-api/internal/model/user/domain/user"
)

type UpdateUser struct {
	User *user.User
}

type UpdateUserHandler decorator.CommandHandler[UpdateUser]

type updateUserHandler struct {
	userRepository user.Repository
}

func NewUpdateUserHandler(userRepository user.Repository) UpdateUserHandler {
	if userRepository == nil {
		panic("nil userRepository")
	}

	return decorator.ApplyCommandDecorators[UpdateUser](
		updateUserHandler{userRepository: userRepository},
	)
}

func (h updateUserHandler) Handle(ctx context.Context, cmd UpdateUser) error {
	_, err := h.userRepository.UpdateUser(ctx, cmd.User) // TODO: !
	if err != nil {
		return err
	}

	return nil
}
