package model

import (
	"log"
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

func UpdateUser(user User) {
	_, err := Db.Exec(
		"UPDATE users SET name = $1, email = $2, mac_address = $3, telegram_id = $4 WHERE id = $5",
		user.Name, user.Email, user.MacAddress, user.TelegramId, user.Id)
	if err != nil {
		panic(err)
	}
}
