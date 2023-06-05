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

func resourceUserRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserRoleCreate,
		ReadContext:   resourceUserRoleRead,
		DeleteContext: resourceUserRoleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"user_id": {
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

func fillUserRole(c *client.UserRole, d *schema.ResourceData) {
	c.UserId = d.Get("user_id").(string)
	c.RoleDisplayName = d.Get("role_display_name").(string)
	c.EntityTypeDisplayName = d.Get("entity_type_display_name").(string)
	c.EntityId = d.Get("entity_id").(string)
	c.EntityRegion = d.Get("entity_region").(string)
}

func fillResourceDataFromUserRole(c *client.UserRole, d *schema.ResourceData) {
	d.Set("user_id", c.UserId)
	d.Set("role_display_name", c.RoleDisplayName)
	d.Set("entity_type_display_name", c.EntityTypeDisplayName)
	d.Set("entity_id", c.EntityId)
	d.Set("entity_region", c.EntityRegion)
}

func resourceUserRoleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	newUserRole := client.UserRole{}
	fillUserRole(&newUserRole, d)
	err := json.NewEncoder(&buf).Encode(newUserRole)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.UserRolePathCreate, newUserRole.UserId)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, false, http.MethodPost, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal := &client.UserRole{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	// Add fields missing from POST response
	retVal.UserId = newUserRole.UserId
	retVal.RoleDisplayName = newUserRole.RoleDisplayName
	retVal.EntityTypeDisplayName = newUserRole.EntityTypeDisplayName
	retVal.EntityId = newUserRole.EntityId
	retVal.EntityRegion = newUserRole.EntityRegion
	d.SetId(retVal.UserRoleEncodeId())
	fillResourceDataFromUserRole(retVal, d)
	return diags
}

func resourceUserRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	userId, id := client.UserRoleDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.UserRolePath, userId)
	body, err := c.HttpRequest(ctx, false, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.UserRoleCollection{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	//Search assigned roles looking for id
	var foundUserRole *client.UserRole
	foundUserRole = nil
	for _, tr := range retVal.UserRoles {
		if tr.Id == id {
			foundUserRole = &tr
			break
		}
	}
	if foundUserRole == nil {
		d.SetId("")
		return diags
	}
	foundUserRole.UserId = userId
	fillResourceDataFromUserRole(foundUserRole, d)
	return diags
}

func resourceUserRoleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	userId, id := client.UserRoleDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.UserRolePathDelete, userId, id)
	_, err := c.HttpRequest(ctx, false, http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
