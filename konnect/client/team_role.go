package client

import "strings"

const (
	TeamRolePath       = "teams/%s/assigned-roles"
	TeamRolePathCreate = TeamRolePath
	TeamRolePathDelete = TeamRolePath + "/%s"
)

type TeamRole struct {
	Id                    string `json:"id,omitempty"`
	TeamId                string `json:"-"`
	RoleDisplayName       string `json:"role_name,omitempty"`
	EntityId              string `json:"entity_id,omitempty"`
	EntityTypeDisplayName string `json:"entity_type_name,omitempty"`
	EntityRegion          string `json:"entity_region,omitempty"`
}
type TeamRoleCollection struct {
	TeamRoles []TeamRole `json:"data"`
}

func (tr *TeamRole) TeamRoleEncodeId() string {
	return tr.TeamId + IdSeparator + tr.Id
}

func TeamRoleDecodeId(s string) (string, string) {
	tokens := strings.Split(s, IdSeparator)
	return tokens[0], tokens[1]
}
