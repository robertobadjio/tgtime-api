package model

import (
	"fmt"
	"time"
)

type Period struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Year      int    `json:"year"`
	BeginDate string `json:"beginDate"`
	EndDate   string `json:"endDate"`
}

type Periods struct {
	Periods []*Period `json:"periods"`
}

type NotFoundPeriod struct {
	periodId int
}

func GetAllPeriods() []*Period {
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

	return periods
}

func GetPeriodById(periodId int) (*Period, error) {
	period := new(Period)
	row := Db.QueryRow("SELECT id, name, year, begin_at, ended_at FROM period WHERE id = $1", periodId)
	err := row.Scan(&period.Id, &period.Name, &period.Year, &period.BeginDate, &period.EndDate)
	if err != nil {
		return nil, &NotFoundUserOfId{int(periodId)}
	}

	return period, nil
}

func (e *NotFoundPeriod) Error() string {
	return fmt.Sprintf("Period with id %d not found", e.periodId)
}
