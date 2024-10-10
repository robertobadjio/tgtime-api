package command

import (
	"context"
	"github.com/robertobadjio/tgtime-api/internal/common/decorator"
	"github.com/robertobadjio/tgtime-api/internal/model/department/domain/department"
)

type UpdateDepartment struct {
	Department *department.Department
}

type UpdateDepartmentHandler decorator.CommandHandler[UpdateDepartment]

type updateDepartmentHandler struct {
	departmentRepository department.Repository
}

func NewUpdateDepartmentHandler(departmentRepository department.Repository) UpdateDepartmentHandler {
	if departmentRepository == nil {
		panic("nil departmentRepository")
	}

	return decorator.ApplyCommandDecorators[UpdateDepartment](
		updateDepartmentHandler{departmentRepository: departmentRepository},
	)
}

func (h updateDepartmentHandler) Handle(ctx context.Context, cmd UpdateDepartment) error {
	_, err := h.departmentRepository.UpdateDepartment(ctx, cmd.Department) // TODO: !
	if err != nil {
		return err
	}

	return nil
}
