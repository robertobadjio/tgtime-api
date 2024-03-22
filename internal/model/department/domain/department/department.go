package department

import "fmt"

type Department struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ErrorDeleteDepartment struct {
	DepartmentId int
}

type ErrorUpdateDepartment struct {
	DepartmentId int
}

func (e *ErrorUpdateDepartment) Error() string {
	return fmt.Sprintf("Error update department %d", e.DepartmentId)
}

func (e *ErrorDeleteDepartment) Error() string {
	return fmt.Sprintf("Errod deleting department %d", e.DepartmentId)
}
