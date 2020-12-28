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

func GetAllDepartments(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(model.GetAllDepartments())
}

func DeleteDepartment(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	userId, _ := strconv.Atoi(id)

	err := model.DeleteDepartment(userId)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	json.NewEncoder(w).Encode(userId)
}

func UpdateDepartment(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the mac address and seconds only in order to update")
	}

	var department model.Department
	err = json.Unmarshal(reqBody, &department)
	if err != nil {
		panic(err)
	}

	model.UpdateDepartment(department)

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(department)
	if err != nil {
		panic(err)
	}
}
