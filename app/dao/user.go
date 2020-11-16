package dao

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
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
