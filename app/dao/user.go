package dao

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"officetime-api/app/model"
	"officetime-api/app/service"
	"strconv"
)

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	json.NewEncoder(w).Encode(model.GetAllUsers(offset, limit))
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	userId, _ := strconv.Atoi(id)
	user, _ := model.GetUser(int64(userId))
	json.NewEncoder(w).Encode(user)
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

	service.SendEmail(user.Name, user.Email, password)

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
