package client

const (
	RolePath                 = "roles"
	RuntimeGroupsDisplayName = "Runtime Groups"
	ServicesDisplayName      = "Services"
)

type Role struct {
	DisplayName string `json:"name"`
	Description string `json:"description"`
}
type RoleGroup struct {
	DisplayName string          `json:"name"`
	RoleMap     map[string]Role `json:"roles"`
}
type RoleCollection map[string]RoleGroup
