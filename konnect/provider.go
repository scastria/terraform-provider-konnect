package konnect

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/scastria/terraform-provider-konnect/konnect/client"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"pat": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("KONNECT_PAT", nil),
			},
			"region": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("KONNECT_REGION", "us"),
				ValidateFunc: validation.StringInSlice([]string{"us", "eu"}, false),
			},
			"default_tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"konnect_runtime_group":          resourceRuntimeGroup(),
			"konnect_authentication_setting": resourceAuthenticationSettings(),
			"konnect_identity_provider":      resourceIdentityProvider(),
			"konnect_user":                   resourceUser(),
			"konnect_team":                   resourceTeam(),
			"konnect_team_user":              resourceTeamUser(),
			"konnect_team_role":              resourceTeamRole(),
			"konnect_user_role":              resourceUserRole(),
			"konnect_service":                resourceService(),
			"konnect_route":                  resourceRoute(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"konnect_runtime_group": dataSourceRuntimeGroup(),
			"konnect_user":          dataSourceUser(),
			"konnect_team":          dataSourceTeam(),
			"konnect_role":          dataSourceRole(),
			"konnect_team_role":     dataSourceTeamRole(),
			"konnect_user_role":     dataSourceUserRole(),
			"konnect_nodes":         dataSourceNodes(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	pat := d.Get("pat").(string)
	region := d.Get("region").(string)
	//defaultTags := []string{}
	//defaultTagsSet, ok := d.GetOk("default_tags")
	//if ok {
	//	defaultTags = convertSetToArray(defaultTagsSet.(*schema.Set))
	//}

	var diags diag.Diagnostics
	//c, err := client.NewClient(pat, region, defaultTags)
	c, err := client.NewClient(pat, region)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return c, diags
}
