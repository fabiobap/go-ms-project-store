package dto

type PermissionResponse struct {
	Id        int64          `json:"id"`
	Name      string         `json:"name"`
	CreatedAt string         `json:"created_at"`
	UpdatedAt string         `json:"updated_at"`
	Roles     []RoleResponse `json:"roles,omitempty"`
}
