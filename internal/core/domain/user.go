package domain

import (
	"time"

	"github.com/go-ms-project-store/internal/adapters/input/http/dto"
	"github.com/go-ms-project-store/internal/pkg/errs"
	"github.com/go-ms-project-store/internal/pkg/helpers"
	"github.com/go-ms-project-store/internal/pkg/pagination"
	"github.com/google/uuid"
)

type User struct {
	Id              int64     `db:"id"`
	UUID            uuid.UUID `db:"uuid"`
	Name            string    `db:"name"`
	Email           string    `db:"email"`
	Password        string    `db:"password"`
	RoleId          int64     `db:"role_id"`
	EmailVerifiedAt time.Time `db:"email_verified_at"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
	Role            Role
}

type Users []User

type UserRepository interface {
	Delete(id string) *errs.AppError
	FindAll(filter pagination.DataDBFilter, roleName string) (Users, int64, *errs.AppError)
	FindAllAdmins(filter pagination.DataDBFilter) (Users, int64, *errs.AppError)
	FindAllCustomers(filter pagination.DataDBFilter) (Users, int64, *errs.AppError)
	FindById(id string) (*User, *errs.AppError)
}

func (u User) ToUserDTO() dto.UserResponse {
	return dto.UserResponse{
		Id:              u.Id,
		UUID:            u.UUID,
		Name:            u.Name,
		Email:           u.Email,
		RoleId:          u.RoleId,
		EmailVerifiedAt: helpers.DatetimeToString(u.EmailVerifiedAt),
		CreatedAt:       helpers.DatetimeToString(u.CreatedAt),
		UpdatedAt:       helpers.DatetimeToString(u.UpdatedAt),
		Role: dto.RoleResponse{
			Id:   u.Role.Id,
			Name: u.Role.Name,
		},
	}
}

func (c Users) ToDTO() []dto.UserResponse {
	dtos := make([]dto.UserResponse, len(c))
	for i, user := range c {
		dtos[i] = user.ToUserDTO()
	}
	return dtos
}
