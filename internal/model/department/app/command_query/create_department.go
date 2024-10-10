package command_query

import (
	"context"
	"github.com/robertobadjio/tgtime-api/internal/common/decorator"
	"github.com/robertobadjio/tgtime-api/internal/model/department/domain/department"
)

type CreateDepartment struct {
	Department *department.Department
}

type CreateDepartmentHandler decorator.CommandQueryHandler[CreateDepartment, *department.Department]

type createDepartmentHandler struct {
	departmentRepository department.Repository
}

func NewCreateDepartmentHandler(departmentRepository department.Repository) CreateDepartmentHandler {
	if departmentRepository == nil {
		panic("nil departmentRepository")
	}

	return decorator.ApplyCommandQueryDecorators[CreateDepartment, *department.Department](
		createDepartmentHandler{departmentRepository: departmentRepository},
	)
}

func (h createDepartmentHandler) Handle(ctx context.Context, cmdQr CreateDepartment) (*department.Department, error) {
	departmentNew, err := h.departmentRepository.CreateDepartment(ctx, cmdQr.Department)
	if err != nil {
		return nil, err
	}

	return departmentNew, nil
}
