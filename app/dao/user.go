package dao

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Users struct {
	Users []*User `json:"users"`
}

type User struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	MacAddress string `json:"macAddress"`
	TelegramId int64  `json:"telegramId"`
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
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
	json.NewEncoder(w).Encode(usersStruct)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	userId := mux.Vars(r)["id"]

	user := new(User)
	row := Db.QueryRow("SELECT u.id, u.name, u.email, u.mac_address, u.telegram_id FROM users u WHERE u.id = $1", userId)
	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.MacAddress, &user.TelegramId)
	if err != nil {
		panic(err)
	}

	json.NewEncoder(w).Encode(user)
}
