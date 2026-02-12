package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tknika/terraform-provider-isardvdi/internal/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &qosNetResource{}
	_ resource.ResourceWithConfigure = &qosNetResource{}
)

// NewQoSNetResource is a helper function to simplify the provider implementation.
func NewQoSNetResource() resource.Resource {
	return &qosNetResource{}
}

// qosNetResource is the resource implementation.
type qosNetResource struct {
	client *client.Client
}

// qosNetResourceModel maps the resource schema data.
type qosNetResourceModel struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	AverageDownload  types.Int64  `tfsdk:"average_download"`
	AverageUpload    types.Int64  `tfsdk:"average_upload"`
	PeakDownload     types.Int64  `tfsdk:"peak_download"`
	PeakUpload       types.Int64  `tfsdk:"peak_upload"`
	BurstDownload    types.Int64  `tfsdk:"burst_download"`
	BurstUpload      types.Int64  `tfsdk:"burst_upload"`
}

// Metadata returns the resource type name.
func (r *qosNetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_qos_net"
}

// Schema defines the schema for the resource.
func (r *qosNetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Gestiona un perfil de QoS de red en Isard VDI.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "ID único del QoS de red.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Nombre del perfil de QoS.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Descripción del perfil de QoS.",
				Optional:    true,
			},
			"average_download": schema.Int64Attribute{
				Description: "Velocidad media de descarga en KB/s.",
				Optional:    true,
			},
			"average_upload": schema.Int64Attribute{
				Description: "Velocidad media de subida en KB/s.",
				Optional:    true,
			},
			"peak_download": schema.Int64Attribute{
				Description: "Velocidad pico de descarga en KB/s.",
				Optional:    true,
			},
			"peak_upload": schema.Int64Attribute{
				Description: "Velocidad pico de subida en KB/s.",
				Optional:    true,
			},
			"burst_download": schema.Int64Attribute{
				Description: "Ráfaga de descarga en KB.",
				Optional:    true,
			},
			"burst_upload": schema.Int64Attribute{
				Description: "Ráfaga de subida en KB.",
				Optional:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *qosNetResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *qosNetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan qosNetResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Construir el objeto bandwidth
	bandwidth := make(map[string]interface{})
	
	if !plan.AverageDownload.IsNull() {
		bandwidth["average_download"] = plan.AverageDownload.ValueInt64()
	}
	if !plan.AverageUpload.IsNull() {
		bandwidth["average_upload"] = plan.AverageUpload.ValueInt64()
	}
	if !plan.PeakDownload.IsNull() {
		bandwidth["peak_download"] = plan.PeakDownload.ValueInt64()
	}
	if !plan.PeakUpload.IsNull() {
		bandwidth["peak_upload"] = plan.PeakUpload.ValueInt64()
	}
	if !plan.BurstDownload.IsNull() {
		bandwidth["burst_download"] = plan.BurstDownload.ValueInt64()
	}
	if !plan.BurstUpload.IsNull() {
		bandwidth["burst_upload"] = plan.BurstUpload.ValueInt64()
	}

	// Crear el QoS de red
	qosID, err := r.client.CreateQoSNet(
		plan.Name.ValueString(),
		plan.Description.ValueString(),
		bandwidth,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creando QoS de red",
			"No se pudo crear el QoS de red: "+err.Error(),
		)
		return
	}

	// Establecer el ID
	plan.ID = types.StringValue(qosID)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *qosNetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state qosNetResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed qos value from Isard
	qos, err := r.client.GetQoSNet(state.ID.ValueString())
	if err != nil {
		if err.Error() == "qos_net not found" {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error leyendo QoS de red",
			"No se pudo leer el QoS de red ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	state.Name = types.StringValue(qos.Name)
	if qos.Description != "" {
		state.Description = types.StringValue(qos.Description)
	}

	// Extraer valores de bandwidth si existen
	if qos.Bandwidth != nil {
		if val, ok := qos.Bandwidth["average_download"].(float64); ok {
			state.AverageDownload = types.Int64Value(int64(val))
		}
		if val, ok := qos.Bandwidth["average_upload"].(float64); ok {
			state.AverageUpload = types.Int64Value(int64(val))
		}
		if val, ok := qos.Bandwidth["peak_download"].(float64); ok {
			state.PeakDownload = types.Int64Value(int64(val))
		}
		if val, ok := qos.Bandwidth["peak_upload"].(float64); ok {
			state.PeakUpload = types.Int64Value(int64(val))
		}
		if val, ok := qos.Bandwidth["burst_download"].(float64); ok {
			state.BurstDownload = types.Int64Value(int64(val))
		}
		if val, ok := qos.Bandwidth["burst_upload"].(float64); ok {
			state.BurstUpload = types.Int64Value(int64(val))
		}
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *qosNetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan qosNetResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state
	var state qosNetResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preparar los valores a actualizar
	var name, description *string

	if !plan.Name.Equal(state.Name) {
		n := plan.Name.ValueString()
		name = &n
	}

	if !plan.Description.Equal(state.Description) {
		d := plan.Description.ValueString()
		description = &d
	}

	// Construir bandwidth solo si algún valor cambió
	var bandwidth map[string]interface{}
	if !plan.AverageDownload.Equal(state.AverageDownload) ||
		!plan.AverageUpload.Equal(state.AverageUpload) ||
		!plan.PeakDownload.Equal(state.PeakDownload) ||
		!plan.PeakUpload.Equal(state.PeakUpload) ||
		!plan.BurstDownload.Equal(state.BurstDownload) ||
		!plan.BurstUpload.Equal(state.BurstUpload) {
		
		bandwidth = make(map[string]interface{})
		if !plan.AverageDownload.IsNull() {
			bandwidth["average_download"] = plan.AverageDownload.ValueInt64()
		}
		if !plan.AverageUpload.IsNull() {
			bandwidth["average_upload"] = plan.AverageUpload.ValueInt64()
		}
		if !plan.PeakDownload.IsNull() {
			bandwidth["peak_download"] = plan.PeakDownload.ValueInt64()
		}
		if !plan.PeakUpload.IsNull() {
			bandwidth["peak_upload"] = plan.PeakUpload.ValueInt64()
		}
		if !plan.BurstDownload.IsNull() {
			bandwidth["burst_download"] = plan.BurstDownload.ValueInt64()
		}
		if !plan.BurstUpload.IsNull() {
			bandwidth["burst_upload"] = plan.BurstUpload.ValueInt64()
		}
	}

	// Actualizar el QoS de red
	err := r.client.UpdateQoSNet(
		plan.ID.ValueString(),
		name,
		description,
		bandwidth,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error actualizando QoS de red",
			"No se pudo actualizar el QoS de red ID "+plan.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Set state with updated values
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *qosNetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state qosNetResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing QoS de red
	err := r.client.DeleteQoSNet(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error eliminando QoS de red",
			"No se pudo eliminar el QoS de red ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}
}
