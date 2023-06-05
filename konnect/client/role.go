package client

const (
	RolePath                 = "roles"
	RuntimeGroupsName        = "runtime_groups"
	RuntimeGroupsDisplayName = "Runtime Groups"
	ServicesName             = "services"
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
type RoleCollection struct {
	RuntimeGroups RoleGroup `json:"runtime_groups"`
	Services      RoleGroup `json:"services"`
}
