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
	_ resource.Resource              = &networkResource{}
	_ resource.ResourceWithConfigure = &networkResource{}
)

// NewNetworkResource is a helper function to simplify the provider implementation.
func NewNetworkResource() resource.Resource {
	return &networkResource{}
}

// networkResource is the resource implementation.
type networkResource struct {
	client *client.Client
}

// networkResourceModel maps the resource schema data.
type networkResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Model       types.String `tfsdk:"model"`
	QoSID       types.String `tfsdk:"qos_id"`
	MetadataID  types.String `tfsdk:"metadata_id"`
}

// Metadata returns the resource type name.
func (r *networkResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network"
}

// Schema defines the schema for the resource.
func (r *networkResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Gestiona una red virtual de usuario en Isard VDI.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "ID único de la red (generado automáticamente).",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Nombre de la red.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Descripción de la red.",
				Optional:    true,
			},
			"model": schema.StringAttribute{
				Description: "Modelo de red (por defecto: virtio).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"qos_id": schema.StringAttribute{
				Description: "ID de QoS para limitar ancho de banda (por defecto: unlimited).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"metadata_id": schema.StringAttribute{
				Description: "ID de metadata generado para OpenFlow (solo lectura).",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *networkResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *networkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan networkResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Crear la red
	model := plan.Model.ValueString()
	qosID := plan.QoSID.ValueString()

	// allowed por defecto: solo propietario
	allowed := map[string]interface{}{
		"roles":      false,
		"categories": false,
		"groups":     false,
		"users":      []interface{}{},
	}

	networkID, err := r.client.CreateNetwork(
		plan.Name.ValueString(),
		plan.Description.ValueString(),
		model,
		qosID,
		allowed,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creando red",
			"No se pudo crear la red: "+err.Error(),
		)
		return
	}

	// Obtener la red creada para leer los valores computados
	network, err := r.client.GetNetwork(networkID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error leyendo red creada",
			"No se pudo leer la red recién creada: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(network.ID)
	plan.Model = types.StringValue(network.Model)
	plan.QoSID = types.StringValue(network.QoSID)
	plan.MetadataID = types.StringValue(network.MetadataID)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *networkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state networkResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed network value from Isard
	network, err := r.client.GetNetwork(state.ID.ValueString())
	if err != nil {
		if err.Error() == "network not found" {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error leyendo red",
			"No se pudo leer la red ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	state.Name = types.StringValue(network.Name)
	state.Description = types.StringValue(network.Description)
	state.Model = types.StringValue(network.Model)
	state.QoSID = types.StringValue(network.QoSID)
	state.MetadataID = types.StringValue(network.MetadataID)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *networkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan networkResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state
	var state networkResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preparar los valores a actualizar (solo los que cambiaron)
	var name, description, qosID *string

	if !plan.Name.Equal(state.Name) {
		n := plan.Name.ValueString()
		name = &n
	}

	if !plan.Description.Equal(state.Description) {
		d := plan.Description.ValueString()
		description = &d
	}

	if !plan.QoSID.Equal(state.QoSID) {
		q := plan.QoSID.ValueString()
		qosID = &q
	}

	// Actualizar la red
	err := r.client.UpdateNetwork(
		plan.ID.ValueString(),
		name,
		description,
		qosID,
		nil, // allowed no se actualiza por ahora
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error actualizando red",
			"No se pudo actualizar la red ID "+plan.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Obtener la red actualizada
	network, err := r.client.GetNetwork(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error leyendo red actualizada",
			"No se pudo leer la red actualizada: "+err.Error(),
		)
		return
	}

	// Update state with refreshed values
	plan.Name = types.StringValue(network.Name)
	plan.Description = types.StringValue(network.Description)
	plan.Model = types.StringValue(network.Model)
	plan.QoSID = types.StringValue(network.QoSID)
	plan.MetadataID = types.StringValue(network.MetadataID)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *networkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state networkResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing network
	err := r.client.DeleteNetwork(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error eliminando red",
			"No se pudo eliminar la red ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}
}
