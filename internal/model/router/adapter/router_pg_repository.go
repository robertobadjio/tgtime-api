package adapter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"officetime-api/internal/db"
	"officetime-api/internal/model/router/domain/router"
	"strings"
	"time"
)

type PgRouterRepository struct {
	db *sql.DB
}

func NewPgRouterRepository(db *sql.DB) *PgRouterRepository {
	if db == nil {
		panic("missing db")
	}

	return &PgRouterRepository{db: db}
}

func (prr PgRouterRepository) CreateRouter(ctx context.Context, router *router.Router) (*router.Router, error) {
	existingRouters, err := prr.existsRouterByFilter([]string{"name = '" + router.Name + "'", "address = '" + router.Address + "'"})
	if existingRouters {
		return nil, fmt.Errorf("router already exists")
	}

	moscowLocation, _ := time.LoadLocation("Europe/Moscow")
	now := time.Now().In(moscowLocation)
	lastInsertId := 0
	err = db.GetDB().QueryRow("INSERT INTO router (name, description, address, login, password, status, work_time, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id", router.Name, router.Description, router.Address, router.Login, router.Password, router.Status, router.WorkTime, now.Format("2006-01-02 15:04:05")).Scan(&lastInsertId)
	if err != nil {
		return nil, fmt.Errorf("error inserting router into db: %v", err)
	}

	if pgErr, ok := err.(*pq.Error); ok {
		if pgErr.Code == "23505" {
			if strings.Contains(err.Error(), "router_address_uindex") {
				return nil, fmt.Errorf("router with address %s already exists", router.Address)
			} else if strings.Contains(err.Error(), "router_name_uindex") {
				return nil, fmt.Errorf("router with name %s already exists", router.Name)
			} else {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return prr.GetRouter(ctx, lastInsertId)
}

func (prr PgRouterRepository) UpdateRouter(_ context.Context, router *router.Router) (*router.Router, error) {
	moscowLocation, _ := time.LoadLocation("Europe/Moscow")
	now := time.Now().In(moscowLocation)

	_, err := db.GetDB().Exec(
		"UPDATE router SET name = $1, description = $2, address = $3, login = $4, password = $5, status = $6, work_time = $7, updated_at = $8 WHERE id = $9",
		router.Name, router.Description, router.Address, router.Login, router.Password, router.Status, router.WorkTime, now.Format("2006-01-02 15:04:05"), router.Id)
	if err != nil {
		return nil, err
	}

	return router, nil
}

func (prr PgRouterRepository) GetRouter(_ context.Context, routerId int) (*router.Router, error) {
	router := new(router.Router)
	if err := db.GetDB().QueryRow(
		"SELECT id, name, description, address, login, password, status, work_time FROM router WHERE id = $1",
		routerId,
	).Scan(&router.Id, &router.Name, &router.Description, &router.Address, &router.Login, &router.Password, &router.Status, &router.WorkTime); err == nil {
		return router, nil
	} else if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else {
		return nil, fmt.Errorf("error getting router from db: %v", err)
	}
}

func (prr PgRouterRepository) DeleteRouter(_ context.Context, routerId int) error {
	_, err := db.GetDB().Exec("DELETE FROM router WHERE id = $1", routerId)
	if err != nil {
		return err
	}

	return nil
}

func (prr PgRouterRepository) GetRouters(_ context.Context) ([]*router.Router, error) {
	rows, err := db.GetDB().Query("SELECT id, name, description, address, login, password, status, work_time FROM router")
	if err != nil {
		return nil, fmt.Errorf("error getting router: %v", err)
	}
	defer rows.Close()

	routers := make([]*router.Router, 0)
	for rows.Next() {
		router := new(router.Router)
		err = rows.Scan(&router.Id, &router.Name, &router.Description, &router.Address, &router.Login, &router.Password, &router.Status, &router.WorkTime)
		if err != nil {
			return nil, fmt.Errorf("error scan router: %v", err)
		}

		routers = append(routers, router)
	}

	return routers, nil
}

func (prr PgRouterRepository) existsRouterByFilter(filters []string) (bool, error) {
	var id int
	err := db.GetDB().QueryRow(
		"SELECT id FROM router WHERE " + buildOrCondition(filters),
	).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// a real error happened! you should change your function return
			// to "(bool, error)" and return "false, err" here
			return false, nil
		}

		return false, fmt.Errorf("error getting router from db: %v", err)
	}

	return true, nil
}

func buildOrCondition(filters []string) string {
	return strings.Join(filters, " OR ")
}

func NewPgConnection() (*sql.DB, error) {
	/*config := mysql.NewConfig()

	config.Net = "tcp"
	config.Addr = os.Getenv("MYSQL_ADDR")
	config.User = os.Getenv("MYSQL_USER")
	config.Passwd = os.Getenv("MYSQL_PASSWORD")
	config.DBName = os.Getenv("MYSQL_DATABASE")
	config.ParseTime = true // with that parameter, we can use time.Time in mysqlHour.Hour

	db, err := sqlx.Connect("mysql", config.FormatDSN())
	if err != nil {
		return nil, errors.Wrap(err, "cannot connect to MySQL")
	}

	return db, nil*/

	return db.GetDB(), nil
}
