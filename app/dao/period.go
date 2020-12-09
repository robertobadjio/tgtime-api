package dao

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"time"
)

func GetAllPeriods(w http.ResponseWriter, r *http.Request) {
	rows, err := Db.Query("SELECT p.id, p.name, p.year, p.begin_at, p.ended_at FROM period p WHERE p.deleted = false") // TODO: const
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	periods := make([]*Period, 0)
	for rows.Next() {
		period := new(Period)
		err := rows.Scan(&period.Id, &period.Name, &period.Year, &period.BeginDate, &period.EndDate)
		if err != nil {
			panic(err)
		}

		periods = append(periods, period)
	}

	// TODO: Костыль 2020-06-01T00:00:00Z -> 2020-06-01
	for key, period := range periods {
		timeTemp, _ := time.Parse(time.RFC3339, period.BeginDate)
		periods[key].BeginDate = timeTemp.Format("2006-01-02")
		timeTemp, _ = time.Parse(time.RFC3339, period.EndDate)
		periods[key].EndDate = timeTemp.Format("2006-01-02")
	}

	var periodsStruct Periods
	periodsStruct.Periods = periods
	json.NewEncoder(w).Encode(periodsStruct)
}

func CreatePeriod(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the mac address and seconds only in order to update") // TODO: описание
	}

	var period Period
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

	var period Period
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

	var period Period
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

	var period Period
	Db.QueryRow("UPDATE period SET deleted = true WHERE id = $1", periodId) // TODO: const

	err := json.NewEncoder(w).Encode(period)
	if err != nil {
		panic(err)
	}
}
