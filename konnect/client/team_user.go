package client

import "strings"

const (
	TeamUserPath       = "teams/%s/users"
	TeamUserPathCreate = TeamUserPath
	TeamUserPathDelete = TeamUserPath + "/%s"
)

type TeamUser struct {
	TeamId string `json:"-"`
	UserId string `json:"id"`
}

func (tm *TeamUser) TeamUserEncodeId() string {
	return tm.TeamId + IdSeparator + tm.UserId
}

func TeamUserDecodeId(s string) (string, string) {
	tokens := strings.Split(s, IdSeparator)
	return tokens[0], tokens[1]
}
