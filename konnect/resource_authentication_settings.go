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

func resourceAuthenticationSettings() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAuthenticationSettingsCreate,
		ReadContext:   resourceAuthenticationSettingsRead,
		UpdateContext: resourceAuthenticationSettingsUpdate,
		DeleteContext: resourceAuthenticationSettingsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"basic_auth_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"oidc_auth_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"idp_mapping_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"konnect_mapping_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func fillAuthenticationSettings(c *client.AuthenticationSettings, d *schema.ResourceData) {
	c.BasicAuthEnabled = d.Get("basic_auth_enabled").(bool)
	c.OIDCAuthEnabled = d.Get("oidc_auth_enabled").(bool)
	c.IDPMappingEnabled = d.Get("idp_mapping_enabled").(bool)
	c.KonnectMappingEnabled = d.Get("konnect_mapping_enabled").(bool)
}

func fillResourceDataFromAuthenticationSettings(c *client.AuthenticationSettings, d *schema.ResourceData) {
	d.Set("basic_auth_enabled", c.BasicAuthEnabled)
	d.Set("oidc_auth_enabled", c.OIDCAuthEnabled)
	d.Set("idp_mapping_enabled", c.IDPMappingEnabled)
	d.Set("konnect_mapping_enabled", c.KonnectMappingEnabled)
}

func resourceAuthenticationSettingsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	err := resourceAuthenticationSettingsCreateUpdate(ctx, d, m)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	return diags
}

func resourceAuthenticationSettingsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	if d.Id() != client.AuthenticationSettingsPath {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("Id must be equal to %s", client.AuthenticationSettingsPath))
	}
	requestPath := fmt.Sprintf(client.AuthenticationSettingsPath)
	body, err := c.HttpRequest(ctx, false, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal := &client.AuthenticationSettings{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	fillResourceDataFromAuthenticationSettings(retVal, d)
	return diags
}

func resourceAuthenticationSettingsCreateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) error {
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	upAuthenticationSettings := client.AuthenticationSettings{}
	fillAuthenticationSettings(&upAuthenticationSettings, d)
	err := json.NewEncoder(&buf).Encode(upAuthenticationSettings)
	if err != nil {
		return err
	}
	requestPath := fmt.Sprintf(client.AuthenticationSettingsPath)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, false, http.MethodPatch, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		return err
	}
	retVal := &client.AuthenticationSettings{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		return err
	}
	fillResourceDataFromAuthenticationSettings(retVal, d)
	return nil
}

func resourceAuthenticationSettingsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	err := resourceAuthenticationSettingsCreateUpdate(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceAuthenticationSettingsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.SetId("")
	return diag.Diagnostics{}
}
