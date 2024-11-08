package dto

type RoleResponse struct {
	Id          int64                `json:"id"`
	Name        string               `json:"name"`
	CreatedAt   string               `json:"created_at,omitempty"`
	UpdatedAt   string               `json:"updated_at,omitempty"`
	Permissions []PermissionResponse `json:"permissions,omitempty"`
}
