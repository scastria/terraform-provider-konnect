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

func resourceIdentityProvider() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityProviderCreate,
		ReadContext:   resourceIdentityProviderRead,
		UpdateContext: resourceIdentityProviderUpdate,
		DeleteContext: resourceIdentityProviderDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"issuer": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"login_path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scopes": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"claim_mappings": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func fillIdentityProvider(c *client.IdentityProvider, d *schema.ResourceData) {
	issuer, ok := d.GetOk("issuer")
	if ok {
		c.Issuer = issuer.(string)
	}
	loginPath, ok := d.GetOk("login_path")
	if ok {
		c.LoginPath = loginPath.(string)
	}
	clientId, ok := d.GetOk("client_id")
	if ok {
		c.ClientId = clientId.(string)
	}
	scopes, ok := d.GetOk("scopes")
	if ok {
		c.Scopes = convertSetToArray(scopes.(*schema.Set))
	}
	cm, ok := d.GetOk("claim_mappings")
	if ok {
		if c.ClaimMappings == nil {
			c.ClaimMappings = map[string]string{}
		}
		claimMappings := cm.(map[string]interface{})
		for name, value := range claimMappings {
			c.ClaimMappings[name] = value.(string)
		}
	}
}

func fillResourceDataFromIdentityProvider(c *client.IdentityProvider, d *schema.ResourceData) {
	d.Set("issuer", c.Issuer)
	d.Set("login_path", c.LoginPath)
	d.Set("client_id", c.ClientId)
	d.Set("scopes", c.Scopes)
	d.Set("claim_mappings", c.ClaimMappings)
}

func resourceIdentityProviderCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	err := resourceIdentityProviderCreateUpdate(ctx, d, m)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	return diags
}

func resourceIdentityProviderRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	if d.Id() != client.IdentityProviderPath {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("Id must be equal to %s", client.IdentityProviderPath))
	}
	requestPath := fmt.Sprintf(client.IdentityProviderPath)
	body, err := c.HttpRequest(ctx, false, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal := &client.IdentityProvider{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	fillResourceDataFromIdentityProvider(retVal, d)
	return diags
}

func resourceIdentityProviderCreateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) error {
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	upIdentityProvider := client.IdentityProvider{}
	fillIdentityProvider(&upIdentityProvider, d)
	err := json.NewEncoder(&buf).Encode(upIdentityProvider)
	if err != nil {
		return err
	}
	requestPath := fmt.Sprintf(client.IdentityProviderPath)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, false, http.MethodPatch, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		return err
	}
	retVal := &client.IdentityProvider{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		return err
	}
	fillResourceDataFromIdentityProvider(retVal, d)
	return nil
}

func resourceIdentityProviderUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	err := resourceIdentityProviderCreateUpdate(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceIdentityProviderDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.SetId("")
	return diag.Diagnostics{}
}
