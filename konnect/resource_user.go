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

func resourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"full_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"preferred_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"active": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func fillUser(c *client.User, d *schema.ResourceData) {
	c.Email = d.Get("email").(string)
	c.Active = d.Get("active").(bool)
	fullName, ok := d.GetOk("full_name")
	if ok {
		c.FullName = fullName.(string)
	}
	preferredName, ok := d.GetOk("preferred_name")
	if ok {
		c.PreferredName = preferredName.(string)
	}
}

func fillResourceDataFromUser(c *client.User, d *schema.ResourceData) {
	d.Set("email", c.Email)
	d.Set("full_name", c.FullName)
	d.Set("preferred_name", c.PreferredName)
	d.Set("active", c.Active)
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	newUser := client.User{}
	fillUser(&newUser, d)
	err := json.NewEncoder(&buf).Encode(newUser)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.UserPath)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, false, http.MethodPost, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal := &client.User{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId(retVal.Id)
	fillResourceDataFromUser(retVal, d)
	return diags
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.UserPathGet, d.Id())
	body, err := c.HttpRequest(ctx, false, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.User{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	fillResourceDataFromUser(retVal, d)
	return diags
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	upUser := client.User{}
	fillUser(&upUser, d)
	// Hide non-updateable fields
	upUser.Email = ""
	upUser.Active = false
	err := json.NewEncoder(&buf).Encode(upUser)
	if err != nil {
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.UserPathGet, d.Id())
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, false, http.MethodPatch, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		return diag.FromErr(err)
	}
	retVal := &client.User{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		return diag.FromErr(err)
	}
	fillResourceDataFromUser(retVal, d)
	return diags
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.UserPathGet, d.Id())
	_, err := c.HttpRequest(ctx, false, http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
