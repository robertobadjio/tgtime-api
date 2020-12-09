package dao

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"officetime-api/app/config"
	"officetime-api/app/model"
	"strconv"
)

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(model.GetAllUsers())
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	userId, _ := strconv.Atoi(id)
	json.NewEncoder(w).Encode(model.GetUser(int64(userId)))
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the mac address and seconds only in order to update")
	}

	var user model.User
	err = json.Unmarshal(reqBody, &user)
	if err != nil {
		panic(err)
	}

	model.UpdateUser(user)

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		panic(err)
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the mac address and seconds only in order to update") // TODO: описание
		return
	}

	var user model.User
	err = json.Unmarshal(reqBody, &user)
	if err != nil {
		panic(err)
	}

	var userId int
	var password string
	password, userId, err = model.CreateUser(user)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	sendEmail(user.Name, user.Email, password)

	user.Id = userId
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		panic(err)
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	userId, _ := strconv.Atoi(id)

	err := model.DeleteUser(userId)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	json.NewEncoder(w).Encode(userId)
}

func sendEmail(name, email, password string) {
	auth := smtp.PlainAuth("", "info@officetime.tech", "E61kTuq7", "smtp.timeweb.ru")

	// Here we do it all: connect to our server, set up a message and send it
	to := []string{email}
	msg := []byte("To: " + email + "\r\n" +
		"Subject: OfficeTime registration\r\n" +
		"Для начала работы напишите '/start' телеграм-боту " + config.Config.TelegramBot + "\r\n" +
		"Hello " + name + "! You password: " + password + "\r\n")
	err := smtp.SendMail("smtp.timeweb.ru:25", auth, "info@officetime.tech", to, msg)
	if err != nil {
		log.Fatal(err)
	}
}
