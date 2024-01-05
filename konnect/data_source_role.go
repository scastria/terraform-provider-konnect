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
)

func dataSourceRole() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRoleRead,
		Schema: map[string]*schema.Schema{
			"entity_type_display_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{client.ControlPlanesDisplayName, client.ServicesDisplayName}, false),
			},
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"entity_type_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	entityTypeDisplayName := d.Get("entity_type_display_name").(string)
	displayName := d.Get("display_name").(string)
	body, err := c.HttpRequest(ctx, false, http.MethodGet, client.RolePath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVals := &client.RoleCollection{}
	err = json.NewDecoder(body).Decode(retVals)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	// Find requested role group by display name
	var foundGroup *client.RoleGroup
	foundGroup = nil
	foundGroupName := ""
	for k, g := range *retVals {
		if g.DisplayName == entityTypeDisplayName {
			foundGroup = &g
			foundGroupName = k
			break
		}
	}
	if foundGroup == nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("No role exists with that filter criteria"))
	}
	// Find requested role by display name
	var foundRole *client.Role
	foundRole = nil
	foundRoleName := ""
	for k, r := range (*retVals)[foundGroupName].RoleMap {
		if r.DisplayName == displayName {
			foundRole = &r
			foundRoleName = k
			break
		}
	}
	if foundRole == nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("No role exists with that filter criteria"))
	}
	d.Set("entity_type_name", foundGroupName)
	d.Set("name", foundRoleName)
	d.Set("description", foundRole.Description)
	d.SetId(foundGroupName + client.IdSeparator + foundRoleName)
	return diags
}
