package department

import (
	"context"
)

type Repository interface {
	CreateDepartment(ctx context.Context, department *Department) (*Department, error)
	UpdateDepartment(ctx context.Context, department *Department) (*Department, error)
	GetDepartment(ctx context.Context, departmentId int) (*Department, error)
	GetDepartments(ctx context.Context) ([]*Department, error)
	DeleteDepartment(ctx context.Context, departmentId int) error
}
