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
	"math"
	"net/http"
)

func resourceRuntimeGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRuntimeGroupCreate,
		ReadContext:   resourceRuntimeGroupRead,
		UpdateContext: resourceRuntimeGroupUpdate,
		DeleteContext: resourceRuntimeGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(2, math.MaxInt),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"control_plane_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"telemetry_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceRuntimeGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	newRunGroup := client.RuntimeGroup{}
	fillRuntimeGroup(&newRunGroup, d)
	err := json.NewEncoder(&buf).Encode(newRunGroup)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.RuntimeGroupPath)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, true, http.MethodPost, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal := &client.RuntimeGroup{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId(retVal.Id)
	fillResourceData(retVal, d)
	return diags
}

func fillRuntimeGroup(c *client.RuntimeGroup, d *schema.ResourceData) {
	c.Name = d.Get("name").(string)
	description, ok := d.GetOk("description")
	if ok {
		c.Description = description.(string)
	}
}

func fillResourceData(c *client.RuntimeGroup, d *schema.ResourceData) {
	d.Set("name", c.Name)
	d.Set("description", c.Description)
	d.Set("cluster_type", c.Config.ClusterType)
	d.Set("control_plane_endpoint", c.Config.ControlPlaneEndpoint)
	d.Set("telemetry_endpoint", c.Config.TelemetryEndpoint)
}

func resourceRuntimeGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.RuntimeGroupPathGet, d.Id())
	body, err := c.HttpRequest(ctx, true, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.RuntimeGroup{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	fillResourceData(retVal, d)
	return diags
}

func resourceRuntimeGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	upRunGroup := client.RuntimeGroup{}
	fillRuntimeGroup(&upRunGroup, d)
	err := json.NewEncoder(&buf).Encode(upRunGroup)
	if err != nil {
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.RuntimeGroupPathGet, d.Id())
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, true, http.MethodPatch, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		return diag.FromErr(err)
	}
	retVal := &client.RuntimeGroup{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		return diag.FromErr(err)
	}
	fillResourceData(retVal, d)
	return diags
}

func resourceRuntimeGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.RuntimeGroupPathGet, d.Id())
	_, err := c.HttpRequest(ctx, true, http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
