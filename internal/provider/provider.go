package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/jmtx1020/go_quicknode/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &quicknodeProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &quicknodeProvider{
			version: version,
		}
	}
}

// quicknodeProvider is the provider implementation.
type quicknodeProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

type quicknodeProviderModel struct {
	Host  types.String `tfsdk:"host"`
	Token types.String `tfsdk:"token"`
}

// Metadata returns the provider type name.
func (p *quicknodeProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "quicknode"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *quicknodeProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Description: "API Hostname",
				Optional:    true,
			},
			"token": schema.StringAttribute{
				Description: "API Token to use to authenticate.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

// Configure prepares a QN API client for data sources and resources.
func (p *quicknodeProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring QuickNode client")
	// retrieve provider data from configuration
	var config quicknodeProviderModel

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// if the practicioner provided a configuration value for an of the attributes
	// it must be a known value
	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown QuickNode API Host",
			"The provider cannot create the QuickNode API Client as there is an unknown configuration value for the QuickNode API host."+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the QUICKNODE_HOST environment variable.")
	}

	if config.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Unknown QuickNode API Password",
			"The provider cannot create the QuickNode API client as there is an unknown configuration value for the QuickNode API token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the HASHICUPS_PASSWORD environment variable.")
	}

	if resp.Diagnostics.HasError() {
		return
	}

	host := os.Getenv("QUICKNODE_API_HOST")
	token := os.Getenv("QUICKNODE_API_TOKEN")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Token.IsNull() {
		token = config.Token.ValueString()
	}

	// if any of the expected configurations are missing, return
	// errors with provider-specific guidance
	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing QuickNode API Host",
			"The provider cannot create the QuickNode API client as there is a missing or empty value for the QuickNode API host. "+
				"Set the host value in the configuration or use the QUICKNODE_API_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.")
	}

	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Missing QuickNode API Token",
			"The provider cannot create the QuickNode API client as there is a missing or empty value for the QuickNode API password. "+
				"Set the password value in the configuration or use the QUICKNODE_API_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.")
	}

	ctx = tflog.SetField(ctx, "quicknode_api_host", host)
	ctx = tflog.SetField(ctx, "quicknode_api_token", token)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "quicknode_api_token")

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating QuickNode client")
	tf_client := client.NewAPIWrapper(token, host)

	// make the quicknode api client available during data source and resource
	resp.DataSourceData = tf_client
	resp.ResourceData = tf_client

	tflog.Info(ctx, "Configured QuickNode client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *quicknodeProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewDestinationsDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *quicknodeProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewDestinationResource,
	}
}
