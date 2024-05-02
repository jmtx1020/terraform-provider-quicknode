package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	destinations "github.com/jmtx1020/go_quicknode/api/destinations"
	"github.com/jmtx1020/go_quicknode/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &destinationResource{}
	_ resource.ResourceWithConfigure   = &destinationResource{}
	_ resource.ResourceWithImportState = &destinationResource{}
)

type destinationResource struct {
	client *client.APIWrapper
}

func NewDestinationResource() resource.Resource {
	return &destinationResource{}
}

type destinationResourceModel struct {
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

// Configure adds the provider configured client to the resource.
func (r *destinationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = apiWrapper
}

// Metadata returns the resource type name.
func (r *destinationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_destination"
}

// Schema defines the schema for the resource.
func (r *destinationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "ID given by API for the destination.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "User supplied name given to the destination.",
				Required:    true,
			},
			"to": schema.StringAttribute{
				Description: "The webhook URL to which QuickAlerts will send alert payloads.",
				Required:    true,
			},
			"webhook_type": schema.StringAttribute{
				Description: "The type of destination. ENUM: 'POST', 'GET'",
				Required:    true,
			},
			"service": schema.StringAttribute{
				Description: "The destination service. Currently only \"webhook\" is supported.",
				Required:    true,
			},
			"token": schema.StringAttribute{
				Description: "The token for this destination. This is used to optionally verify a QuickAlerts payload.",
				Computed:    true,
			},
			"payload_type": schema.Int64Attribute{
				Description: "The type of payload to send. ENUM: 1,2,3,4,5,6,7",
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
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *destinationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan destinationResourceModel
	diags := req.Plan.Get(ctx, &plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	destinationsAPI := &destinations.DestinationAPI{API: r.client}
	dest, err := destinationsAPI.CreateDestination(
		plan.Name.ValueString(),
		plan.To.ValueString(),
		plan.WebhookType.ValueString(),
		plan.Service.ValueString(),
		int(plan.PayloadType.ValueInt64()),
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating destination",
			"Could not create destination, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(dest.ID)
	plan.Token = types.StringValue(dest.Token)
	plan.CreatedAt = types.StringValue(dest.CreatedAt.Format("2006-01-02 15:04:05"))
	plan.UpdatedAt = types.StringValue(dest.UpdatedAt.Format("2006-01-02 15:04:05"))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *destinationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state destinationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	destinationsAPI := &destinations.DestinationAPI{API: r.client}
	// Get refreshed order value from HashiCups
	dest, err := destinationsAPI.GetDestinationByID(fmt.Sprint(state.ID.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading QuickNode Destination",
			"Could not read QuickNode ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	state.ID = types.StringValue(dest.ID)
	state.Token = types.StringValue(dest.Token)
	state.Name = types.StringValue(dest.Name)
	state.To = types.StringValue(dest.To)
	state.WebhookType = types.StringValue(dest.WebhookType)
	state.Service = types.StringValue(dest.Service)
	state.PayloadType = types.Int64Value(int64(dest.PayloadType))
	state.CreatedAt = types.StringValue(dest.CreatedAt.Format("2006-01-02 15:04:05"))
	state.UpdatedAt = types.StringValue(dest.UpdatedAt.Format("2006-01-02 15:04:05"))

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *destinationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan destinationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state destinationResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	destinationsAPI := &destinations.DestinationAPI{API: r.client}

	err := destinationsAPI.DeleteDestinationByID(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Update - Error Deleting QuickNode Destination.",
			"Update - Could not delete QuickNode ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	dest, err := destinationsAPI.CreateDestination(
		plan.Name.ValueString(),
		plan.To.ValueString(),
		plan.WebhookType.ValueString(),
		plan.Service.ValueString(),
		int(plan.PayloadType.ValueInt64()),
	)

	plan.ID = types.StringValue(dest.ID)
	plan.Token = types.StringValue(dest.Token)
	plan.CreatedAt = types.StringValue(dest.CreatedAt.Format("2006-01-02 15:04:05"))
	plan.UpdatedAt = types.StringValue(dest.UpdatedAt.Format("2006-01-02 15:04:05"))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *destinationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state destinationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	destinationsAPI := &destinations.DestinationAPI{API: r.client}

	// Delete existing order
	err := destinationsAPI.DeleteDestinationByID(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting QuickNode Destination"+state.ID.ValueString(),
			"Could not delete destination, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *destinationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
