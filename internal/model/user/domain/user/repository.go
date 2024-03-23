package user

import (
	"context"
	"officetime-api/internal/model/user/app/command_query"
)

type Repository interface {
	CreateUser(ctx context.Context, user *User) (*command_query.User, error)
	UpdateUser(ctx context.Context, user *User) (*User, error)
	GetUser(ctx context.Context, userId int) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUsersByDepartment(ctx context.Context, departmentId int) ([]*User, error)
	GetUserPasswordHashByEmail(ctx context.Context, email string) (string, error)
	GetUsers(ctx context.Context, offset, limit int) ([]*User, error)
	DeleteUser(ctx context.Context, userId int) error
}
