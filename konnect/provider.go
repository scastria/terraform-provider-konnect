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
			//"default_tags": {
			//	Type:     schema.TypeSet,
			//	Optional: true,
			//	Elem: &schema.Schema{
			//		Type: schema.TypeString,
			//	},
			//},
		},
		ResourcesMap: map[string]*schema.Resource{
			"konnect_control_plane":           resourceControlPlane(),
			"konnect_authentication_settings": resourceAuthenticationSettings(),
			"konnect_identity_provider":       resourceIdentityProvider(),
			"konnect_user":                    resourceUser(),
			"konnect_team":                    resourceTeam(),
			"konnect_team_user":               resourceTeamUser(),
			"konnect_team_role":               resourceTeamRole(),
			"konnect_team_mappings":           resourceTeamMappings(),
			"konnect_user_role":               resourceUserRole(),
			"konnect_service":                 resourceService(),
			"konnect_route":                   resourceRoute(),
			"konnect_consumer":                resourceConsumer(),
			"konnect_consumer_key":            resourceConsumerKey(),
			"konnect_consumer_acl":            resourceConsumerACL(),
			"konnect_plugin":                  resourcePlugin(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"konnect_control_plane": dataSourceControlPlane(),
			"konnect_user":          dataSourceUser(),
			"konnect_team":          dataSourceTeam(),
			"konnect_role":          dataSourceRole(),
			"konnect_team_role":     dataSourceTeamRole(),
			"konnect_user_role":     dataSourceUserRole(),
			"konnect_nodes":         dataSourceNodes(),
			"konnect_consumer":      dataSourceConsumer(),
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
