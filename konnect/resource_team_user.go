package konnect

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-http-utils/headers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scastria/terraform-provider-konnect/konnect/client"
	"net/http"
)

func resourceTeamUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTeamUserCreate,
		ReadContext:   resourceTeamUserRead,
		DeleteContext: resourceTeamUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"team_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func fillTeamUser(c *client.TeamUser, d *schema.ResourceData) {
	c.TeamId = d.Get("team_id").(string)
	c.UserId = d.Get("user_id").(string)
}

func fillResourceDataFromTeamUser(c *client.TeamUser, d *schema.ResourceData) {
	d.Set("team_id", c.TeamId)
	d.Set("user_id", c.UserId)
}

func resourceTeamUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	newTeamUser := client.TeamUser{}
	fillTeamUser(&newTeamUser, d)
	err := json.NewEncoder(&buf).Encode(newTeamUser)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.TeamUserPathCreate, newTeamUser.TeamId)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	_, err = c.HttpRequest(ctx, false, http.MethodPost, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId(newTeamUser.TeamUserEncodeId())
	return diags
}

func resourceTeamUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	teamId, userId := client.TeamUserDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.TeamUserPath, teamId)
	body, err := c.HttpRequest(ctx, false, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.UserCollection{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	//Search users looking for userId
	foundUser := false
	for _, u := range retVal.Users {
		if u.Id == userId {
			foundUser = true
			break
		}
	}
	if !foundUser {
		d.SetId("")
		return diags
	}
	d.Set("team_id", teamId)
	d.Set("user_id", userId)
	return diags
}

func resourceTeamUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	teamId, userId := client.TeamUserDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.TeamUserPathDelete, teamId, userId)
	_, err := c.HttpRequest(ctx, false, http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
