package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jmtx1020/go_quicknode/api/notifications"
	"github.com/jmtx1020/go_quicknode/client"
)

var (
	_ datasource.DataSource              = &notificationDataSource{}
	_ datasource.DataSourceWithConfigure = &notificationDataSource{}
)

type notificationDataResourceModel struct {
	ID           types.String               `tfsdk:"id"`
	Name         types.String               `tfsdk:"name"`
	Expression   types.String               `tfsdk:"expression"`
	Network      types.String               `tfsdk:"network"`
	Enabled      types.Bool                 `tfsdk:"enabled"`
	Destinations []destinationResourceModel `tfsdk:"destinations"`
	CreatedAt    types.String               `tfsdk:"created_at"`
	UpdatedAt    types.String               `tfsdk:"updated_at"`
}

func NewNotificationDataSource() datasource.DataSource {
	return &notificationDataSource{}
}

type notificationDataSource struct {
	client *client.APIWrapper
}

func (n *notificationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification"
}

func (n *notificationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The notification ID.",
				Required:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the notification is enabled.",
				Computed:    true,
			},
			"expression": schema.StringAttribute{
				Description: "The expression for the notification.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the notification.",
				Computed:    true,
			},
			"network": schema.StringAttribute{
				Description: "The webhook URL to which QuickAlerts will send alert payloads.",
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
			"destinations": schema.ListNestedAttribute{
				Description: "The destinations for the notification returned as arrays.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The destination ID.",
							Computed:    true,
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
				},
			},
		},
	}
}

func (n *notificationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state notificationDataResourceModel

	req.Config.GetAttribute(ctx, path.Root("id"), &state.ID)
	notificationsAPI := &notifications.NotificationAPI{API: n.client}

	notif, err := notificationsAPI.GetNotificationByID(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read QuickNode Notification",
			err.Error())
	}

	var destinationModels []destinationResourceModel
	for _, dest := range notif.Destinations {
		destModel := destinationResourceModel{
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
		destinationModels = append(destinationModels, destModel)
	}

	state.ID = types.StringValue(notif.ID)
	state.Name = types.StringValue(notif.Name)
	state.Network = types.StringValue(notif.Network)
	state.Expression = types.StringValue(notif.Expression)
	state.Enabled = types.BoolValue(notif.Enabled)
	state.Destinations = destinationModels
	state.CreatedAt = types.StringValue(notif.CreatedAt.Format("2006-01-02 15:04:05"))
	state.UpdatedAt = types.StringValue(notif.UpdatedAt.Format("2006-01-02 15:04:05"))

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (n *notificationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	n.client = apiWrapper
}
