package client

import "strings"

const (
	CustomPluginSchemaPath    = ControlPlanePathGet + "/core-entities/plugin-schemas"
	CustomPluginSchemaPathGet = CustomPluginSchemaPath + "/%s"
)

type CustomPluginSchema struct {
	ControlPlaneId string `json:"-"`
	Name           string `json:"name"`
	SchemaLua      string `json:"lua_schema"`
}
type CustomPluginSchemaItem struct {
	Item CustomPluginSchema `json:"item"`
}

func (s *CustomPluginSchema) CustomPluginSchemaEncodeId() string {
	return s.ControlPlaneId + IdSeparator + s.Name
}

func CustomPluginSchemaDecodeId(s string) (string, string) {
	tokens := strings.Split(s, IdSeparator)
	return tokens[0], tokens[1]
}
