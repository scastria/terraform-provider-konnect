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
	"regexp"
	"strconv"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserRead,
		Schema: map[string]*schema.Schema{
			"search_email": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"email": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`), "must be a valid email address"),
			},
			"search_full_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"full_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"active": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"preferred_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	requestQuery := url.Values{}
	searchEmail, ok := d.GetOk("search_email")
	if ok {
		requestQuery[client.FilterEmailContains] = []string{searchEmail.(string)}
	}
	email, ok := d.GetOk("email")
	if ok {
		requestQuery[client.FilterEmail] = []string{email.(string)}
	}
	searchFullName, ok := d.GetOk("search_full_name")
	if ok {
		requestQuery[client.FilterFullNameContains] = []string{searchFullName.(string)}
	}
	fullName, ok := d.GetOk("full_name")
	if ok {
		requestQuery[client.FilterFullName] = []string{fullName.(string)}
	}
	active := d.Get("active").(bool)
	requestQuery[client.FilterActive] = []string{strconv.FormatBool(active)}
	body, err := c.HttpRequest(ctx, false, http.MethodGet, client.UserPath, requestQuery, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVals := &client.UserCollection{}
	err = json.NewDecoder(body).Decode(retVals)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	numUsers := len(retVals.Users)
	if numUsers > 1 {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("Filter criteria does not result in a single user"))
	} else if numUsers != 1 {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("No user exists with that filter criteria"))
	}
	retVal := retVals.Users[0]
	d.Set("email", retVal.Email)
	d.Set("full_name", retVal.FullName)
	d.Set("preferred_name", retVal.PreferredName)
	d.Set("active", retVal.Active)
	d.SetId(retVal.Id)
	return diags
}
