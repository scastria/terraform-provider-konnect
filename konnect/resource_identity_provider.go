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
			"client_secret": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"scopes": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"email_claim_mapping": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name_claim_mapping": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"groups_claim_mapping": {
				Type:     schema.TypeString,
				Optional: true,
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
	clientSecret, ok := d.GetOk("client_secret")
	if ok {
		c.ClientSecret = clientSecret.(string)
	}
	scopes, ok := d.GetOk("scopes")
	if ok {
		c.Scopes = convertSetToArray(scopes.(*schema.Set))
	}
	c.ClaimMappings = map[string]string{}
	emailMapping, ok := d.GetOk("email_claim_mapping")
	if ok {
		c.ClaimMappings[client.EmailClaim] = emailMapping.(string)
	}
	nameMapping, ok := d.GetOk("name_claim_mapping")
	if ok {
		c.ClaimMappings[client.NameClaim] = nameMapping.(string)
	}
	groupsMapping, ok := d.GetOk("groups_claim_mapping")
	if ok {
		c.ClaimMappings[client.GroupsClaim] = groupsMapping.(string)
	}
}

func fillResourceDataFromIdentityProvider(c *client.IdentityProvider, d *schema.ResourceData) {
	d.Set("issuer", c.Issuer)
	d.Set("login_path", c.LoginPath)
	d.Set("client_id", c.ClientId)
	//Do not set client_secret in state since it can never be read back.  Let previous value propagate forward for No Changes
	//d.Set("client_secret", c.ClientSecret)
	d.Set("scopes", c.Scopes)
	if c.ClaimMappings == nil {
		d.Set("email_claim_mapping", "")
		d.Set("name_claim_mapping", "")
		d.Set("groups_claim_mapping", "")
	} else {
		d.Set("email_claim_mapping", c.ClaimMappings[client.EmailClaim])
		d.Set("name_claim_mapping", c.ClaimMappings[client.NameClaim])
		d.Set("groups_claim_mapping", c.ClaimMappings[client.GroupsClaim])
	}
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
