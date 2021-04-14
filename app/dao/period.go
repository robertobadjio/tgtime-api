package dao

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"officetime-api/app/model"
	"time"
)

func GetAllPeriods(w http.ResponseWriter, r *http.Request) {
	var periodsStruct model.Periods
	periodsStruct.Periods = model.GetAllPeriods()
	json.NewEncoder(w).Encode(periodsStruct)
}

func CreatePeriod(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the mac address and seconds only in order to update") // TODO: описание
	}

	var period model.Period
	err = json.Unmarshal(reqBody, &period)
	if err != nil {
		panic(err)
	}

	moscowLocation, _ := time.LoadLocation("Europe/Moscow")
	now := time.Now().In(moscowLocation)

	_, err = Db.Exec("INSERT INTO period (name, description, year, begin_at, ended_at, created_at) VALUES ($1, $2, $3, $4, $5, $6)", period.Name, period.Name, period.Year, period.BeginDate, period.EndDate, now.Format("2006-01-02 15:04:05"))
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(period)
	if err != nil {
		panic(err)
	}
}

func GetPeriod(w http.ResponseWriter, r *http.Request) {
	periodId := mux.Vars(r)["id"]

	var period model.Period
	row := Db.QueryRow("SELECT p.id, p.name, p.year, p.begin_at, p.ended_at FROM period p WHERE p.id = $1", periodId)
	err := row.Scan(&period.Id, &period.Name, &period.Year, &period.BeginDate, &period.EndDate)
	if err != nil {
		panic(err)
	}

	// TODO: Костыль 2020-06-01T00:00:00Z -> 2020-06-01
	timeTemp, _ := time.Parse(time.RFC3339, period.BeginDate)
	period.BeginDate = timeTemp.Format("2006-01-02")
	timeTemp, _ = time.Parse(time.RFC3339, period.EndDate)
	period.EndDate = timeTemp.Format("2006-01-02")

	err = json.NewEncoder(w).Encode(period)
	if err != nil {
		panic(err)
	}
}

func UpdatePeriod(w http.ResponseWriter, r *http.Request) {
	periodId := mux.Vars(r)["id"]
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the mac address and seconds only in order to update")
	}

	var period model.Period
	err = json.Unmarshal(reqBody, &period)
	if err != nil {
		panic(err)
	}

	_, err = Db.Exec("UPDATE period SET name = $1, year = $2, begin_at = $3, ended_at = $4 WHERE id = $5",
		period.Name, period.Year, period.BeginDate, period.EndDate, periodId)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(period)
	if err != nil {
		panic(err)
	}
}

func DeletePeriod(w http.ResponseWriter, r *http.Request) {
	periodId := mux.Vars(r)["id"]

	var period model.Period
	Db.QueryRow("UPDATE period SET deleted = true WHERE id = $1", periodId) // TODO: const

	err := json.NewEncoder(w).Encode(period)
	if err != nil {
		panic(err)
	}
}
