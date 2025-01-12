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

func dataSourceService() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceServiceRead,
		Schema: map[string]*schema.Schema{
			"control_plane_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"search_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"search_name"},
			},
			"host": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"retries": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"protocol": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"connect_timeout": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"read_timeout": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"write_timeout": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeSet,
				Computed: true,
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

func dataSourceServiceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	controlPlaneId := d.Get("control_plane_id").(string)
	requestPath := fmt.Sprintf(client.ServicePath, controlPlaneId)
	requestQuery := url.Values{}
	searchName, ok := d.GetOk("search_name")
	if ok {
		requestQuery[client.FilterNameContains] = []string{searchName.(string)}
	}
	name, ok := d.GetOk("name")
	if ok {
		requestQuery[client.FilterNameEquals] = []string{name.(string)}
	}
	body, err := c.HttpRequest(ctx, false, http.MethodGet, requestPath, requestQuery, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVals := &client.ListServiceCollection{}
	err = json.NewDecoder(body).Decode(retVals)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	numServices := len(retVals.Services)
	if numServices > 1 {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("Filter criteria does not result in a single service"))
	} else if numServices != 1 {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("No service exists with that filter criteria"))
	}
	retVal := retVals.Services[0]
	retVal.ControlPlaneId = controlPlaneId
	d.Set("control_plane_id", retVal.ControlPlaneId)
	d.Set("host", retVal.Host)
	d.Set("name", retVal.Name)
	d.Set("retries", retVal.Retries)
	d.Set("protocol", retVal.Protocol)
	d.Set("port", retVal.Port)
	d.Set("path", retVal.Path)
	d.Set("connect_timeout", retVal.ConnectTimeout)
	d.Set("read_timeout", retVal.ReadTimeout)
	d.Set("write_timeout", retVal.WriteTimeout)
	d.Set("enabled", retVal.Enabled)
	d.Set("all_tags", retVal.AllTags)
	d.Set("tags", subtractArrays(retVal.AllTags, c.DefaultTags))
	d.Set("service_id", retVal.Id)
	d.SetId(retVal.ServiceEncodeId())
	return diags
}
