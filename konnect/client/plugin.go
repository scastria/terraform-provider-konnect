package client

import "strings"

const (
	PluginPath       = ControlPlanePathGet + "/core-entities/plugins"
	PluginPathGet    = PluginPath + "/%s"
	PluginSchemaPath = ControlPlanePathGet + "/core-entities/schemas/plugins/%s"
)

type Plugin struct {
	ControlPlaneId string                 `json:"-"`
	Id             string                 `json:"id"`
	Name           string                 `json:"name"`
	InstanceName   string                 `json:"instance_name"`
	Protocols      []string               `json:"protocols"`
	Enabled        bool                   `json:"enabled"`
	Config         map[string]interface{} `json:"config"`
	ConfigAll      map[string]interface{} `json:"-"`
	Route          *EntityId              `json:"route,omitempty"`
	Service        *EntityId              `json:"service,omitempty"`
	Consumer       *EntityId              `json:"consumer,omitempty"`
	AllTags        []string               `json:"tags"`
	Tags           []string               `json:"-"`
}
type PluginCollection struct {
	Plugins []Plugin `json:"data"`
}
type PluginSchema struct {
	Fields []map[string]PluginField `json:"fields"`
}
type PluginField struct {
	Type            string                   `json:"type"`
	Default         interface{}              `json:"default"`
	Fields          []map[string]PluginField `json:"fields"`
	ShorthandFields []map[string]PluginField `json:"shorthand_fields"`
}

func (s *Plugin) PluginEncodeId() string {
	return s.ControlPlaneId + IdSeparator + s.Id
}

func PluginDecodeId(s string) (string, string) {
	tokens := strings.Split(s, IdSeparator)
	return tokens[0], tokens[1]
}
