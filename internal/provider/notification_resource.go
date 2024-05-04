package provider

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/jmtx1020/go_quicknode/api/notifications"
	"github.com/jmtx1020/go_quicknode/client"
)

var (
	_ resource.Resource                = &notificationResource{}
	_ resource.ResourceWithConfigure   = &notificationResource{}
	_ resource.ResourceWithImportState = &notificationResource{}
)

type notificationResource struct {
	client *client.APIWrapper
}

func NewNotificationResource() resource.Resource {
	return &notificationResource{}
}

type notificationResourceModel struct {
	ID             types.String   `tfsdk:"id"`
	Name           types.String   `tfsdk:"name"`
	Expression     types.String   `tfsdk:"expression"`
	Network        types.String   `tfsdk:"network"`
	Enabled        types.Bool     `tfsdk:"enabled"`
	DestinationIDs []types.String `tfsdk:"destination_ids"`
	// Destinations   []destinationResourceModel `tfsdk:"destinations"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

// Configure adds the provider configured client to the resource.
func (n *notificationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Metadata returns the resource type name.
func (n *notificationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification"
}

func (n *notificationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
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
			"destination_ids": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (n *notificationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan notificationResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert []types.String to []string
	destinationIds := make([]string, len(plan.DestinationIDs))
	for i, dest := range plan.DestinationIDs {
		destinationIds[i] = dest.ValueString()
	}

	notificationsAPI := &notifications.NotificationAPI{API: n.client}
	notification, err := notificationsAPI.CreateNotification(
		plan.Name.ValueString(),
		plan.Expression.ValueString(),
		plan.Network.ValueString(),
		destinationIds,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating destination",
			"Could not create destination, unexpected error: "+err.Error(),
		)
		return
	}

	// toggle the notification (enabled or disabled) based on plan values
	if plan.Enabled.ValueBool() == true {
		tflog.Debug(ctx, "Enabling Notification")
		err = notificationsAPI.ToggleNotificationByID(
			notification.ID,
			true,
		)
	} else {
		tflog.Debug(ctx, "Disabling Notification")
		err = notificationsAPI.ToggleNotificationByID(
			notification.ID,
			false,
		)
	}

	plan.ID = types.StringValue(notification.ID)
	plan.CreatedAt = types.StringValue(notification.CreatedAt.Format("2006-01-02 15:04:05"))
	plan.UpdatedAt = types.StringValue(notification.UpdatedAt.Format("2006-01-02 15:04:05"))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (n *notificationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state notificationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	notificationsAPI := &notifications.NotificationAPI{API: n.client}
	notif, err := notificationsAPI.GetNotificationByID(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading QuickNode Notification",
			"Could not read QuickNode ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	destinationIds := make([]types.String, len(notif.Destinations))
	for i, dest := range notif.Destinations {
		destinationIds[i] = types.StringValue(dest.ID)
	}

	notif_bytes := []byte(notif.Expression)
	notif_expr_b64 := base64.StdEncoding.EncodeToString(notif_bytes)

	tflog.Debug(ctx, "READ METHOD")

	state.ID = types.StringValue(notif.ID)
	state.Name = types.StringValue(notif.Name)
	state.Enabled = types.BoolValue(notif.Enabled)
	state.Expression = types.StringValue(notif_expr_b64)
	state.Network = types.StringValue(notif.Network)
	state.DestinationIDs = destinationIds

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (n *notificationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan notificationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state notificationResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	destinationIDs := make([]string, len(plan.DestinationIDs))
	for i, id := range plan.DestinationIDs {
		destinationIDs[i] = id.ValueString()
	}

	notificationsAPI := &notifications.NotificationAPI{API: n.client}
	notif, err := notificationsAPI.UpdateNotificationByID(
		state.ID.ValueString(),
		plan.Name.ValueString(),
		plan.Expression.ValueString(),
		destinationIDs,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating QuickNode Notification.",
			"Could not update QuickNode ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// toggle the notification (enabled or disabled) based on plan values
	if plan.Enabled.ValueBool() == true {
		tflog.Debug(ctx, "Enabling Notification")
		err = notificationsAPI.ToggleNotificationByID(
			notif.ID,
			true,
		)
	} else {
		tflog.Debug(ctx, "Disabling Notification")
		err = notificationsAPI.ToggleNotificationByID(
			notif.ID,
			false,
		)
	}

	notif_bytes := []byte(notif.Expression)
	notif_expr_b64 := base64.StdEncoding.EncodeToString(notif_bytes)

	tflog.Debug(ctx, "THIS IS ORIGINAL PLAN DEBUG: "+plan.Expression.ValueString())
	tflog.Debug(ctx, "THIS IS FROM API PLAN DEBUG: "+notif.Expression)
	tflog.Debug(ctx, "THIS IS FROM API CONVERTED TO BASE64"+notif_expr_b64)

	plan.ID = types.StringValue(notif.ID)
	plan.Name = types.StringValue(notif.Name)
	plan.CreatedAt = types.StringValue(notif.CreatedAt.Format("2006-01-02 15:04:05"))
	plan.UpdatedAt = types.StringValue(notif.UpdatedAt.Format("2006-01-02 15:04:05"))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (n *notificationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state notificationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	notificationsAPI := &notifications.NotificationAPI{API: n.client}
	err := notificationsAPI.DeleteNotificationByID(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting QuickNode Notification"+state.ID.ValueString(),
			"Could not delete destination, unexpected error: "+err.Error(),
		)
		return
	}
}

func (n *notificationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
