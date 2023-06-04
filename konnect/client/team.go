package client

const (
	TeamPath    = "teams"
	TeamPathGet = TeamPath + "/%s"
)

type Team struct {
	Id           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Description  string `json:"description,omitempty"`
	IsPredefined bool   `json:"system_team,omitempty"`
}
type TeamCollection struct {
	Teams []Team `json:"data"`
}
