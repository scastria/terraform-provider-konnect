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

func resourceControlPlane() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceControlPlaneCreate,
		ReadContext:   resourceControlPlaneRead,
		UpdateContext: resourceControlPlaneUpdate,
		DeleteContext: resourceControlPlaneDelete,
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

func fillControlPlane(c *client.ControlPlane, d *schema.ResourceData) {
	c.Name = d.Get("name").(string)
	description, ok := d.GetOk("description")
	if ok {
		c.Description = description.(string)
	}
}

func fillResourceDataFromControlPlane(c *client.ControlPlane, d *schema.ResourceData) {
	d.Set("name", c.Name)
	d.Set("description", c.Description)
	d.Set("cluster_type", c.Config.ClusterType)
	d.Set("control_plane_endpoint", c.Config.ControlPlaneEndpoint)
	d.Set("telemetry_endpoint", c.Config.TelemetryEndpoint)
}

func resourceControlPlaneCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	newControlPlane := client.ControlPlane{}
	fillControlPlane(&newControlPlane, d)
	err := json.NewEncoder(&buf).Encode(newControlPlane)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.ControlPlanePath)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, true, http.MethodPost, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal := &client.ControlPlane{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId(retVal.Id)
	fillResourceDataFromControlPlane(retVal, d)
	return diags
}

func resourceControlPlaneRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.ControlPlanePathGet, d.Id())
	body, err := c.HttpRequest(ctx, true, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.ControlPlane{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	fillResourceDataFromControlPlane(retVal, d)
	return diags
}

func resourceControlPlaneUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	upControlPlane := client.ControlPlane{}
	fillControlPlane(&upControlPlane, d)
	err := json.NewEncoder(&buf).Encode(upControlPlane)
	if err != nil {
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.ControlPlanePathGet, d.Id())
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, true, http.MethodPatch, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		return diag.FromErr(err)
	}
	retVal := &client.ControlPlane{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		return diag.FromErr(err)
	}
	fillResourceDataFromControlPlane(retVal, d)
	return diags
}

func resourceControlPlaneDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.ControlPlanePathGet, d.Id())
	_, err := c.HttpRequest(ctx, true, http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
