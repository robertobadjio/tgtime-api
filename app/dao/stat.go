package dao

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"officetime-api/app/model"
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

func checkDate(format, date string) bool {
	t, err := time.Parse(format, date)
	if err != nil {
		return false
	}
	return t.Format(format) == date
}