package dao

import (
	"encoding/json"
	"log"
	"net/http"
)

func GetAllPeriods(w http.ResponseWriter, r *http.Request) {
	rows, err := Db.Query("SELECT p.id, p.begin_at, p.ended_at FROM period p")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	periods := make([]*Period, 0)
	for rows.Next() {
		period := new(Period)
		err := rows.Scan(&period.Id, &period.BeginDate, &period.EndDate)
		if err != nil {
			log.Fatal(err)
		}

		periods = append(periods, period)
	}

	json.NewEncoder(w).Encode(periods)
}
