package dao

import (
	"encoding/json"
	"github.com/gorilla/mux"
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
