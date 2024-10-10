package command

import (
	"context"
	"github.com/robertobadjio/tgtime-api/internal/common/decorator"
	"github.com/robertobadjio/tgtime-api/internal/model/department/domain/department"
)

type DeleteDepartment struct {
	DepartmentId int
}

type DeleteDepartmentHandler decorator.CommandHandler[DeleteDepartment]

type deleteDepartmentHandler struct {
	departmentRepository department.Repository
}

func NewDeleteDepartmentHandler(departmentRepository department.Repository) DeleteDepartmentHandler {
	if departmentRepository == nil {
		panic("nil departmentRepository")
	}

	return decorator.ApplyCommandDecorators[DeleteDepartment](
		deleteDepartmentHandler{departmentRepository: departmentRepository},
	)
}

func (h deleteDepartmentHandler) Handle(ctx context.Context, cmd DeleteDepartment) error {
	err := h.departmentRepository.DeleteDepartment(ctx, cmd.DepartmentId)
	if err != nil {
		return err
	}

	return nil
}
