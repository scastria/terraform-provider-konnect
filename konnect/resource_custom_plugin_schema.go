package konnect

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-http-utils/headers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scastria/terraform-provider-konnect/konnect/client"
	"net/http"
)

func resourceCustomPluginSchema() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCustomPluginSchemaCreate,
		ReadContext:   resourceCustomPluginSchemaRead,
		UpdateContext: resourceCustomPluginSchemaUpdate,
		DeleteContext: resourceCustomPluginSchemaDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"control_plane_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"schema_lua": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func fillCustomPluginSchema(c *client.CustomPluginSchema, d *schema.ResourceData) {
	c.ControlPlaneId = d.Get("control_plane_id").(string)
	c.Name = d.Get("name").(string)
	c.SchemaLua = d.Get("schema_lua").(string)
}

func fillResourceDataFromCustomPluginSchema(c *client.CustomPluginSchema, d *schema.ResourceData) {
	d.Set("control_plane_id", c.ControlPlaneId)
	d.Set("name", c.Name)
	d.Set("schema_lua", c.SchemaLua)
}

func resourceCustomPluginSchemaCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	newCustomPluginSchema := client.CustomPluginSchema{}
	fillCustomPluginSchema(&newCustomPluginSchema, d)
	err := json.NewEncoder(&buf).Encode(newCustomPluginSchema)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.CustomPluginSchemaPath, newCustomPluginSchema.ControlPlaneId)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, true, http.MethodPost, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal := &client.CustomPluginSchemaItem{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal.Item.ControlPlaneId = newCustomPluginSchema.ControlPlaneId
	d.SetId(retVal.Item.CustomPluginSchemaEncodeId())
	fillResourceDataFromCustomPluginSchema(&(retVal.Item), d)
	return diags
}

func resourceCustomPluginSchemaRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	controlPlaneId, name := client.CustomPluginSchemaDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.CustomPluginSchemaPathGet, controlPlaneId, name)
	body, err := c.HttpRequest(ctx, true, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.CustomPluginSchemaItem{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal.Item.ControlPlaneId = controlPlaneId
	fillResourceDataFromCustomPluginSchema(&(retVal.Item), d)
	return diags
}

func resourceCustomPluginSchemaUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	controlPlaneId, name := client.CustomPluginSchemaDecodeId(d.Id())
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	upCustomPluginSchema := client.CustomPluginSchema{}
	fillCustomPluginSchema(&upCustomPluginSchema, d)
	err := json.NewEncoder(&buf).Encode(upCustomPluginSchema)
	if err != nil {
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.CustomPluginSchemaPathGet, controlPlaneId, name)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, true, http.MethodPut, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		return diag.FromErr(err)
	}
	retVal := &client.CustomPluginSchemaItem{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		return diag.FromErr(err)
	}
	retVal.Item.ControlPlaneId = controlPlaneId
	fillResourceDataFromCustomPluginSchema(&(retVal.Item), d)
	return diags
}

func resourceCustomPluginSchemaDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	controlPlaneId, name := client.CustomPluginSchemaDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.CustomPluginSchemaPathGet, controlPlaneId, name)
	_, err := c.HttpRequest(ctx, true, http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
