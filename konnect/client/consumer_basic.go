package client

import "strings"

const (
	ConsumerBasicPath    = ControlPlanePathGet + "/core-entities/consumers/%s/basic-auth"
	ConsumerBasicPathGet = ConsumerBasicPath + "/%s"
)

type ConsumerBasic struct {
	ControlPlaneId string `json:"-"`
	ConsumerId     string `json:"-"`
	Id             string `json:"id"`
	Username       string `json:"username"`
	Password       string `json:"password"`
}

func (ck *ConsumerBasic) ConsumerBasicEncodeId() string {
	return ck.ControlPlaneId + IdSeparator + ck.ConsumerId + IdSeparator + ck.Id
}

func ConsumerBasicDecodeId(s string) (string, string, string) {
	tokens := strings.Split(s, IdSeparator)
	return tokens[0], tokens[1], tokens[2]
}
