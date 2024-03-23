package command_query

import (
	"officetime-api/internal/model/user/domain/user"
)

type User struct {
	Password string
	User     *user.User
}
