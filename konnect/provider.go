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
				ValidateFunc: validation.StringInSlice([]string{"us", "eu", "au"}, false),
			},
			"num_retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("KONNECT_NUM_RETRIES", 3),
			},
			"retry_delay": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("KONNECT_RETRY_DELAY", 30),
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
			"konnect_consumer_basic":          resourceConsumerBasic(),
			"konnect_consumer_hmac":           resourceConsumerHMAC(),
			"konnect_consumer_jwt":            resourceConsumerJWT(),
			"konnect_plugin":                  resourcePlugin(),
			"konnect_custom_plugin_schema":    resourceCustomPluginSchema(),
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
			"konnect_service":       dataSourceService(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	pat := d.Get("pat").(string)
	region := d.Get("region").(string)
	numRetries := d.Get("num_retries").(int)
	retryDelay := d.Get("retry_delay").(int)
	defaultTags := []string{}
	defaultTagsSet, ok := d.GetOk("default_tags")
	if ok {
		defaultTags = convertSetToArray(defaultTagsSet.(*schema.Set))
	}

	var diags diag.Diagnostics
	c, err := client.NewClient(pat, region, numRetries, retryDelay, defaultTags)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return c, diags
}
