package client

import "strings"

const (
	ConsumerHMACPath    = ControlPlanePathGet + "/core-entities/consumers/%s/hmac-auth"
	ConsumerHMACPathGet = ConsumerHMACPath + "/%s"
)

type ConsumerHMAC struct {
	ControlPlaneId string `json:"-"`
	ConsumerId     string `json:"-"`
	Id             string `json:"id"`
	Username       string `json:"username"`
	Secret         string `json:"secret,omitempty"`
}

func (ch *ConsumerHMAC) ConsumerHMACEncodeId() string {
	return ch.ControlPlaneId + IdSeparator + ch.ConsumerId + IdSeparator + ch.Id
}

func ConsumerHMACDecodeId(s string) (string, string, string) {
	tokens := strings.Split(s, IdSeparator)
	return tokens[0], tokens[1], tokens[2]
}

//TAGS
