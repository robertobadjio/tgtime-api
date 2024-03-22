package api

import (
	"context"
	"errors"
	"officetime-api/internal/model/department/app/command"
	"officetime-api/internal/model/department/app/command_query"
	"officetime-api/internal/model/department/app/query"
	"officetime-api/internal/model/department/domain/department"
)

func (s *apiService) GetDepartments(ctx context.Context) ([]*department.Department, error) {
	qr := query.GetDepartments{}
	departments, err := s.departmentApp.Queries.GetDepartments.Handle(ctx, qr)
	if err != nil {
		return nil, err
	}

	return departments, nil
}

func (s *apiService) GetDepartment(ctx context.Context, departmentId int) (*department.Department, error) {
	qr := query.GetDepartment{DepartmentId: departmentId}
	departmentUpdated, err := s.departmentApp.Queries.GetDepartment.Handle(ctx, qr)
	if err != nil {
		return nil, err
	}

	return departmentUpdated, nil
}

func (s *apiService) CreateDepartment(ctx context.Context, department *department.Department) (*department.Department, error) {
	cmd := command_query.CreateDepartment{Department: department}
	departmentNew, err := s.departmentApp.CommandsQueries.CreateDepartment.Handle(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return departmentNew, nil
}

func (s *apiService) UpdateDepartment(ctx context.Context, departmentId int, department *department.Department) (*department.Department, error) {
	if departmentId != department.Id {
		return nil, errors.New("error update ids not equals")
	}

	cmd := command.UpdateDepartment{Department: department}
	err := s.departmentApp.Commands.UpdateDepartment.Handle(ctx, cmd)
	if err != nil {
		return nil, err
	}

	qr := query.GetDepartment{DepartmentId: departmentId}
	departmentNew, err := s.departmentApp.Queries.GetDepartment.Handle(ctx, qr)
	if err != nil {
		return nil, err
	}

	return departmentNew, nil
}

func (s *apiService) DeleteDepartment(ctx context.Context, departmentId int) error {
	cmd := command.DeleteDepartment{DepartmentId: departmentId}
	err := s.departmentApp.Commands.DeleteDepartment.Handle(ctx, cmd)
	if err != nil {
		return err
	}

	return nil
}
