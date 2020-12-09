package model

import (
	"fmt"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
	"strings"
	"time"
)

type User struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	MacAddress string `json:"macAddress"`
	TelegramId int64  `json:"telegramId"`
}

type Users struct {
	Users []*User `json:"users"`
}

type UserAlreadyExists struct {
	email string
}

type TelegramAlreadyExists struct {
	tgId int64
}

type MacAddressAlreadyExists struct {
	macAddress string
}

type ErrorDeleteUser struct {
	userId int
}

func GetAllUsers() Users {
	rows, err := Db.Query("SELECT u.id, u.name, u.email, u.mac_address, u.telegram_id FROM users u")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	users := make([]*User, 0)
	for rows.Next() {
		user := new(User)
		err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.MacAddress, &user.TelegramId)
		if err != nil {
			log.Fatal(err)
		}

		users = append(users, user)
	}

	var usersStruct Users
	usersStruct.Users = users

	return usersStruct
}

func GetUser(userId int64) *User {
	user := new(User)
	row := Db.QueryRow("SELECT u.id, u.name, u.email, u.mac_address, u.telegram_id FROM users u WHERE u.id = $1", userId)
	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.MacAddress, &user.TelegramId)
	if err != nil {
		panic(err)
	}

	return user
}

func DeleteUser(userId int) error {
	_, err := Db.Exec("UPDATE users SET deleted = true WHERE id = $1", userId)
	if err != nil {
		return &ErrorDeleteUser{userId}
	}

	return nil
}

func UpdateUser(user User) {
	_, err := Db.Exec(
		"UPDATE users SET name = $1, email = $2, mac_address = $3, telegram_id = $4 WHERE id = $5",
		user.Name, user.Email, user.MacAddress, user.TelegramId, user.Id)
	if err != nil {
		panic(err)
	}
}

func CreateUser(user User) (string, int, error) {
	moscowLocation, _ := time.LoadLocation("Europe/Moscow")
	now := time.Now().In(moscowLocation)
	password := randomString(10)
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	lastInsertId := 0
	err := Db.QueryRow("INSERT INTO users (name, email, mac_address, telegram_id, password, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", user.Name, user.Email, user.MacAddress, user.TelegramId, passwordHash, now.Format("2006-01-02 15:04:05")).Scan(&lastInsertId)

	if pgerr, ok := err.(*pq.Error); ok {
		if pgerr.Code == "23505" {
			if strings.Contains(err.Error(), "users_email_unique") {
				return "", 0, &UserAlreadyExists{user.Email}
			} else if strings.Contains(err.Error(), "users_telegram_id_unique") {
				return "", 0, &TelegramAlreadyExists{user.TelegramId}
			} else if strings.Contains(err.Error(), "users_mac_address_unique") {
				return "", 0, &MacAddressAlreadyExists{user.MacAddress}
			} else {
				return "", 0, err
			}
		} else {
			return "", 0, err
		}
	}

	return password, lastInsertId, nil
}

func (e *UserAlreadyExists) Error() string {
	return fmt.Sprintf("User with email %s already register", e.email)
}

func (e *TelegramAlreadyExists) Error() string {
	return fmt.Sprintf("User with telegram id %d already register", e.tgId)
}

func (e *MacAddressAlreadyExists) Error() string {
	return fmt.Sprintf("User with mac-address %s already register", e.macAddress)
}

func (e *ErrorDeleteUser) Error() string {
	return fmt.Sprintf("User with id %d not deleted", e.userId)
}

func randomString(n int) string {
	rand.Seed(time.Now().UnixNano())
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
