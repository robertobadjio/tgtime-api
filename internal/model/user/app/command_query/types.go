package command_query

import (
	"github.com/robertobadjio/tgtime-api/internal/model/user/domain/user"
)

type User struct {
	Password string
	User     *user.User
}
