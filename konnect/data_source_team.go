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

func dataSourceTeam() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTeamRead,
		Schema: map[string]*schema.Schema{
			"search_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_predefined": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceTeamRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	requestQuery := url.Values{}
	searchName, ok := d.GetOk("search_name")
	if ok {
		requestQuery[client.FilterNameContains] = []string{searchName.(string)}
	}
	name, ok := d.GetOk("name")
	if ok {
		requestQuery[client.FilterName] = []string{name.(string)}
	}
	body, err := c.HttpRequest(ctx, false, http.MethodGet, client.TeamPath, requestQuery, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVals := &client.TeamCollection{}
	err = json.NewDecoder(body).Decode(retVals)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	numTeams := len(retVals.Teams)
	if numTeams > 1 {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("Filter criteria does not result in a single team"))
	} else if numTeams != 1 {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("No team exists with that filter criteria"))
	}
	retVal := retVals.Teams[0]
	d.Set("name", retVal.Name)
	d.Set("description", retVal.Description)
	d.Set("is_predefined", retVal.IsPredefined)
	d.SetId(retVal.Id)
	return diags
}
