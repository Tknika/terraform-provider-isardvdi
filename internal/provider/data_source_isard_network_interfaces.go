package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tknika/terraform-provider-isard/internal/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &networkInterfacesDataSource{}
	_ datasource.DataSourceWithConfigure = &networkInterfacesDataSource{}
)

// NewNetworkInterfacesDataSource is a helper function to simplify the provider implementation.
func NewNetworkInterfacesDataSource() datasource.DataSource {
	return &networkInterfacesDataSource{}
}

// networkInterfacesDataSource is the data source implementation.
type networkInterfacesDataSource struct {
	client *client.Client
}

// networkInterfacesDataSourceModel maps the data source schema data.
type networkInterfacesDataSourceModel struct {
	ID     types.String                  `tfsdk:"id"`
	Name   types.String                  `tfsdk:"name"`
	Filter *networkInterfaceFilterModel  `tfsdk:"filter"`
	Items  []networkInterfaceDetailModel `tfsdk:"interfaces"`
}

// networkInterfaceFilterModel maps the filter schema.
type networkInterfaceFilterModel struct {
	Name types.String `tfsdk:"name"`
	Kind types.String `tfsdk:"kind"`
	Net  types.String `tfsdk:"net"`
}

// networkInterfaceDetailModel maps individual interface details.
type networkInterfaceDetailModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Net         types.String `tfsdk:"net"`
	Kind        types.String `tfsdk:"kind"`
	Model       types.String `tfsdk:"model"`
	QoSID       types.String `tfsdk:"qos_id"`
}

// Metadata returns the data source type name.
func (d *networkInterfacesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_interfaces"
}

// Schema defines the schema for the data source.
func (d *networkInterfacesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Obtiene una lista de interfaces de red del sistema. Permite filtrar por nombre, tipo o red.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identificador del data source.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Nombre exacto de la interfaz a buscar (opcional, deprecated, usar filter en su lugar).",
				Optional:    true,
			},
			"filter": schema.SingleNestedAttribute{
				Description: "Filtros para buscar interfaces de red.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Description: "Nombre de la interfaz (búsqueda parcial, case-insensitive).",
						Optional:    true,
					},
					"kind": schema.StringAttribute{
						Description: "Tipo de interfaz (bridge, network, ovs, personal).",
						Optional:    true,
					},
					"net": schema.StringAttribute{
						Description: "Red/bridge del sistema asociada.",
						Optional:    true,
					},
				},
			},
			"interfaces": schema.ListNestedAttribute{
				Description: "Lista de interfaces de red encontradas.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "ID único de la interfaz.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Nombre de la interfaz.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "Descripción de la interfaz.",
							Computed:    true,
						},
						"net": schema.StringAttribute{
							Description: "Red/bridge del sistema.",
							Computed:    true,
						},
						"kind": schema.StringAttribute{
							Description: "Tipo de interfaz.",
							Computed:    true,
						},
						"model": schema.StringAttribute{
							Description: "Modelo de la interfaz.",
							Computed:    true,
						},
						"qos_id": schema.StringAttribute{
							Description: "ID del perfil QoS.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *networkInterfacesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state networkInterfacesDataSourceModel

	// Get configuration
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Obtener todas las interfaces
	interfaces, err := d.client.ListNetworkInterfaces()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error obteniendo interfaces de red",
			"No se pudo obtener la lista de interfaces: "+err.Error(),
		)
		return
	}

	// Aplicar filtros
	var filteredInterfaces []client.NetworkInterface
	
	// Si hay un filtro por nombre (deprecated)
	if !state.Name.IsNull() && state.Name.ValueString() != "" {
		searchName := state.Name.ValueString()
		for _, iface := range interfaces {
			if iface.Name == searchName {
				filteredInterfaces = append(filteredInterfaces, iface)
			}
		}
	} else if state.Filter != nil {
		// Usar filtros nuevos
		filteredInterfaces = d.applyFilters(interfaces, state.Filter)
	} else {
		// Sin filtros, devolver todas
		filteredInterfaces = interfaces
	}

	// Mapear a la estructura del state
	state.Items = []networkInterfaceDetailModel{}
	for _, iface := range filteredInterfaces {
		item := networkInterfaceDetailModel{
			ID:          types.StringValue(iface.ID),
			Name:        types.StringValue(iface.Name),
			Description: types.StringValue(iface.Description),
			Net:         types.StringValue(iface.Net),
			Kind:        types.StringValue(iface.Kind),
			Model:       types.StringValue(iface.Model),
			QoSID:       types.StringValue(iface.QoSID),
		}
		state.Items = append(state.Items, item)
	}

	// Set ID del data source
	if !state.Name.IsNull() && state.Name.ValueString() != "" {
		state.ID = types.StringValue(fmt.Sprintf("network-interface-%s", state.Name.ValueString()))
	} else {
		state.ID = types.StringValue("network-interfaces-all")
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// applyFilters aplica los filtros especificados a la lista de interfaces
func (d *networkInterfacesDataSource) applyFilters(interfaces []client.NetworkInterface, filter *networkInterfaceFilterModel) []client.NetworkInterface {
	var result []client.NetworkInterface

	for _, iface := range interfaces {
		match := true

		// Filtro por nombre (búsqueda parcial, case-insensitive)
		if !filter.Name.IsNull() && filter.Name.ValueString() != "" {
			filterName := filter.Name.ValueString()
			if !contains(iface.Name, filterName) {
				match = false
			}
		}

		// Filtro por tipo (exacto)
		if !filter.Kind.IsNull() && filter.Kind.ValueString() != "" {
			if iface.Kind != filter.Kind.ValueString() {
				match = false
			}
		}

		// Filtro por red (exacto)
		if !filter.Net.IsNull() && filter.Net.ValueString() != "" {
			if iface.Net != filter.Net.ValueString() {
				match = false
			}
		}

		if match {
			result = append(result, iface)
		}
	}

	return result
}

// Configure adds the provider configured client to the data source.
func (d *networkInterfacesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}
