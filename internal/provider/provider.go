package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-konnect/internal/client"
	"os"
)

const (
	DefaultRegion = "us"
)

// Ensure KonnectProvider satisfies various provider interfaces.
var _ provider.Provider = &KonnectProvider{}

// KonnectProvider defines the provider implementation.
type KonnectProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// KonnectProviderModel describes the provider data model.
type KonnectProviderModel struct {
	pat    types.String `tfsdk:"pat"`
	region types.String `tfsdk:"region"`
}

func (p *KonnectProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "konnect"
	resp.Version = p.version
}

func (p *KonnectProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"pat": schema.StringAttribute{
				MarkdownDescription: "Personal Access Token. Can be specified via env variable KONNECT_PAT.",
				Required:            true,
				Sensitive:           true,
			},
			"region": schema.StringAttribute{
				MarkdownDescription: "Region used for all region specific resources. Can be specified via env variable KONNECT_REGION. Default: us",
				Optional:            true,
			},
		},
	}
}

func (p *KonnectProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config KonnectProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	pat := os.Getenv("KONNECT_PAT")
	region := os.Getenv("KONNECT_REGION")
	if region == "" {
		region = DefaultRegion
	}
	if !config.pat.IsNull() {
		pat = config.pat.ValueString()
	}
	if !config.region.IsNull() {
		region = config.region.ValueString()
	}

	if pat == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("pat"),
			"Missing pat",
			"pat is required",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	c, err := client.NewClient(ctx, pat, region)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Konnect http client",
			"Unable to create Konnect http client",
		)
	}
	resp.DataSourceData = c
	resp.ResourceData = c
}

func (p *KonnectProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		//NewExampleResource,
	}
}

func (p *KonnectProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		//NewExampleDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &KonnectProvider{
			version: version,
		}
	}
}
