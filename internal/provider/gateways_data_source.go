package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	gateways "github.com/jmtx1020/go_quicknode/api/ipfs/gateway"
	"github.com/jmtx1020/go_quicknode/client"
)

var (
	_ datasource.DataSource              = &gatewaysDataSource{}
	_ datasource.DataSourceWithConfigure = &gatewaysDataSource{}
)

func NewGatewaysDataSource() datasource.DataSource {
	return &gatewaysDataSource{}
}

type gatewaysDataSource struct {
	client *client.APIWrapper
}

type gatewaysDataSourceModel struct {
	Gateways []gatewayResourceModel `tfsdk:"gateways"`
}

func (g *gatewaysDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gateways"
}

func (g *gatewaysDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"gateways": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
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
							Computed:    true,
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
				},
			},
		},
	}
}

func (g *gatewaysDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state gatewaysDataSourceModel

	gatewayAPI := &gateways.GatewayAPI{API: g.client}
	gateways, err := gatewayAPI.GetAllGateways()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read QuickNode Destinations",
			err.Error())
	}

	for _, gateway := range gateways {
		gatewayState := gatewayResourceModel{
			ID:        types.StringValue(fmt.Sprintf("%v", gateway.ID)),
			Domain:    types.StringValue(gateway.Domain),
			Name:      types.StringValue(gateway.Name),
			UUID:      types.StringValue(gateway.UUID),
			IsEnabled: types.BoolValue(gateway.IsEnabled),
			IsPrivate: types.BoolValue(gateway.IsPrivate),
			Status:    types.StringValue(gateway.Status),
			CreatedAt: types.StringValue(gateway.CreatedAT.Format("2006-01-02 15:04:05")),
			UpdatedAt: types.StringValue(gateway.UpdatedAt.Format("2006-01-02 15:04:05")),
		}

		state.Gateways = append(state.Gateways, gatewayState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (g *gatewaysDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
