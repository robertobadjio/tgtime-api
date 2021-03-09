package model

import (
	"fmt"
	"log"
)

type Department struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ErrorDeleteDepartment struct {
	departmentId int
}

type ErrorUpdateDepartment struct {
	departmentId int
}

func GetAllDepartments() []*Department {
	rows, err := Db.Query("SELECT id, name, description FROM department")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	departments := make([]*Department, 0)
	for rows.Next() {
		department := new(Department)
		err := rows.Scan(&department.Id, &department.Name, &department.Description)
		if err != nil {
			log.Fatal(err)
		}

		departments = append(departments, department)
	}

	return departments
}

func UpdateDepartment(department Department) error {
	_, err := Db.Exec(
		"UPDATE department SET name = $1, description = $2  WHERE id = $3",
		department.Name, department.Description, department.Id)
	if err != nil {
		return &ErrorUpdateDepartment{department.Id}
	}

	return nil
}

func DeleteDepartment(departmentId int) error {
	_, err := Db.Exec("UPDATE department SET deleted = true WHERE id = $1", departmentId)
	if err != nil {
		return &ErrorDeleteDepartment{departmentId}
	}

	return nil
}

func (e *ErrorUpdateDepartment) Error() string {
	return fmt.Sprintf("Error update department %d", e.departmentId)
}

func (e *ErrorDeleteDepartment) Error() string {
	return fmt.Sprintf("Errod deleting department %d", e.departmentId)
}
