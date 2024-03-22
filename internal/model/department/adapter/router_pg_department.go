package adapter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"officetime-api/internal/db"
	"officetime-api/internal/model/department/domain/department"
	"strings"
	"time"
)

type PgDepartmentRepository struct {
	db *sql.DB
}

func NewPgDepartmentRepository(db *sql.DB) *PgDepartmentRepository {
	if db == nil {
		panic("missing db")
	}

	return &PgDepartmentRepository{db: db}
}

func (prr PgDepartmentRepository) CreateDepartment(ctx context.Context, department *department.Department) (*department.Department, error) {
	moscowLocation, _ := time.LoadLocation("Europe/Moscow")
	now := time.Now().In(moscowLocation)

	_, err := prr.db.Exec("INSERT INTO department (name, description, created_at, updated_at, deleted) VALUES ($1, $2, $3, $4, $5)", department.Name, department.Description, now.Format("2006-01-02 15:04:05"), now.Format("2006-01-02 15:04:05"), false)
	if err != nil {
		return nil, err
	}

	return department, nil
}

func (prr PgDepartmentRepository) UpdateDepartment(_ context.Context, d *department.Department) (*department.Department, error) {
	_, err := prr.db.Exec(
		"UPDATE department SET name = $1, description = $2  WHERE id = $3",
		d.Name, d.Description, d.Id)
	if err != nil {
		return nil, &department.ErrorUpdateDepartment{DepartmentId: d.Id}
	}

	return d, nil
}

func (prr PgDepartmentRepository) GetDepartment(_ context.Context, departmentId int) (*department.Department, error) {
	var d department.Department
	row := prr.db.QueryRow("SELECT d.id, d.name, d.description FROM department d WHERE d.id = $1", departmentId)
	err := row.Scan(&d.Id, &d.Name, &d.Description)
	if err != nil {
		panic(err)
	}

	return &d, nil
}

func (prr PgDepartmentRepository) DeleteDepartment(_ context.Context, departmentId int) error {
	_, err := prr.db.Exec("UPDATE department SET deleted = true WHERE id = $1", departmentId)
	if err != nil {
		return &department.ErrorDeleteDepartment{DepartmentId: departmentId}
	}

	return nil
}

func (prr PgDepartmentRepository) GetDepartments(_ context.Context) ([]*department.Department, error) {
	rows, err := prr.db.Query("SELECT id, name, description FROM department")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	departments := make([]*department.Department, 0)
	for rows.Next() {
		d := new(department.Department)
		err := rows.Scan(&d.Id, &d.Name, &d.Description)
		if err != nil {
			log.Fatal(err)
		}

		departments = append(departments, d)
	}

	return departments, nil
}

func (prr PgDepartmentRepository) existsDepartmentByFilter(filters []string) (bool, error) {
	var id int
	err := prr.db.QueryRow(
		"SELECT id FROM department WHERE " + buildOrCondition(filters),
	).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// a real error happened! you should change your function return
			// to "(bool, error)" and return "false, err" here
			return false, nil
		}

		return false, fmt.Errorf("error getting department from db: %v", err)
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
