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

type Data struct {
	Content []*model.Router `json:"content"`
}

func GetAllRouters(w http.ResponseWriter, r *http.Request) {
	data := new(Data)
	data.Content = model.GetAllRouters()
	json.NewEncoder(w).Encode(data)
}

func GetRouter(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	routerId, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Param not number")
		return
	}

	router, err := model.GetRouter(routerId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, err.Error())
		return
	}
	json.NewEncoder(w).Encode(router)
}

func CreateRouter(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error body")
		return
	}

	var router model.Router
	err = json.Unmarshal(reqBody, &router)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, err.Error())
		return
	}

	var routerId int
	routerId, err = model.CreateRouter(router)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		fmt.Fprintf(w, err.Error())
		return
	}

	router.Id = routerId
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(router)
	if err != nil {
		panic(err)
	}
}

func UpdateRouter(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	routerId, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error id router")
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error body")
		return
	}

	var router model.Router
	err = json.Unmarshal(reqBody, &router)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		fmt.Fprintf(w, err.Error())
		return
	}
	router.Id = routerId

	err = model.UpdateRouter(router)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		fmt.Fprintf(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(router)
	if err != nil {
		panic(err)
	}
}

func DeleteRouter(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	routerId, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error id router")
		return
	}

	err = model.DeleteRouter(routerId)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, err.Error())
		return
	}

	json.NewEncoder(w).Encode(routerId)
}
