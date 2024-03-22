package dao

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"io/ioutil"
	"net/http"
	"officetime-api/app/model"
	"strconv"
	"time"
)

var Db *sql.DB

func GetTimeDayAll(w http.ResponseWriter, r *http.Request) {
	user := mux.Vars(r)["id"]
	date := mux.Vars(r)["date"]
	if !checkDate("2006-01-02", date) {
		fmt.Fprintf(w, "Invalid date")
		return
	}

	userId, _ := strconv.Atoi(user)

	err := json.NewEncoder(w).Encode(model.GetTimeDayAll(userId, buildTimeStructFromDate(date)))
	if err != nil {
		panic(err)
	}
}

func GetTimeByPeriod(w http.ResponseWriter, r *http.Request) {
	user := mux.Vars(r)["id"]
	period := mux.Vars(r)["period"]

	userId, err := strconv.Atoi(user)
	if err != nil {
		panic(err)
	}
	periodId, err := strconv.Atoi(period)
	if err != nil {
		panic(err)
	}

	json.NewEncoder(w).Encode(model.GetTimeByPeriod(userId, periodId))
}

func CreateTime(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the mac address and seconds only in order to update") // TODO: Описание
	}

	var timeUser model.TimeUser
	err = json.Unmarshal(reqBody, &timeUser)
	if err != nil {
		panic(err)
	}

	_, err = Db.Exec("INSERT INTO time (mac_address, second, router_id) VALUES ($1, $2, $3)", timeUser.MacAddress, timeUser.Second, timeUser.RouterId)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(timeUser)
	if err != nil {
		panic(err)
	}
}

func buildTimeStructFromDate(date string) time.Time {
	timeStruct, err := time.Parse("2006-01-02", date)
	if err != nil {
		panic(err)
	}

	return timeStruct
}

/*func GetAllBreaksByTimes(macAddress, date string) []*model.Break {
	var breaksJson string
	row := Db.QueryRow("SELECT ts.breaks FROM time_summary ts WHERE ts.mac_address = $1 AND ts.date = $2", macAddress, date)
	err := row.Scan(&breaksJson)
	var s []*model.Break // TODO: Переименовать переменную
	if err == sql.ErrNoRows {
		return s
	}
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal([]byte(breaksJson), &s)
	if err != nil {
		panic(err)
	}
	return s
}*/

/*func GetDayTime(macAddress string, date string, sort string) int64 {
	var second int64
	// TODO: Разделить на отдельные методы
	var field string
	if "ASC" == sort {
		field = "ts.seconds_begin"
	} else {
		field = "ts.seconds_end"
	}

	row := Db.QueryRow("SELECT "+field+" FROM time_summary ts WHERE ts.mac_address = $1 AND ts.date = $2", macAddress, date)
	err := row.Scan(&second)

	if err == sql.ErrNoRows {
		return 0
	}

	if err != nil {
		panic(err)
	}

	return second
}*/
