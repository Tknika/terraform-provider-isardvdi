package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tknika/terraform-provider-isard/internal/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &networkInterfaceResource{}
	_ resource.ResourceWithConfigure = &networkInterfaceResource{}
)

// NewNetworkInterfaceResource is a helper function to simplify the provider implementation.
func NewNetworkInterfaceResource() resource.Resource {
	return &networkInterfaceResource{}
}

// networkInterfaceResource is the resource implementation.
type networkInterfaceResource struct {
	client *client.Client
}

// networkInterfaceResourceModel maps the resource schema data.
type networkInterfaceResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Net         types.String `tfsdk:"net"`
	Kind        types.String `tfsdk:"kind"`
	Model       types.String `tfsdk:"model"`
	QoSID       types.String `tfsdk:"qos_id"`
}

// Metadata returns the resource type name.
func (r *networkInterfaceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_interface"
}

// Schema defines the schema for the resource.
func (r *networkInterfaceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Gestiona una interfaz de red del sistema en Isard VDI (solo administradores).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "ID único de la interfaz de red.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Nombre de la interfaz (ej. 'bridge-100').",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Descripción de la interfaz.",
				Optional:    true,
			},
			"net": schema.StringAttribute{
				Description: "Nombre de la red/bridge del sistema a la que se asocia la interfaz.",
				Required:    true,
			},
			"kind": schema.StringAttribute{
				Description: "Tipo de interfaz (ej. 'bridge', 'network').",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"model": schema.StringAttribute{
				Description: "Modelo de interfaz de red (por defecto: 'virtio').",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"qos_id": schema.StringAttribute{
				Description: "ID del perfil QoS de red (por defecto: 'unlimited').",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *networkInterfaceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Create creates the resource and sets the initial Terraform state.
func (r *networkInterfaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan networkInterfaceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Crear la interfaz de red
	err := r.client.CreateNetworkInterface(
		plan.ID.ValueString(),
		plan.Name.ValueString(),
		plan.Description.ValueString(),
		plan.Net.ValueString(),
		plan.Kind.ValueString(),
		plan.Model.ValueString(),
		plan.QoSID.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creando interfaz de red",
			"No se pudo crear la interfaz de red: "+err.Error(),
		)
		return
	}

	// Obtener la interfaz creada para leer los valores computados
	iface, err := r.client.GetNetworkInterface(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error leyendo interfaz de red creada",
			"No se pudo leer la interfaz recién creada: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	if iface.Kind != "" {
		plan.Kind = types.StringValue(iface.Kind)
	}
	if iface.Model != "" {
		plan.Model = types.StringValue(iface.Model)
	}
	if iface.QoSID != "" {
		plan.QoSID = types.StringValue(iface.QoSID)
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *networkInterfaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state networkInterfaceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed interface value from Isard
	iface, err := r.client.GetNetworkInterface(state.ID.ValueString())
	if err != nil {
		if err.Error() == "network interface not found" {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error leyendo interfaz de red",
			"No se pudo leer la interfaz de red ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	state.Name = types.StringValue(iface.Name)
	state.Net = types.StringValue(iface.Net)
	if iface.Description != "" {
		state.Description = types.StringValue(iface.Description)
	}
	if iface.Kind != "" {
		state.Kind = types.StringValue(iface.Kind)
	}
	if iface.Model != "" {
		state.Model = types.StringValue(iface.Model)
	}
	if iface.QoSID != "" {
		state.QoSID = types.StringValue(iface.QoSID)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *networkInterfaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan networkInterfaceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state
	var state networkInterfaceResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preparar los valores a actualizar (solo los que cambiaron)
	var name, description, net, kind, model, qosID *string

	if !plan.Name.Equal(state.Name) {
		n := plan.Name.ValueString()
		name = &n
	}

	if !plan.Description.Equal(state.Description) {
		d := plan.Description.ValueString()
		description = &d
	}

	if !plan.Net.Equal(state.Net) {
		n := plan.Net.ValueString()
		net = &n
	}

	if !plan.Kind.Equal(state.Kind) {
		k := plan.Kind.ValueString()
		kind = &k
	}

	if !plan.Model.Equal(state.Model) {
		m := plan.Model.ValueString()
		model = &m
	}

	if !plan.QoSID.Equal(state.QoSID) {
		q := plan.QoSID.ValueString()
		qosID = &q
	}

	// Actualizar la interfaz de red
	err := r.client.UpdateNetworkInterface(
		plan.ID.ValueString(),
		name,
		description,
		net,
		kind,
		model,
		qosID,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error actualizando interfaz de red",
			"No se pudo actualizar la interfaz de red ID "+plan.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Obtener la interfaz actualizada
	iface, err := r.client.GetNetworkInterface(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error leyendo interfaz de red actualizada",
			"No se pudo leer la interfaz actualizada: "+err.Error(),
		)
		return
	}

	// Update state with refreshed values
	plan.Name = types.StringValue(iface.Name)
	plan.Net = types.StringValue(iface.Net)
	if iface.Description != "" {
		plan.Description = types.StringValue(iface.Description)
	}
	if iface.Kind != "" {
		plan.Kind = types.StringValue(iface.Kind)
	}
	if iface.Model != "" {
		plan.Model = types.StringValue(iface.Model)
	}
	if iface.QoSID != "" {
		plan.QoSID = types.StringValue(iface.QoSID)
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *networkInterfaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state networkInterfaceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing interface
	err := r.client.DeleteNetworkInterface(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error eliminando interfaz de red",
			"No se pudo eliminar la interfaz de red ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}
}
