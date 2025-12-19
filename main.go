package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/scastria/terraform-provider-konnect/konnect"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: konnect.Provider,
	})
	//c, _ := client.NewClient("PAT_GOES_HERE", "us", 3, 30, []string{})
	//requestPath := fmt.Sprintf(client.PluginPathGet, "CONTROL_PLANE_ID_GOES_HERE", "PLUGIN_ID_GOES_HERE")
	//body, _ := c.HttpRequestDebug(true, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	//retVal := &client.Plugin{}
	//json.NewDecoder(body).Decode(retVal)
	//pluginSchema, _ := getPluginSchema(c, "CONTROL_PLANE_ID_GOES_HERE", "rate-limiting")
	//var fieldSchema []map[string]client.PluginField
	//fieldSchema = nil
	//var shorthandFieldSchema []map[string]client.PluginField
	//shorthandFieldSchema = nil
	//for _, value := range pluginSchema.Fields {
	//	valueMap, ok := value["config"]
	//	if ok {
	//		fieldSchema = valueMap.Fields
	//		shorthandFieldSchema = valueMap.ShorthandFields
	//		break
	//	}
	//}
	//if fieldSchema != nil {
	//	pruneConfigValues(retVal.Config, fieldSchema, false)
	//}
	//if shorthandFieldSchema != nil {
	//	pruneConfigValues(retVal.Config, shorthandFieldSchema, true)
	//}
}

//func pruneConfigValues(config map[string]interface{}, configSchema []map[string]client.PluginField, isShorthand bool) {
//	for _, fieldConfig := range configSchema {
//		for fieldKey, field := range fieldConfig {
//			configValue, ok := config[fieldKey]
//			if ok {
//				// If field type is record, then recurse
//				if field.Type == "record" {
//					if configValue == nil {
//						delete(config, fieldKey)
//						continue
//					}
//					childConfigValue := configValue.(map[string]interface{})
//					pruneConfigValues(childConfigValue, field.Fields, false)
//					pruneConfigValues(childConfigValue, field.ShorthandFields, true)
//					// If all fields of child config are defaults, then remove entire child config from parent
//					if len(childConfigValue) == 0 {
//						delete(config, fieldKey)
//					}
//				}
//			}
//		}
//	}
//}
//
//func getPluginSchema(c *client.Client, controlPlaneId string, pluginName string) (*client.PluginSchema, error) {
//	requestPath := fmt.Sprintf(client.PluginSchemaPath, controlPlaneId, pluginName)
//	body, err := c.HttpRequestDebug(true, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
//	if err != nil {
//		return nil, err
//	}
//	retVal := &client.PluginSchema{}
//	err = json.NewDecoder(body).Decode(retVal)
//	if err != nil {
//		return nil, err
//	}
//	return retVal, nil
//}
