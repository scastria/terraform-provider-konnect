package client

import "strings"

const (
	TeamUserPath       = "v3/teams/%s/users"
	TeamUserPathCreate = TeamUserPath
	TeamUserPathDelete = TeamUserPath + "/%s"
)

type TeamUser struct {
	TeamId string `json:"-"`
	UserId string `json:"id"`
}

func (tu *TeamUser) TeamUserEncodeId() string {
	return tu.TeamId + IdSeparator + tu.UserId
}

func TeamUserDecodeId(s string) (string, string) {
	tokens := strings.Split(s, IdSeparator)
	return tokens[0], tokens[1]
}
