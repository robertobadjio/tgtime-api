package model

import (
	"fmt"
	"github.com/lib/pq"
	"log"
	"strings"
	"time"
)

type Router struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Address     string `json:"address"`
	Login       string `json:"login"`
	Password    string `json:"password"`
	Status      bool   `json:"status"`
	WorkTime    bool   `json:"workTime"`
}

type NotFoundRouter struct {
	routerId int
}

type DuplicateAddress struct {
	address string
}

type DuplicateName struct {
	name string
}

func (e *NotFoundRouter) Error() string {
	return fmt.Sprintf("Router with id %d not found", e.routerId)
}

func (e *DuplicateAddress) Error() string {
	return fmt.Sprintf("Router with address %s already exists", e.address)
}

func (e *DuplicateName) Error() string {
	return fmt.Sprintf("Router with name %s already exists", e.name)
}

func GetAllRouters() []*Router {
	rows, err := Db.Query("SELECT id, name, description, address, login, password, status, work_time FROM router")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	routers := make([]*Router, 0)
	for rows.Next() {
		router := new(Router)
		err := rows.Scan(&router.Id, &router.Name, &router.Description, &router.Address, &router.Login, &router.Password, &router.Status, &router.WorkTime)
		if err != nil {
			log.Fatal(err)
		}

		routers = append(routers, router)
	}

	return routers
}

func GetRouter(routerId int) (*Router, error) {
	router := new(Router)
	row := Db.QueryRow("SELECT id, name, description, address, login, password, status, work_time FROM router WHERE id = $1", routerId)
	err := row.Scan(&router.Id, &router.Name, &router.Description, &router.Address, &router.Login, &router.Password, &router.Status, &router.WorkTime)
	if err != nil {
		return nil, &NotFoundRouter{routerId}
	}

	return router, nil
}

func CreateRouter(router Router) (int, error) {
	moscowLocation, _ := time.LoadLocation("Europe/Moscow")
	now := time.Now().In(moscowLocation)
	lastInsertId := 0
	err := Db.QueryRow("INSERT INTO router (name, description, address, login, password, status, work_time, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id", router.Name, router.Description, router.Address, router.Login, router.Password, router.Status, router.WorkTime, now.Format("2006-01-02 15:04:05")).Scan(&lastInsertId)

	if pgErr, ok := err.(*pq.Error); ok {
		if pgErr.Code == "23505" {
			if strings.Contains(err.Error(), "router_address_uindex") {
				return 0, &DuplicateAddress{router.Address}
			} else if strings.Contains(err.Error(), "router_name_uindex") {
				return 0, &DuplicateName{router.Name}
			} else {
				return 0, err
			}
		} else {
			return 0, err
		}
	}

	return lastInsertId, nil
}

func UpdateRouter(router Router) error {
	moscowLocation, _ := time.LoadLocation("Europe/Moscow")
	now := time.Now().In(moscowLocation)

	_, err := Db.Exec(
		"UPDATE router SET name = $1, description = $2, address = $3, login = $4, password = $5, status = $6, work_time = $7, updated_at = $8 WHERE id = $9",
		router.Name, router.Description, router.Address, router.Login, router.Password, router.Status, router.WorkTime, now.Format("2006-01-02 15:04:05"), router.Id)
	if err != nil {
		return err
	}

	return nil
}

func DeleteRouter(routerId int) error {
	_, err := Db.Exec("DELETE FROM router WHERE id = $1", routerId)
	if err != nil {
		return err
	}

	return nil
}
