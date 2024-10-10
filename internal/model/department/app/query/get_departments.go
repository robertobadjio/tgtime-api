package query

import (
	"context"
	"github.com/robertobadjio/tgtime-api/internal/common/decorator"
	"github.com/robertobadjio/tgtime-api/internal/model/department/domain/department"
)

type GetDepartments struct {
}

type GetDepartmentsHandler decorator.QueryHandler[GetDepartments, []*department.Department]

/*type GetDepartmentsReadModel interface {
	GetDepartments(ctx context.Context) ([]*department.Department, error)
}*/

type getDepartmentsHandler struct {
	departmentRepo department.Repository
}

func NewGetDepartmentsHandler(departmentRepo department.Repository) GetDepartmentsHandler {
	return decorator.ApplyQueryDecorators[GetDepartments, []*department.Department](
		getDepartmentsHandler{departmentRepo: departmentRepo},
	)
}

func (h getDepartmentsHandler) Handle(ctx context.Context, _ GetDepartments) ([]*department.Department, error) {
	return h.departmentRepo.GetDepartments(ctx)
}
