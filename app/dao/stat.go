package dao

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"officetime-api/app/model"
	"strconv"
	"time"
)

func GetStatByPeriodsAndRouters(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(model.GetAllTimesByPeriodsAndRouters())
}

func GetAllTimesDepartmentsByDate(w http.ResponseWriter, r *http.Request) {
	date := mux.Vars(r)["date"]
	if !checkDate("2006-01-02", date) {
		fmt.Fprintf(w, "Invalid date")
		return
	}
	json.NewEncoder(w).Encode(model.GetAllTimesDepartmentsByDate(date))
}

func GetStatWorkingPeriod(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	period := mux.Vars(r)["period"]
	userId, _ := strconv.Atoi(id)
	periodId, _ := strconv.Atoi(period)
	json.NewEncoder(w).Encode(model.GetStatWorkingPeriod(userId, periodId))
}

func checkDate(format, date string) bool {
	t, err := time.Parse(format, date)
	if err != nil {
		return false
	}
	return t.Format(format) == date
}