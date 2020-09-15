package dao

import (
	"encoding/json"
	"net/http"
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

	var periodsStruct Periods
	periodsStruct.Periods = periods
	json.NewEncoder(w).Encode(periodsStruct)
}
