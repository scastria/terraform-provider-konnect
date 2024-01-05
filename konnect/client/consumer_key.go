package client

import "strings"

const (
	ConsumerKeyPath    = ControlPlanePathGet + "/core-entities/consumers/%s/key-auth"
	ConsumerKeyPathGet = ConsumerKeyPath + "/%s"
)

type ConsumerKey struct {
	ControlPlaneId string `json:"-"`
	ConsumerId     string `json:"-"`
	Id             string `json:"id"`
	Key            string `json:"key,omitempty"`
}

func (ck *ConsumerKey) ConsumerKeyEncodeId() string {
	return ck.ControlPlaneId + IdSeparator + ck.ConsumerId + IdSeparator + ck.Id
}

func ConsumerKeyDecodeId(s string) (string, string, string) {
	tokens := strings.Split(s, IdSeparator)
	return tokens[0], tokens[1], tokens[2]
}
