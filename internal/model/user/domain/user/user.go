package user

import "fmt"

type User struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Lastname   string `json:"lastname"`
	BirthDate  string `json:"birthDate"`
	Email      string `json:"email"`
	MacAddress string `json:"macAddress"`
	TelegramId int64  `json:"telegramId"`
	Role       string `json:"role"`
	Department int64  `json:"department"`
	Position   string `json:"position"`
}

type Users struct {
	Users []*User `json:"users"`
}

type UserAlreadyExists struct {
	Email string
}

type TelegramAlreadyExists struct {
	TgId int64
}

type MacAddressAlreadyExists struct {
	MacAddress string
}

type ErrorDeleteUser struct {
	UserId int
}

type NotFoundUser struct {
	UserEmail string
}

type NotFoundUserByMacAddress struct {
	MacAddress string
}

type NotFoundUserOfId struct {
	UserId int
}

func NewUser() (*User, error) {
	return &User{}, nil // TODO: !
}

func (e *UserAlreadyExists) Error() string {
	return fmt.Sprintf("User with email %s already register", e.Email)
}

func (e *TelegramAlreadyExists) Error() string {
	return fmt.Sprintf("User with telegram id %d already register", e.TgId)
}

func (e *MacAddressAlreadyExists) Error() string {
	return fmt.Sprintf("User with mac-address %s already register", e.MacAddress)
}

func (e *NotFoundUser) Error() string {
	return fmt.Sprintf("User with email %s not found", e.UserEmail)
}

func (e *NotFoundUserByMacAddress) Error() string {
	return fmt.Sprintf("User with mac address %s not found", e.MacAddress)
}

func (e *NotFoundUserOfId) Error() string {
	return fmt.Sprintf("User with id %d not found", e.UserId)
}

func (e *ErrorDeleteUser) Error() string {
	return fmt.Sprintf("User with id %d not deleted", e.UserId)
}
