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

func resourceTeamRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTeamRoleCreate,
		ReadContext:   resourceTeamRoleRead,
		DeleteContext: resourceTeamRoleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"team_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role_display_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"entity_type_display_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{client.RuntimeGroupsDisplayName, client.ServicesDisplayName}, false),
			},
			"entity_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"entity_region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func fillTeamRole(c *client.TeamRole, d *schema.ResourceData) {
	c.TeamId = d.Get("team_id").(string)
	c.RoleDisplayName = d.Get("role_display_name").(string)
	c.EntityTypeDisplayName = d.Get("entity_type_display_name").(string)
	c.EntityId = d.Get("entity_id").(string)
	c.EntityRegion = d.Get("entity_region").(string)
}

func fillResourceDataFromTeamRole(c *client.TeamRole, d *schema.ResourceData) {
	d.Set("team_id", c.TeamId)
	d.Set("role_display_name", c.RoleDisplayName)
	d.Set("entity_type_display_name", c.EntityTypeDisplayName)
	d.Set("entity_id", c.EntityId)
	d.Set("entity_region", c.EntityRegion)
}

func resourceTeamRoleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	newTeamRole := client.TeamRole{}
	fillTeamRole(&newTeamRole, d)
	err := json.NewEncoder(&buf).Encode(newTeamRole)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.TeamRolePathCreate, newTeamRole.TeamId)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, false, http.MethodPost, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal := &client.TeamRole{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	// Add fields missing from POST response
	retVal.TeamId = newTeamRole.TeamId
	retVal.RoleDisplayName = newTeamRole.RoleDisplayName
	retVal.EntityTypeDisplayName = newTeamRole.EntityTypeDisplayName
	retVal.EntityId = newTeamRole.EntityId
	retVal.EntityRegion = newTeamRole.EntityRegion
	d.SetId(retVal.TeamRoleEncodeId())
	fillResourceDataFromTeamRole(retVal, d)
	return diags
}

func resourceTeamRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	teamId, id := client.TeamRoleDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.TeamRolePath, teamId)
	body, err := c.HttpRequest(ctx, false, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.TeamRoleCollection{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	//Search assigned roles looking for id
	var foundTeamRole *client.TeamRole
	foundTeamRole = nil
	for _, tr := range retVal.TeamRoles {
		if tr.Id == id {
			foundTeamRole = &tr
			break
		}
	}
	if foundTeamRole == nil {
		d.SetId("")
		return diags
	}
	foundTeamRole.TeamId = teamId
	fillResourceDataFromTeamRole(foundTeamRole, d)
	return diags
}

func resourceTeamRoleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	teamId, id := client.TeamRoleDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.TeamRolePathDelete, teamId, id)
	_, err := c.HttpRequest(ctx, false, http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
