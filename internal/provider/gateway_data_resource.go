package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jmtx1020/go_quicknode/api/ipfs/gateway"
	"github.com/jmtx1020/go_quicknode/client"
)

var (
	_ datasource.DataSource              = &gatewayDataSource{}
	_ datasource.DataSourceWithConfigure = &gatewayDataSource{}
)

func NewGatewayDataSource() datasource.DataSource {
	return &gatewayDataSource{}
}

type gatewayDataSource struct {
	client *client.APIWrapper
}

func (g *gatewayDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gateway"
}

func (g *gatewayDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "An integer that represents the unique identifier of a specified gateway.",
				Computed:    true,
			},
			"uuid": schema.StringAttribute{
				Description: `A string that represents the universally unique identifier (UUID) of the new dedicated gateway.
				UUIDs are used to identify resources uniquely.`,
				Computed: true,
			},
			"name": schema.StringAttribute{
				Description: "A string that specifies the name of the specified gateway. It is a human-readable identifier for the gateway.",
				Required:    true,
			},
			"domain": schema.StringAttribute{
				Description: "The domain associated with the gateway.",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "The status of the gateway.",
				Computed:    true,
			},
			"private": schema.BoolAttribute{
				Description: `A boolean value that indicates whether the specified gateway is private or not.
				If set to true, the gateway is private and not publicly accessible.
				If set to false, the gateway is public and can be accessed by authorized users isEnabled.`,
				Computed: true,
			},
			"enabled": schema.BoolAttribute{
				Description: `A boolean value that indicates whether the specified gateway is enabled or not.
				If set to true, it means the gateway is currently enabled and operational.
				If set to false, it means the gateway is disabled and not functioning`,
				Computed: true,
			},
			"created_at": schema.StringAttribute{
				Description: "The date and time the destination was created.",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "The date and time the destination was last updated.",
				Computed:    true,
			},
		},
	}
}

func (g *gatewayDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state gatewayResourceModel
	req.Config.GetAttribute(ctx, path.Root("name"), &state.Name)

	gatewayAPI := &gateway.GatewayAPI{API: g.client}
	gateway, err := gatewayAPI.GetGetwayByName(state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read QuickNode Gateway",
			err.Error())
	}

	state.ID = types.StringValue(fmt.Sprintf("%v", gateway.ID))
	state.UUID = types.StringValue(gateway.UUID)
	state.Domain = types.StringValue(gateway.Domain)
	state.Status = types.StringValue(gateway.Status)
	state.IsEnabled = types.BoolValue(gateway.IsEnabled)
	state.IsPrivate = types.BoolValue(gateway.IsPrivate)
	state.CreatedAt = types.StringValue(gateway.CreatedAT.Format("2006-01-02 15:04:05"))
	state.UpdatedAt = types.StringValue(gateway.UpdatedAt.Format("2006-01-02 15:04:05"))

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (g *gatewayDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	apiWrapper, ok := req.ProviderData.(*client.APIWrapper)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.APIWrapper, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	g.client = apiWrapper
}
