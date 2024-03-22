package adapter

import (
	"context"
	"database/sql"
	"time"
)

type PgWeekendRepository struct {
	db *sql.DB
}

func NewPgWeekendRepository(db *sql.DB) *PgWeekendRepository {
	if db == nil {
		panic("missing db")
	}

	return &PgWeekendRepository{db: db}
}

func (prr PgWeekendRepository) GetWeekends(_ context.Context) ([]string, error) {
	rows, err := prr.db.Query("SELECT date FROM weekend")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var weekends []string

	for rows.Next() {
		var date string
		err := rows.Scan(&date)
		if err != nil {
			panic(err)
		}
		dateObject, err := time.Parse(time.RFC3339, date)
		if err != nil {
			panic(err)
		}

		weekends = append(weekends, dateObject.Format("2006-01-02"))
	}

	return weekends, nil
}

func (prr PgWeekendRepository) GetWeekendsByPeriod(_ context.Context, begin, end time.Time) (map[string]bool, error) {
	rows, err := prr.db.Query("SELECT date FROM weekend WHERE date BETWEEN $1 AND $2", begin.Format("2006-01-02"), end.Format("2006-01-02"))
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	weekend := make(map[string]bool)

	for rows.Next() {
		var date string
		err := rows.Scan(&date)
		if err != nil {
			panic(err)
		}
		dateObject, err := time.Parse(time.RFC3339, date)
		if err != nil {
			panic(err)
		}
		weekend[dateObject.Format("2006-01-02")] = true
	}

	return weekend, nil
}
