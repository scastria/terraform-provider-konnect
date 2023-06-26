package konnect

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-http-utils/headers"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/scastria/terraform-provider-konnect/konnect/client"
	"net/http"
)

func resourcePlugin() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePluginCreate,
		ReadContext:   resourcePluginRead,
		UpdateContext: resourcePluginUpdate,
		DeleteContext: resourcePluginDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"runtime_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "-",
			},
			"protocols": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"grpc", "grpcs", "http", "https", "tcp", "tls", "tls_passthrough", "udp", "ws", "wss"}, false),
				},
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"config_json": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsJSON,
				StateFunc: func(v interface{}) string {
					jsonStr, _ := structure.NormalizeJsonString(v)
					return jsonStr
				},
			},
			"service_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"route_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"consumer_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"config_all_json": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"plugin_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func fillPlugin(c *client.Plugin, d *schema.ResourceData) {
	c.RuntimeGroupId = d.Get("runtime_group_id").(string)
	c.Enabled = d.Get("enabled").(bool)
	name, ok := d.GetOk("name")
	if ok {
		c.Name = name.(string)
	}
	instanceName, ok := d.GetOk("instance_name")
	if ok {
		c.InstanceName = instanceName.(string)
	}
	protocols, ok := d.GetOk("protocols")
	if ok {
		c.Protocols = convertSetToArray(protocols.(*schema.Set))
	}
	service, ok := d.GetOk("service_id")
	if ok {
		c.Service = &client.EntityId{}
		c.Service.Id = service.(string)
	}
	route, ok := d.GetOk("route_id")
	if ok {
		c.Route = &client.EntityId{}
		c.Route.Id = route.(string)
	}
	consumer, ok := d.GetOk("consumer_id")
	if ok {
		c.Consumer = &client.EntityId{}
		c.Consumer.Id = consumer.(string)
	}
	configJson, ok := d.GetOk("config_json")
	if ok {
		c.Config = map[string]interface{}{}
		json.Unmarshal([]byte(configJson.(string)), &c.Config)
	}
	configAllJson, ok := d.GetOk("config_all_json")
	if ok {
		c.ConfigAll = map[string]interface{}{}
		json.Unmarshal([]byte(configAllJson.(string)), &c.ConfigAll)
	}
}

func fillResourceDataFromPlugin(ctx context.Context, c *client.Plugin, d *schema.ResourceData, pluginSchema *client.PluginSchema) {
	// Grab just the schema for config parameter from plugin schema
	var configSchema []map[string]client.PluginField
	configSchema = nil
	for _, value := range pluginSchema.Fields {
		valueMap, ok := value["config"]
		if ok {
			configSchema = valueMap.Fields
			break
		}
	}
	d.Set("runtime_group_id", c.RuntimeGroupId)
	d.Set("name", c.Name)
	d.Set("instance_name", c.InstanceName)
	d.Set("protocols", c.Protocols)
	d.Set("enabled", c.Enabled)
	serviceId := ""
	if c.Service != nil {
		serviceId = c.Service.Id
	}
	d.Set("service_id", serviceId)
	routeId := ""
	if c.Route != nil {
		routeId = c.Route.Id
	}
	d.Set("route_id", routeId)
	consumerId := ""
	if c.Consumer != nil {
		consumerId = c.Consumer.Id
	}
	d.Set("consumer_id", consumerId)
	d.Set("plugin_id", c.Id)
	// Remove all default values from config based on plugin schema
	if configSchema != nil {
		pruneConfigValues(ctx, c.Config, configSchema)
	}
	bytes, _ := json.Marshal(c.Config)
	d.Set("config_json", string(bytes[:]))
	bytes, _ = json.Marshal(c.ConfigAll)
	d.Set("config_all_json", string(bytes[:]))
}

func pruneConfigValues(ctx context.Context, config map[string]interface{}, configSchema []map[string]client.PluginField) {
	for _, fieldConfig := range configSchema {
		for fieldKey, field := range fieldConfig {
			configValue, ok := config[fieldKey]
			if ok {
				// If field type is record, then recurse
				if field.Type == "record" {
					childConfigValue := configValue.(map[string]interface{})
					pruneConfigValues(ctx, childConfigValue, field.Fields)
					// If all fields of child config are defaults, then remove entire child config from parent
					if len(childConfigValue) == 0 {
						delete(config, fieldKey)
					}
				} else if areConfigValuesEqual(ctx, fieldKey, configValue, field.Default) {
					delete(config, fieldKey)
				}
			}
		}
	}
}

func areConfigValuesEqual(ctx context.Context, key string, configValue interface{}, schemaDefault interface{}) bool {
	configValueJSON, _ := json.Marshal(configValue)
	configValueString := string(configValueJSON[:])
	schemaFieldDefaultJSON, _ := json.Marshal(schemaDefault)
	schemaFieldString := string(schemaFieldDefaultJSON[:])
	tflog.Info(ctx, "Comparing config:", map[string]any{
		"key":         key,
		"configValue": configValueString,
		"schemaValue": schemaFieldString,
	})
	return configValueString == schemaFieldString
}

func getPluginSchema(ctx context.Context, c *client.Client, runtimeGroupId string, pluginName string) (*client.PluginSchema, error) {
	requestPath := fmt.Sprintf(client.PluginSchemaPath, runtimeGroupId, pluginName)
	body, err := c.HttpRequest(ctx, true, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return nil, err
	}
	retVal := &client.PluginSchema{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		return nil, err
	}
	return retVal, nil
}

func resourcePluginCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	newPlugin := client.Plugin{}
	fillPlugin(&newPlugin, d)
	err := json.NewEncoder(&buf).Encode(newPlugin)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.PluginPath, newPlugin.RuntimeGroupId)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, true, http.MethodPost, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal := &client.Plugin{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	pluginSchema, err := getPluginSchema(ctx, c, newPlugin.RuntimeGroupId, retVal.Name)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal.ConfigAll = copyMapByJSON(retVal.Config)
	retVal.RuntimeGroupId = newPlugin.RuntimeGroupId
	d.SetId(retVal.PluginEncodeId())
	fillResourceDataFromPlugin(ctx, retVal, d, pluginSchema)
	return diags
}

func resourcePluginRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	runtimeGroupId, id := client.PluginDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.PluginPathGet, runtimeGroupId, id)
	body, err := c.HttpRequest(ctx, true, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.Plugin{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	pluginSchema, err := getPluginSchema(ctx, c, runtimeGroupId, retVal.Name)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal.ConfigAll = copyMapByJSON(retVal.Config)
	retVal.RuntimeGroupId = runtimeGroupId
	fillResourceDataFromPlugin(ctx, retVal, d, pluginSchema)
	return diags
}

func resourcePluginUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	runtimeGroupId, id := client.PluginDecodeId(d.Id())
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	upPlugin := client.Plugin{}
	fillPlugin(&upPlugin, d)
	err := json.NewEncoder(&buf).Encode(upPlugin)
	if err != nil {
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.PluginPathGet, runtimeGroupId, id)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, true, http.MethodPut, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		return diag.FromErr(err)
	}
	retVal := &client.Plugin{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		return diag.FromErr(err)
	}
	pluginSchema, err := getPluginSchema(ctx, c, runtimeGroupId, retVal.Name)
	if err != nil {
		return diag.FromErr(err)
	}
	retVal.ConfigAll = copyMapByJSON(retVal.Config)
	retVal.RuntimeGroupId = runtimeGroupId
	fillResourceDataFromPlugin(ctx, retVal, d, pluginSchema)
	return diags
}

func resourcePluginDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	runtimeGroupId, id := client.PluginDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.PluginPathGet, runtimeGroupId, id)
	_, err := c.HttpRequest(ctx, true, http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
