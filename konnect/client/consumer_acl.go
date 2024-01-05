package client

import "strings"

const (
	ConsumerACLPath    = ControlPlanePathGet + "/core-entities/consumers/%s/acls"
	ConsumerACLPathGet = ConsumerACLPath + "/%s"
)

type ConsumerACL struct {
	ControlPlaneId string `json:"-"`
	ConsumerId     string `json:"-"`
	Id             string `json:"id"`
	Group          string `json:"group"`
}

func (ck *ConsumerACL) ConsumerACLEncodeId() string {
	return ck.ControlPlaneId + IdSeparator + ck.ConsumerId + IdSeparator + ck.Id
}

func ConsumerACLDecodeId(s string) (string, string, string) {
	tokens := strings.Split(s, IdSeparator)
	return tokens[0], tokens[1], tokens[2]
}
