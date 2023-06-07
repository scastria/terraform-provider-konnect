package client

import "strings"

const (
	ServicePath    = RuntimeGroupPathGet + "/core-entities/services"
	ServicePathGet = ServicePath + "/%s"
)

type Service struct {
	RuntimeGroupId string `json:"-"`
	Id             string `json:"id"`
	Name           string `json:"name"`
	Retries        int    `json:"retries"`
	Protocol       string `json:"protocol"`
	Host           string `json:"host"`
	Port           int    `json:"port"`
	Path           string `json:"path"`
	ConnectTimeout int    `json:"connect_timeout"`
	ReadTimeout    int    `json:"read_timeout"`
	WriteTimeout   int    `json:"write_timeout"`
	Enabled        bool   `json:"enabled"`
}
type ServiceCollection struct {
	Services []Service `json:"data"`
}

func (s *Service) ServiceEncodeId() string {
	return s.RuntimeGroupId + IdSeparator + s.Id
}

func ServiceDecodeId(s string) (string, string) {
	tokens := strings.Split(s, IdSeparator)
	return tokens[0], tokens[1]
}
