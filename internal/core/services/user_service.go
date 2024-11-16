package services

import (
	"net/http"

	"github.com/go-ms-project-store/internal/core/domain"
	"github.com/go-ms-project-store/internal/core/ports"
	"github.com/go-ms-project-store/internal/pkg/errs"
	"github.com/go-ms-project-store/internal/pkg/logger"
	"github.com/go-ms-project-store/internal/pkg/pagination"
)

type DefaultUserService struct {
	repo ports.UserRepository
}

func (s DefaultUserService) GetAllUsers(r *http.Request) (domain.Users, int64, pagination.DataDBFilter, *errs.AppError) {
	allowedOrderBy := map[string]bool{
		"id": true, "name": true, "slug": true, "role_id": true, "created_at": true, "updated_at": true,
	}

	filter := pagination.GetBaseFilterParams(r, allowedOrderBy)
	users, totalRows, err := s.repo.FindAll(filter, "")

	if err != nil {
		logger.Error("Error while finding all users")
		return nil, 0, pagination.DataDBFilter{}, errs.NewUnexpectedError("unexpected database error")
	}

	return users, totalRows, filter, nil
}

func (s DefaultUserService) GetAllUserCustomers(r *http.Request) (domain.Users, int64, pagination.DataDBFilter, *errs.AppError) {
	allowedOrderBy := map[string]bool{
		"id": true, "name": true, "slug": true, "role_id": true, "created_at": true, "updated_at": true,
	}

	filter := pagination.GetBaseFilterParams(r, allowedOrderBy)
	users, totalRows, err := s.repo.FindAllCustomers(filter)

	if err != nil {
		logger.Error("Error while finding all users")
		return nil, 0, pagination.DataDBFilter{}, errs.NewUnexpectedError("unexpected database error")
	}

	return users, totalRows, filter, nil
}

func (s DefaultUserService) GetAllUserAdmins(r *http.Request) (domain.Users, int64, pagination.DataDBFilter, *errs.AppError) {
	allowedOrderBy := map[string]bool{
		"id": true, "name": true, "slug": true, "role_id": true, "created_at": true, "updated_at": true,
	}

	filter := pagination.GetBaseFilterParams(r, allowedOrderBy)
	users, totalRows, err := s.repo.FindAllAdmins(filter)

	if err != nil {
		logger.Error("Error while finding all users")
		return nil, 0, pagination.DataDBFilter{}, errs.NewUnexpectedError("unexpected database error")
	}

	return users, totalRows, filter, nil
}

func (s DefaultUserService) FindUserById(id string) (*domain.User, *errs.AppError) {
	user, err := s.repo.FindByUuid(id)
	if err != nil {
		if err.Code != http.StatusNotFound {
			return nil, errs.NewUnexpectedError("unexpected database error")
		} else {
			return nil, err
		}
	}

	return user, nil
}

func (s DefaultUserService) DeleteUser(id string) (bool, *errs.AppError) {
	err := s.repo.Delete(id)
	if err != nil {
		if err.Code != http.StatusNotFound {
			return false, errs.NewUnexpectedError("unexpected database error")
		} else {
			return false, err
		}
	}

	return true, nil
}

func NewUserService(repository ports.UserRepository) DefaultUserService {
	return DefaultUserService{repo: repository}
}
