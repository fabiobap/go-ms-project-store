package domain

import (
	"time"

	"github.com/go-ms-project-store/internal/adapters/input/http/dto"
	"github.com/go-ms-project-store/internal/pkg/helpers"
)

type Permission struct {
	Id        int64     `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Roles     []Role
}

type Permissions []Permission

type PermissionRepository interface {
}

func (u Permission) ToPermissionDTO() dto.PermissionResponse {
	return dto.PermissionResponse{
		Id:        u.Id,
		Name:      u.Name,
		CreatedAt: helpers.DatetimeToString(u.CreatedAt),
		UpdatedAt: helpers.DatetimeToString(u.UpdatedAt),
	}
}

func (p Permissions) ToDTO() []dto.PermissionResponse {
	dtos := make([]dto.PermissionResponse, len(p))
	for i, permission := range p {
		dtos[i] = permission.ToPermissionDTO()
	}
	return dtos
}
