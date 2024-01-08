package client

import "strings"

const (
	UserRolePath       = "v3/users/%s/assigned-roles"
	UserRolePathCreate = UserRolePath
	UserRolePathDelete = UserRolePath + "/%s"
)

type UserRole struct {
	Id                    string `json:"id,omitempty"`
	UserId                string `json:"-"`
	RoleDisplayName       string `json:"role_name,omitempty"`
	EntityId              string `json:"entity_id,omitempty"`
	EntityTypeDisplayName string `json:"entity_type_name,omitempty"`
	EntityRegion          string `json:"entity_region,omitempty"`
}
type UserRoleCollection struct {
	UserRoles []UserRole `json:"data"`
}

func (ur *UserRole) UserRoleEncodeId() string {
	return ur.UserId + IdSeparator + ur.Id
}

func UserRoleDecodeId(s string) (string, string) {
	tokens := strings.Split(s, IdSeparator)
	return tokens[0], tokens[1]
}
