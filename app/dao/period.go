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
	rows, err := Db.Query("SELECT p.id, p.name, p.year, p.begin_at, p.ended_at FROM period p")
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
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the mac address and seconds only in order to update") // TODO: описание
	}

	var period Period
	err = json.Unmarshal(reqBody, &period)
	if err != nil {
		panic(err)
	}

	_, err = Db.Exec("INSERT INTO period (name, year, begin_at, ended_at) VALUES ($1, $2, $3, $4)", period.Name, period.Year, period.BeginDate, period.EndDate)
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
	row := Db.QueryRow("SELECT p.name, p.year, p.begin_at, p.ended_at FROM period p WHERE p.id = $1", periodId)
	err := row.Scan(&period.Id, &period.Name, &period.Year, &period.BeginDate, &period.EndDate)
	if err != nil {
		panic(err)
	}

	err = json.NewEncoder(w).Encode(period)
	if err != nil {
		panic(err)
	}
}

func UpdatePeriod(w http.ResponseWriter, r *http.Request) {
	periodId := mux.Vars(r)["id"]

	var period Period
	Db.QueryRow("UPDATE period p SET p.name, p.year, p.begin_at, p.ended_at WHERE p.id = $1", periodId)

	err := json.NewEncoder(w).Encode(period)
	if err != nil {
		panic(err)
	}
}

func DeletePeriod(w http.ResponseWriter, r *http.Request) {
	periodId := mux.Vars(r)["id"]

	var period Period
	Db.QueryRow("UPDATE period p SET p.deleted = $1 WHERE p.id = $2", periodId, 1) // TODO: const

	err := json.NewEncoder(w).Encode(period)
	if err != nil {
		panic(err)
	}
}
