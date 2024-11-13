package domain

import (
	"time"

	"github.com/go-ms-project-store/internal/adapters/input/http/dto"
	"github.com/go-ms-project-store/internal/pkg/helpers"
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

func (u Users) ToDTO() []dto.UserResponse {
	dtos := make([]dto.UserResponse, len(u))
	for i, user := range u {
		dtos[i] = user.ToUserDTO()
	}
	return dtos
}

func (u User) ToMeDTO() dto.UserMeResponse {
	return dto.UserMeResponse{
		ID:        u.UUID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: helpers.DatetimeToString(u.CreatedAt),
	}
}
