package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	destinations "github.com/jmtx1020/go_quicknode/api/destination"
	"github.com/jmtx1020/go_quicknode/client"
)

var (
	_ datasource.DataSource              = &destinationDataSource{}
	_ datasource.DataSourceWithConfigure = &destinationDataSource{}
)

func NewDestinationDataSource() datasource.DataSource {
	return &destinationDataSource{}
}

type destinationDataSource struct {
	client *client.APIWrapper
}

func (d *destinationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_destination"
}

func (d *destinationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "ID given by API for the destination.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "User supplied name given to the destination.",
				Computed:    true,
			},
			"to": schema.StringAttribute{
				Description: "The webhook URL to which QuickAlerts will send alert payloads.",
				Computed:    true,
			},
			"webhook_type": schema.StringAttribute{
				Description: "The type of destination. ENUM: 'POST', 'GET'",
				Computed:    true,
			},
			"service": schema.StringAttribute{
				Description: "The destination service. Currently only \"webhook\" is supported.",
				Computed:    true,
			},
			"token": schema.StringAttribute{
				Description: "The token for this destination. This is used to optionally verify a QuickAlerts payload.",
				Computed:    true,
			},
			"payload_type": schema.Int64Attribute{
				Description: "The type of payload to send. ENUM: 1,2,3,4,5,6,7",
				Computed:    true,
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

func (d *destinationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state destinationResourceModel
	// idPath := path.Root("id")
	req.Config.GetAttribute(ctx, path.Root("id"), &state.ID)

	destinationAPI := &destinations.DestinationAPI{API: d.client}
	dest, err := destinationAPI.GetDestinationByID(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read QuickNode Destination",
			err.Error())
	}

	state.ID = types.StringValue(dest.ID)
	state.Name = types.StringValue(dest.Name)
	state.To = types.StringValue(dest.To)
	state.Token = types.StringValue(dest.Token)
	state.WebhookType = types.StringValue(dest.WebhookType)
	state.Service = types.StringValue(dest.Service)
	state.PayloadType = types.Int64Value(int64(dest.PayloadType))
	state.CreatedAt = types.StringValue(dest.CreatedAt.Format("2006-01-02 15:04:05"))
	state.UpdatedAt = types.StringValue(dest.UpdatedAt.Format("2006-01-02 15:04:05"))

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *destinationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
