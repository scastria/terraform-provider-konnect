package client

import "strings"

const (
	ConsumerJWTPath    = ControlPlanePathGet + "/core-entities/consumers/%s/jwt"
	ConsumerJWTPathGet = ConsumerJWTPath + "/%s"
)

type ConsumerJWT struct {
	ControlPlaneId string `json:"-"`
	ConsumerId     string `json:"-"`
	Id             string `json:"id"`
	Key            string `json:"key,omitempty"`
	Algorithm      string `json:"algorithm"`
	RSAPublicKey   string `json:"rsa_public_key,omitempty"`
	Secret         string `json:"secret,omitempty"`
}

func (cj *ConsumerJWT) ConsumerJWTEncodeId() string {
	return cj.ControlPlaneId + IdSeparator + cj.ConsumerId + IdSeparator + cj.Id
}

func ConsumerJWTDecodeId(s string) (string, string, string) {
	tokens := strings.Split(s, IdSeparator)
	return tokens[0], tokens[1], tokens[2]
}

//TAGS
