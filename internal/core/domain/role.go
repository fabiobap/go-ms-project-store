package domain

import (
	"time"

	"github.com/go-ms-project-store/internal/adapters/input/http/dto"
	"github.com/go-ms-project-store/internal/pkg/helpers"
)

type Role struct {
	Id          int64     `db:"id"`
	Name        string    `db:"name"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	Permissions []Permission
}

type Roles []Role

type RoleRepository interface {
}

func (u Role) ToRoleDTO() dto.RoleResponse {
	permissions := make([]dto.PermissionResponse, len(u.Permissions))
	for i, permission := range u.Permissions {
		permissions[i] = permission.ToPermissionDTO()
	}

	return dto.RoleResponse{
		Id:          u.Id,
		Name:        u.Name,
		CreatedAt:   helpers.DatetimeToString(u.CreatedAt),
		UpdatedAt:   helpers.DatetimeToString(u.UpdatedAt),
		Permissions: permissions,
	}
}

func (r Roles) ToDTO() []dto.RoleResponse {
	dtos := make([]dto.RoleResponse, len(r))
	for i, role := range r {
		dtos[i] = role.ToRoleDTO()
	}
	return dtos
}
