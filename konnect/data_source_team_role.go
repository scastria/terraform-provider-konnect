package konnect

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/scastria/terraform-provider-konnect/konnect/client"
	"net/http"
	"net/url"
)

func dataSourceTeamRole() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTeamRoleRead,
		Schema: map[string]*schema.Schema{
			"team_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"search_role_display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"role_display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"search_entity_type_display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"entity_type_display_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{client.ControlPlanesDisplayName, client.ServicesDisplayName}, false),
			},
			"entity_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"entity_region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceTeamRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	teamId := d.Get("team_id").(string)
	requestQuery := url.Values{}
	searchRoleDisplayName, ok := d.GetOk("search_role_display_name")
	if ok {
		requestQuery[client.FilterRoleNameContains] = []string{searchRoleDisplayName.(string)}
	}
	roleDisplayName, ok := d.GetOk("role_display_name")
	if ok {
		requestQuery[client.FilterRoleName] = []string{roleDisplayName.(string)}
	}
	searchEntityTypeDisplayName, ok := d.GetOk("search_entity_type_display_name")
	if ok {
		requestQuery[client.FilterEntityTypeNameContains] = []string{searchEntityTypeDisplayName.(string)}
	}
	entityTypeDisplayName, ok := d.GetOk("entity_type_display_name")
	if ok {
		requestQuery[client.FilterEntityTypeName] = []string{entityTypeDisplayName.(string)}
	}
	requestPath := fmt.Sprintf(client.TeamRolePath, teamId)
	body, err := c.HttpRequest(ctx, false, http.MethodGet, requestPath, requestQuery, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVals := &client.TeamRoleCollection{}
	err = json.NewDecoder(body).Decode(retVals)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	numTeamRoles := len(retVals.TeamRoles)
	if numTeamRoles > 1 {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("Filter criteria does not result in a single team role"))
	} else if numTeamRoles != 1 {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("No team role exists with that filter criteria"))
	}
	retVal := retVals.TeamRoles[0]
	retVal.TeamId = teamId
	d.Set("role_display_name", retVal.RoleDisplayName)
	d.Set("entity_type_display_name", retVal.EntityTypeDisplayName)
	d.Set("entity_id", retVal.EntityId)
	d.Set("entity_region", retVal.EntityRegion)
	d.SetId(retVal.TeamRoleEncodeId())
	return diags
}
