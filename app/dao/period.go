package dao

import (
	"encoding/json"
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
