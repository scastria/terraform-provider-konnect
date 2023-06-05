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
			"group_display_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{client.RuntimeGroupsDisplayName, client.ServicesDisplayName}, false),
			},
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_name": {
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
	groupDisplayName := d.Get("group_display_name").(string)
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
	// Find requested role by display names
	var roleGroup client.RoleGroup
	var groupName string
	if groupDisplayName == client.RuntimeGroupsDisplayName {
		roleGroup = retVals.RuntimeGroups
		groupName = client.RuntimeGroupsName
	} else {
		roleGroup = retVals.Services
		groupName = client.ServicesName
	}
	var foundRole *client.Role
	foundRole = nil
	foundRoleName := ""
	for k, r := range roleGroup.RoleMap {
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
	d.Set("group_name", groupName)
	d.Set("name", foundRoleName)
	d.Set("description", foundRole.Description)
	d.SetId(groupName + client.IdSeparator + foundRoleName)
	return diags
}
