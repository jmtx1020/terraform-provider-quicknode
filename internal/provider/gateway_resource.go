package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	gateways "github.com/jmtx1020/go_quicknode/api/ipfs/gateway"
	"github.com/jmtx1020/go_quicknode/client"
)

var (
	_ resource.Resource                = &gatewayResource{}
	_ resource.ResourceWithConfigure   = &gatewayResource{}
	_ resource.ResourceWithImportState = &gatewayResource{}
)

type gatewayResource struct {
	client *client.APIWrapper
}

func NewGatewayResource() resource.Resource {
	return &gatewayResource{}
}

type gatewayResourceModel struct {
	ID        types.String `tfsdk:"id"`
	UUID      types.String `tfsdk:"uuid"`
	Name      types.String `tfsdk:"name"`
	Domain    types.String `tfsdk:"domain"`
	Status    types.String `tfsdk:"status"`
	IsPrivate types.Bool   `tfsdk:"private"`
	IsEnabled types.Bool   `tfsdk:"enabled"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

// Configure adds the provider configured client to the resource.
func (g *gatewayResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (g *gatewayResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gateway"
}

func (g *gatewayResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
				Required: true,
			},
			"enabled": schema.BoolAttribute{
				Description: `A boolean value that indicates whether the specified gateway is enabled or not.
				If set to true, it means the gateway is currently enabled and operational.
				If set to false, it means the gateway is disabled and not functioning`,
				Required: true,
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

func (g *gatewayResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan gatewayResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	gatewayAPI := &gateways.GatewayAPI{API: g.client}
	gateway, err := gatewayAPI.CreateGateway(
		plan.Name.ValueString(),
		plan.IsPrivate.ValueBool(),
		plan.IsEnabled.ValueBool(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating ipfs gateway",
			"Could not create ipfs gateway, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%v", gateway.ID))
	plan.Domain = types.StringValue(gateway.Domain)
	plan.UUID = types.StringValue(gateway.UUID)
	plan.IsEnabled = types.BoolValue(gateway.IsEnabled)
	plan.IsPrivate = types.BoolValue(gateway.IsPrivate)
	plan.Status = types.StringValue(gateway.Status)
	plan.CreatedAt = types.StringValue(gateway.CreatedAT.Format("2006-01-02 15:04:05"))
	plan.UpdatedAt = types.StringValue(gateway.UpdatedAt.Format("2006-01-02 15:04:05"))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (g *gatewayResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state gatewayResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	gatewayAPI := &gateways.GatewayAPI{API: g.client}
	gateway, err := gatewayAPI.GetGetwayByName(state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading QuickNode Gateway",
			"Could not read QuickNode ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	state.ID = types.StringValue(fmt.Sprintf("%v", gateway.ID))
	state.UUID = types.StringValue(gateway.UUID)
	state.Domain = types.StringValue(gateway.Domain)
	state.Status = types.StringValue(gateway.Status)
	state.IsEnabled = types.BoolValue(gateway.IsEnabled)
	state.IsPrivate = types.BoolValue(gateway.IsPrivate)
	state.CreatedAt = types.StringValue(gateway.CreatedAT.Format("2006-01-02 15:04:05"))
	state.UpdatedAt = types.StringValue(gateway.UpdatedAt.Format("2006-01-02 15:04:05"))

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (g *gatewayResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan gatewayResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state gatewayResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if the field that cannot be updated has changed
	if state.Name != plan.Name {
		// Add a diagnostic to indicate that the field cannot be updated
		resp.Diagnostics.AddError("Name field cannot be updated", "The 'Name' field cannot be updated in the API")

		// Return from the method without making any changes to the API
		return
	}

	gatewayAPI := &gateways.GatewayAPI{API: g.client}
	gateway, err := gatewayAPI.UpdateGatewayByName(
		state.Name.ValueString(),
		plan.IsPrivate.ValueBool(),
		plan.IsEnabled.ValueBool(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating QuickNode Gateway.",
			"Could not update gateway QuickNode: "+state.Name.ValueString()+": "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%v", gateway.ID))
	plan.UUID = types.StringValue(gateway.UUID)
	plan.Name = types.StringValue(gateway.Name)
	plan.Domain = types.StringValue(gateway.Domain)
	plan.Status = types.StringValue(gateway.Status)
	plan.IsEnabled = types.BoolValue(gateway.IsEnabled)
	plan.IsPrivate = types.BoolValue(gateway.IsPrivate)
	plan.CreatedAt = types.StringValue(gateway.CreatedAT.Format("2006-01-02 15:04:05"))
	plan.UpdatedAt = types.StringValue(gateway.UpdatedAt.Format("2006-01-02 15:04:05"))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (g *gatewayResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state gatewayResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	gatewayAPI := &gateways.GatewayAPI{API: g.client}
	err := gatewayAPI.DeleteGatewayByName(state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting QuickNode Gateway"+state.ID.ValueString(),
			"Could not delete gateway, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *gatewayResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
