package api

import (
	"context"
	"github.com/robertobadjio/tgtime-api/app/service"
	"github.com/robertobadjio/tgtime-api/internal/model/user/app/command"
	"github.com/robertobadjio/tgtime-api/internal/model/user/app/command_query"
	"github.com/robertobadjio/tgtime-api/internal/model/user/app/query"
	"github.com/robertobadjio/tgtime-api/internal/model/user/domain/user"
)

func (s *apiService) CreateUser(ctx context.Context, u *user.User) (*user.User, error) {
	cmd := command_query.CreateUser{User: u}
	userCQ, err := s.userApp.CommandsQueries.CreateUser.Handle(ctx, cmd)
	if err != nil {
		return nil, err
	}

	service.SendEmail(userCQ.User.Name, userCQ.User.Email, userCQ.Password)

	return userCQ.User, nil
}

func (s *apiService) GetUsers(ctx context.Context, offset, limit int) ([]*user.User, error) {
	// TODO: Валидация offset и limit
	qr := query.GetUsers{Offset: offset, Limit: limit}
	users, err := s.userApp.Queries.GetUsers.Handle(ctx, qr)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *apiService) GetUser(ctx context.Context, userId int) (*user.User, error) {
	qr := query.GetUser{UserId: userId}
	u, err := s.userApp.Queries.GetUser.Handle(ctx, qr)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *apiService) GetUserByMacAddress(ctx context.Context, macAddress string) (*user.User, error) {
	qr := query.GetUserByMacAddress{MacAddress: macAddress}
	u, err := s.userApp.Queries.GetUserByMacAddress.Handle(ctx, qr)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *apiService) GetUserByTelegramId(ctx context.Context, telegramId int64) (*user.User, error) {
	qr := query.GetUserByTelegramId{TelegramId: telegramId}
	u, err := s.userApp.Queries.GetUserByTelegramId.Handle(ctx, qr)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *apiService) UpdateUser(ctx context.Context, u *user.User) (*user.User, error) {
	cmd := command.UpdateUser{User: u}
	err := s.userApp.Commands.UpdateUser.Handle(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *apiService) DeleteUser(ctx context.Context, userId int) error {
	cmd := command.DeleteUser{UserId: userId}
	err := s.userApp.Commands.DeleteUser.Handle(ctx, cmd)
	if err != nil {
		return err
	}

	return nil
}
