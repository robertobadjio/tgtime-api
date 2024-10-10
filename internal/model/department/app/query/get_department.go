package query

import (
	"context"
	"github.com/robertobadjio/tgtime-api/internal/common/decorator"
	"github.com/robertobadjio/tgtime-api/internal/model/department/domain/department"
)

type GetDepartment struct {
	DepartmentId int
}

type GetDepartmentHandler decorator.QueryHandler[GetDepartment, *department.Department]

/*type GetDepartmentReadModel interface {
	GetDepartment(ctx context.Context, departmentId int) (*department.Department, error)
}*/

type getDepartmentHandler struct {
	departmentRepo department.Repository
}

func NewGetDepartmentHandler(departmentRepo department.Repository) GetDepartmentHandler {
	return decorator.ApplyQueryDecorators[GetDepartment, *department.Department](
		getDepartmentHandler{departmentRepo: departmentRepo},
	)
}

func (h getDepartmentHandler) Handle(ctx context.Context, qr GetDepartment) (*department.Department, error) {
	return h.departmentRepo.GetDepartment(ctx, qr.DepartmentId)
}
