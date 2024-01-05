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
)

func dataSourceNodes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNodesRead,
		Schema: map[string]*schema.Schema{
			"control_plane_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"nodes": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hostname": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_ping": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"config_hash": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data_plane_cert_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNodesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	controlPlaneId := d.Get("control_plane_id").(string)
	requestPath := fmt.Sprintf(client.NodePath, controlPlaneId)
	body, err := c.HttpRequest(ctx, true, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVals := &client.NodeCollection{}
	err = json.NewDecoder(body).Decode(retVals)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	var nodes []map[string]interface{}
	for _, item := range retVals.Nodes {
		itemMap := map[string]interface{}{}
		itemMap["id"] = item.Id
		itemMap["version"] = item.Version
		itemMap["hostname"] = item.Hostname
		itemMap["last_ping"] = item.LastPing
		itemMap["type"] = item.Type
		itemMap["config_hash"] = item.ConfigHash
		itemMap["data_plane_cert_id"] = item.DataPlaneCertId
		nodes = append(nodes, itemMap)
	}
	d.Set("nodes", nodes)
	d.SetId(controlPlaneId)
	return diags
}
