package konnect

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scastria/terraform-provider-konnect/konnect/client"
	"net/http"
	"net/url"
)

func dataSourceRuntimeGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRuntimeGroupRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
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

func dataSourceRuntimeGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	name := d.Get("name").(string)
	c := m.(*client.Client)
	requestQuery := url.Values{}
	requestQuery[client.FilterName] = []string{name}
	body, err := c.HttpRequest(ctx, true, http.MethodGet, client.RuntimeGroupPath, requestQuery, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVals := &client.RuntimeGroupCollection{}
	err = json.NewDecoder(body).Decode(retVals)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	numRuntimeGroups := len(retVals.RuntimeGroups)
	if numRuntimeGroups != 1 {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("No runtime group exists with name: %s", name))
	}
	retVal := retVals.RuntimeGroups[0]
	d.Set("description", retVal.Description)
	d.Set("cluster_type", retVal.Config.ClusterType)
	d.Set("control_plane_endpoint", retVal.Config.ControlPlaneEndpoint)
	d.Set("telemetry_endpoint", retVal.Config.TelemetryEndpoint)
	d.SetId(retVal.Id)
	return diags
}
