package client

const (
	TeamMappingsId   = "team-mappings"
	TeamMappingsPath = IdentityProviderPath + "/" + TeamMappingsId
)

type TeamMapping struct {
	Group   string   `json:"group,omitempty"`
	TeamIds []string `json:"team_ids,omitempty"`
}
type TeamMappings struct {
	MappingsRead  []TeamMapping `json:"data,omitempty"`
	MappingsWrite []TeamMapping `json:"mappings,omitempty"`
}
