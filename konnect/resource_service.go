package konnect

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-http-utils/headers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/scastria/terraform-provider-konnect/konnect/client"
	"net/http"
)

func resourceService() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServiceCreate,
		ReadContext:   resourceServiceRead,
		UpdateContext: resourceServiceUpdate,
		DeleteContext: resourceServiceDelete,
		CustomizeDiff: resourceServiceDiff,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"control_plane_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"host": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"retries": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      5,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "http",
				ValidateFunc: validation.StringInSlice([]string{"grpc", "grpcs", "http", "https", "tcp", "tls", "tls_passthrough", "udp", "ws", "wss"}, false),
			},
			"port": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      80,
				ValidateFunc: validation.IntBetween(0, 65535),
			},
			"path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"connect_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      60000,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"read_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      60000,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"write_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      60000,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"all_tags": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"service_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceServiceDiff(ctx context.Context, diff *schema.ResourceDiff, m interface{}) error {
	c := m.(*client.Client)
	tags := []string{}
	tagsSet, ok := diff.GetOk("tags")
	if ok {
		tags = convertSetToArray(tagsSet.(*schema.Set))
	}
	allTags := unionArrays(tags, c.DefaultTags)
	diff.SetNew("all_tags", allTags)
	return nil
}

func fillService(c *client.Service, d *schema.ResourceData, defaultTags []string) {
	c.ControlPlaneId = d.Get("control_plane_id").(string)
	c.Host = d.Get("host").(string)
	c.Enabled = d.Get("enabled").(bool)
	name, ok := d.GetOk("name")
	if ok {
		c.Name = name.(string)
	}
	retries, ok := d.GetOk("retries")
	if ok {
		c.Retries = retries.(int)
	}
	protocol, ok := d.GetOk("protocol")
	if ok {
		c.Protocol = protocol.(string)
	}
	port, ok := d.GetOk("port")
	if ok {
		c.Port = port.(int)
	}
	path, ok := d.GetOk("path")
	if ok {
		c.Path = path.(string)
	}
	connectTimeout, ok := d.GetOk("connect_timeout")
	if ok {
		c.ConnectTimeout = connectTimeout.(int)
	}
	readTimeout, ok := d.GetOk("read_timeout")
	if ok {
		c.ReadTimeout = readTimeout.(int)
	}
	writeTimeout, ok := d.GetOk("write_timeout")
	if ok {
		c.WriteTimeout = writeTimeout.(int)
	}
	tags := []string{}
	tagsSet, ok := d.GetOk("tags")
	if ok {
		tags = convertSetToArray(tagsSet.(*schema.Set))
		c.Tags = tags
	}
	c.AllTags = unionArrays(tags, defaultTags)
}

func fillResourceDataFromService(c *client.Service, d *schema.ResourceData, defaultTags []string) {
	d.Set("control_plane_id", c.ControlPlaneId)
	d.Set("host", c.Host)
	d.Set("name", c.Name)
	d.Set("retries", c.Retries)
	d.Set("protocol", c.Protocol)
	d.Set("port", c.Port)
	d.Set("path", c.Path)
	d.Set("connect_timeout", c.ConnectTimeout)
	d.Set("read_timeout", c.ReadTimeout)
	d.Set("write_timeout", c.WriteTimeout)
	d.Set("enabled", c.Enabled)
	d.Set("all_tags", c.AllTags)
	d.Set("tags", subtractArrays(c.AllTags, defaultTags))
	d.Set("service_id", c.Id)
}

func resourceServiceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	newService := client.Service{}
	fillService(&newService, d, c.DefaultTags)
	err := json.NewEncoder(&buf).Encode(newService)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.ServicePath, newService.ControlPlaneId)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, true, http.MethodPost, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal := &client.Service{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal.ControlPlaneId = newService.ControlPlaneId
	d.SetId(retVal.ServiceEncodeId())
	fillResourceDataFromService(retVal, d, c.DefaultTags)
	return diags
}

func resourceServiceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	controlPlaneId, id := client.ServiceDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.ServicePathGet, controlPlaneId, id)
	body, err := c.HttpRequest(ctx, true, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.Service{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal.ControlPlaneId = controlPlaneId
	fillResourceDataFromService(retVal, d, c.DefaultTags)
	return diags
}

func resourceServiceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	controlPlaneId, id := client.ServiceDecodeId(d.Id())
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	upService := client.Service{}
	fillService(&upService, d, c.DefaultTags)
	err := json.NewEncoder(&buf).Encode(upService)
	if err != nil {
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.ServicePathGet, controlPlaneId, id)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, true, http.MethodPut, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		return diag.FromErr(err)
	}
	retVal := &client.Service{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		return diag.FromErr(err)
	}
	retVal.ControlPlaneId = controlPlaneId
	fillResourceDataFromService(retVal, d, c.DefaultTags)
	return diags
}

func resourceServiceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	controlPlaneId, id := client.ServiceDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.ServicePathGet, controlPlaneId, id)
	_, err := c.HttpRequest(ctx, true, http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
