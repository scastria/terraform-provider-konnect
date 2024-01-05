package client

import "strings"

const (
	ConsumerPath    = ControlPlanePathGet + "/core-entities/consumers"
	ConsumerPathGet = ConsumerPath + "/%s"
)

type Consumer struct {
	ControlPlaneId string `json:"-"`
	Id             string `json:"id"`
	Username       string `json:"username"`
	CustomId       string `json:"custom_id"`
}
type ConsumerCollection struct {
	Consumers []Consumer `json:"data"`
}

func (c *Consumer) ConsumerEncodeId() string {
	return c.ControlPlaneId + IdSeparator + c.Id
}

func ConsumerDecodeId(s string) (string, string) {
	tokens := strings.Split(s, IdSeparator)
	return tokens[0], tokens[1]
}
