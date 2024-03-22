package adapter

import (
	"context"
	"database/sql"
	"fmt"
	"officetime-api/internal/model/period/domain/period"
	"time"
)

type PgPeriodRepository struct {
	db *sql.DB
}

type NotFoundPeriod struct {
	periodId int
}

func NewPgPeriodRepository(db *sql.DB) *PgPeriodRepository {
	if db == nil {
		panic("missing db")
	}

	return &PgPeriodRepository{db: db}
}

func (r PgPeriodRepository) GetPeriods(ctx context.Context) ([]*period.Period, error) {
	rows, err := r.db.Query("SELECT p.id, p.name, p.year, p.begin_at, p.ended_at FROM period p WHERE p.deleted = false") // TODO: const
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	periods := make([]*period.Period, 0)
	for rows.Next() {
		period := new(period.Period)
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

	return periods, nil
}

func (r PgPeriodRepository) GetPeriod(_ context.Context, periodId int) (*period.Period, error) {
	periodNew := new(period.Period)
	row := r.db.QueryRow("SELECT id, name, year, begin_at, ended_at FROM period WHERE id = $1", periodId)
	err := row.Scan(&periodNew.Id, &periodNew.Name, &periodNew.Year, &periodNew.BeginDate, &periodNew.EndDate)
	if err != nil {
		return nil, fmt.Errorf("period not found")
	}

	return periodNew, nil
}

func (r PgPeriodRepository) CreatePeriod(_ context.Context, period *period.Period) (*period.Period, error) {
	moscowLocation, _ := time.LoadLocation("Europe/Moscow")
	now := time.Now().In(moscowLocation)

	_, err := r.db.Exec("INSERT INTO period (name, description, year, begin_at, ended_at, created_at) VALUES ($1, $2, $3, $4, $5, $6)", period.Name, period.Name, period.Year, period.BeginDate, period.EndDate, now.Format("2006-01-02 15:04:05"))
	if err != nil {
		panic(err)
	}

	return period, nil // TODO: !
}

func (r PgPeriodRepository) UpdatePeriod(_ context.Context, period *period.Period) (*period.Period, error) {
	_, err := r.db.Exec("UPDATE period SET name = $1, year = $2, begin_at = $3, ended_at = $4 WHERE id = $5",
		period.Name, period.Year, period.BeginDate, period.EndDate, period.Id)
	if err != nil {
		return nil, err
	}

	return period, nil // TODO: !
}

func (r PgPeriodRepository) DeletePeriod(ctx context.Context, periodId int) error {
	r.db.QueryRow("UPDATE period SET deleted = true WHERE id = $1", periodId) // TODO: !

	return nil
}

func (e *NotFoundPeriod) Error() string {
	return fmt.Sprintf("Period with id %d not found", e.periodId)
}
