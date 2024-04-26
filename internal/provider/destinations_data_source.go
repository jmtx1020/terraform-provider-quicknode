package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	destinations "github.com/jmtx1020/go_quicknode/api/destination"
	"github.com/jmtx1020/go_quicknode/client"
)

var (
	_ datasource.DataSource              = &destinationsDataSource{}
	_ datasource.DataSourceWithConfigure = &destinationsDataSource{}
)

func NewDestinationsDataSource() datasource.DataSource {
	return &destinationsDataSource{}
}

type destinationsDataSource struct {
	client *client.APIWrapper
}

type destinationsDataSourceModel struct {
	Destinations []destinationModel `tfsdk:"destinations"`
}

type destinationModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	To          types.String `tfsdk:"to"`
	WebhookType types.String `tfsdk:"webhook_type"`
	Service     types.String `tfsdk:"service"`
	Token       types.String `tfsdk:"token"`
	PayloadType types.Int64  `tfsdk:"payload_type"`
	CreatedAt   types.String `tfsdk:"created_at"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
}

func (d *destinationsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_destinations"
}

func (d *destinationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"destinations": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"to": schema.StringAttribute{
							Computed: true,
						},
						"webhook_type": schema.StringAttribute{
							Computed: true,
						},
						"service": schema.StringAttribute{
							Computed: true,
						},
						"token": schema.StringAttribute{
							Computed: true,
						},
						"payload_type": schema.Int64Attribute{
							Computed: true,
						},
						"created_at": schema.StringAttribute{
							Computed: true,
						},
						"updated_at": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (d *destinationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state destinationsDataSourceModel

	destinationAPI := &destinations.DestinationAPI{API: d.client}

	dests, err := destinationAPI.GetAllDestinations()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read HashiCups Coffees",
			err.Error())
	}

	for _, dest := range dests {
		destState := destinationModel{
			ID:          types.StringValue(dest.ID),
			Name:        types.StringValue(dest.Name),
			To:          types.StringValue(dest.To),
			WebhookType: types.StringValue(dest.WebhookType),
			Service:     types.StringValue(dest.Service),
			Token:       types.StringValue(dest.Token),
			PayloadType: types.Int64Value(int64(dest.PayloadType)),
			CreatedAt:   types.StringValue(dest.CreatedAt.Format("2006-01-02 15:04:05")),
			UpdatedAt:   types.StringValue(dest.UpdatedAt.Format("2006-01-02 15:04:05")),
		}

		state.Destinations = append(state.Destinations, destState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *destinationsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = apiWrapper
}
