package dao

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"officetime-api/app/model"
	"strconv"
	"time"
)

func GetAllDepartments(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(model.GetAllDepartments())
}

func GetDepartment(w http.ResponseWriter, r *http.Request) {
	departmentId, _ := strconv.Atoi(mux.Vars(r)["id"])
	json.NewEncoder(w).Encode(model.GetDepartment(departmentId))
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

func CreateDepartment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the mac address and seconds only in order to update") // TODO: описание
	}

	var department model.Department
	err = json.Unmarshal(reqBody, &department)
	if err != nil {
		panic(err)
	}

	moscowLocation, _ := time.LoadLocation("Europe/Moscow")
	now := time.Now().In(moscowLocation)

	_, err = Db.Exec("INSERT INTO department (name, description, created_at, updated_at, deleted) VALUES ($1, $2, $3, $4, $5)", department.Name, department.Description, now.Format("2006-01-02 15:04:05"), now.Format("2006-01-02 15:04:05"), false)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(department)
	if err != nil {
		panic(err)
	}
}
