package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	notifications "github.com/jmtx1020/go_quicknode/api/notifications"
	"github.com/jmtx1020/go_quicknode/client"
)

var (
	_ datasource.DataSource              = &notificationsDataSource{}
	_ datasource.DataSourceWithConfigure = &notificationsDataSource{}
)

type notificationsDataSource struct {
	client *client.APIWrapper
}

func NewNotificationsDataSource() datasource.DataSource {
	return &notificationsDataSource{}
}

type notificationsDataSourceModel struct {
	Notifications []notificationsModel `tfsdk:"notifications"`
}

type notificationsModel struct {
	ID           types.String       `tfsdk:"id"`
	Name         types.String       `tfsdk:"name"`
	Expression   types.String       `tfsdk:"expression"`
	Network      types.String       `tfsdk:"network"`
	Enabled      types.Bool         `tfsdk:"enabled"`
	Destinations []destinationModel `tfsdk:"destinations"`
	CreatedAt    types.String       `tfsdk:"created_at"`
	UpdatedAt    types.String       `tfsdk:"updated_at"`
}

func (n *notificationsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notifications"
}

func (d *notificationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"notifications": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The notification ID.",
							Computed:    true,
						},
						"expression": schema.StringAttribute{
							Description: "The expression for the notification.",
							Required:    true,
						},
						"enabled": schema.BoolAttribute{
							Description: "Whether the notification is enabled.",
							Required:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of the notification.",
							Required:    true,
						},
						"network": schema.StringAttribute{
							Description: "The webhook URL to which QuickAlerts will send alert payloads.",
							Required:    true,
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
				},
			},
		},
	}
}

func (n *notificationsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (n *notificationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state notificationsDataSourceModel

	notificationsAPI := &notifications.NotificationAPI{API: n.client}
	notifications, err := notificationsAPI.GetAllNotifications()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read QuickNode Notifications",
			err.Error())
	}

	for _, notification := range notifications {
		notificationState := notificationsModel{
			ID:         types.StringValue(notification.ID),
			Name:       types.StringValue(notification.Name),
			Expression: types.StringValue(notification.Expression),
			Network:    types.StringValue(notification.Network),
			Enabled:    types.BoolValue(notification.Enabled),
			CreatedAt:  types.StringValue(notification.CreatedAt.Format("2006-01-02 15:04:05")),
			UpdatedAt:  types.StringValue(notification.UpdatedAt.Format("2006-01-02 15:04:05")),
		}

		for _, dest := range notification.Destinations {
			destinationState := destinationModel{
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

			notificationState.Destinations = append(notificationState.Destinations, destinationState)
		}

		state.Notifications = append(state.Notifications, notificationState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
