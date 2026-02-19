package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tknika/terraform-provider-isardvdi/internal/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &deploymentResource{}
	_ resource.ResourceWithConfigure = &deploymentResource{}
)

// NewDeploymentResource is a helper function to simplify the provider implementation.
func NewDeploymentResource() resource.Resource {
	return &deploymentResource{}
}

// deploymentResource is the resource implementation.
type deploymentResource struct {
	client *client.Client
}

// deploymentResourceModel maps the resource schema data.
type deploymentResourceModel struct {
	ID                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	Description        types.String `tfsdk:"description"`
	TemplateID         types.String `tfsdk:"template_id"`
	DesktopName        types.String `tfsdk:"desktop_name"`
	Visible            types.Bool   `tfsdk:"visible"`
	Allowed            types.Object `tfsdk:"allowed"`
	VCPUs              types.Int64  `tfsdk:"vcpus"`
	Memory             types.Float64 `tfsdk:"memory"`
	NetworkInterfaces  types.List   `tfsdk:"network_interfaces"`
	ISOs               types.List   `tfsdk:"isos"`
	Floppies           types.List   `tfsdk:"floppies"`
	UserPermissions    types.List   `tfsdk:"user_permissions"`
	Viewers            types.List   `tfsdk:"viewers"`
	ForceStopOnDestroy types.Bool   `tfsdk:"force_stop_on_destroy"`
}

// Metadata returns the resource type name.
func (r *deploymentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_deployment"
}

// Schema defines the schema for the resource.
func (r *deploymentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Gestiona un deployment en Isard VDI. Los deployments permiten crear múltiples desktops a partir de una plantilla para diferentes usuarios.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Identificador único del deployment",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Nombre del deployment (mínimo 4 caracteres, máximo 50)",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Descripción del deployment (máximo 255 caracteres)",
			},
			"template_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "ID de la plantilla a utilizar para crear los desktops del deployment",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"desktop_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Nombre base para los desktops creados en el deployment",
			},
			"visible": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Si los desktops del deployment son visibles para los usuarios (por defecto: false)",
			},
			"allowed": schema.SingleNestedAttribute{
				Required:            true,
				MarkdownDescription: "Configuración de usuarios, grupos y categorías permitidos para acceder a este deployment",
				Attributes: map[string]schema.Attribute{
					"roles": schema.ListAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						MarkdownDescription: "Lista de roles permitidos",
					},
					"categories": schema.ListAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						MarkdownDescription: "Lista de IDs de categorías permitidas",
					},
					"groups": schema.ListAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						MarkdownDescription: "Lista de IDs de grupos permitidos",
					},
					"users": schema.ListAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						MarkdownDescription: "Lista de IDs de usuarios permitidos",
					},
				},
			},
			"vcpus": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(2),
				MarkdownDescription: "Número de CPUs virtuales para los desktops (por defecto: 2)",
			},
			"memory": schema.Float64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             float64default.StaticFloat64(2.0),
				MarkdownDescription: "Memoria RAM en GB para los desktops (por defecto: 2.0 GB)",
			},
			"network_interfaces": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{
					types.StringValue("default"),
					types.StringValue("wireguard"),
				})),
				MarkdownDescription: "Lista de IDs de interfaces de red a utilizar (por defecto: ['default', 'wireguard'])",
			},
			"isos": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "Lista de IDs de medios ISO a adjuntar a los desktops del deployment",
			},
			"floppies": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "Lista de IDs de medios floppy a adjuntar a los desktops del deployment",
			},
			"user_permissions": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "Lista de permisos de usuario para el deployment",
			},
			"viewers": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "Lista de viewers habilitados (ej: ['browser_vnc', 'file_spice', 'file_rdpgw', 'browser_rdp']). Si no se especifica, se usan los del template.",
			},
			"force_stop_on_destroy": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Si es true, detiene todas las máquinas virtuales del deployment antes de eliminarlo (por defecto: false)",
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *deploymentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates a new resource.
func (r *deploymentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan deploymentResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extraer el objeto allowed
	allowedAttrs := plan.Allowed.Attributes()

	// Construir el mapa allowed para la API
	// La API espera false para campos no utilizados, no omitirlos
	allowed := make(map[string]interface{})
	
	if rolesAttr, ok := allowedAttrs["roles"]; ok {
		if roles, ok := rolesAttr.(types.List); ok && !roles.IsNull() {
			var rolesList []string
			diags := roles.ElementsAs(ctx, &rolesList, false)
			resp.Diagnostics.Append(diags...)
			if !resp.Diagnostics.HasError() {
				if len(rolesList) > 0 {
					allowed["roles"] = rolesList
				} else {
					allowed["roles"] = false
				}
			}
		}
	}
	
	if categoriesAttr, ok := allowedAttrs["categories"]; ok {
		if categories, ok := categoriesAttr.(types.List); ok && !categories.IsNull() {
			var categoriesList []string
			diags := categories.ElementsAs(ctx, &categoriesList, false)
			resp.Diagnostics.Append(diags...)
			if !resp.Diagnostics.HasError() {
				if len(categoriesList) > 0 {
					allowed["categories"] = categoriesList
				} else {
					allowed["categories"] = false
				}
			}
		}
	}
	
	if groupsAttr, ok := allowedAttrs["groups"]; ok {
		if groups, ok := groupsAttr.(types.List); ok && !groups.IsNull() {
			var groupsList []string
			diags := groups.ElementsAs(ctx, &groupsList, false)
			resp.Diagnostics.Append(diags...)
			if !resp.Diagnostics.HasError() {
				if len(groupsList) > 0 {
					allowed["groups"] = groupsList
				} else {
					allowed["groups"] = false
				}
			}
		}
	}
	
	if usersAttr, ok := allowedAttrs["users"]; ok {
		if users, ok := usersAttr.(types.List); ok && !users.IsNull() {
			var usersList []string
			diags := users.ElementsAs(ctx, &usersList, false)
			resp.Diagnostics.Append(diags...)
			if !resp.Diagnostics.HasError() {
				if len(usersList) > 0 {
					allowed["users"] = usersList
				} else {
					allowed["users"] = false
				}
			}
		}
	}

	// Preparar hardware personalizado si se especifica
	var vcpus *int64
	var memory *float64
	var interfaces []string
	var isos []string
	var floppies []string
	var userPermissions []string
	var viewers []string
	
	if !plan.VCPUs.IsNull() && !plan.VCPUs.IsUnknown() {
		v := plan.VCPUs.ValueInt64()
		vcpus = &v
	}
	
	if !plan.Memory.IsNull() && !plan.Memory.IsUnknown() {
		m := plan.Memory.ValueFloat64()
		memory = &m
	}
	
	if !plan.NetworkInterfaces.IsNull() && !plan.NetworkInterfaces.IsUnknown() {
		diags := plan.NetworkInterfaces.ElementsAs(ctx, &interfaces, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	
	if !plan.ISOs.IsNull() && !plan.ISOs.IsUnknown() {
		diags := plan.ISOs.ElementsAs(ctx, &isos, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	
	if !plan.Floppies.IsNull() && !plan.Floppies.IsUnknown() {
		diags := plan.Floppies.ElementsAs(ctx, &floppies, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	
	if !plan.UserPermissions.IsNull() && !plan.UserPermissions.IsUnknown() {
		diags := plan.UserPermissions.ElementsAs(ctx, &userPermissions, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	
	if !plan.Viewers.IsNull() && !plan.Viewers.IsUnknown() {
		diags := plan.Viewers.ElementsAs(ctx, &viewers, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Construir guest_properties si se especifican viewers
	var guestProperties map[string]interface{}
	if len(viewers) > 0 {
		viewersMap := make(map[string]interface{})
		for _, viewer := range viewers {
			viewersMap[viewer] = map[string]interface{}{"options": nil}
		}
		guestProperties = map[string]interface{}{
			"viewers": viewersMap,
		}
	}

	// Crear el deployment usando la API
	deploymentID, err := r.client.CreateDeployment(
		plan.Name.ValueString(),
		plan.Description.ValueString(),
		plan.TemplateID.ValueString(),
		plan.DesktopName.ValueString(),
		plan.Visible.ValueBool(),
		allowed,
		vcpus,
		memory,
		interfaces,
		guestProperties,
		nil, // image
		userPermissions,
		isos,
		floppies,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creando el deployment",
			fmt.Sprintf("No se pudo crear el deployment: %s", err.Error()),
		)
		return
	}

	// Actualizar el plan con el ID devuelto por la API
	plan.ID = types.StringValue(deploymentID)

	// Obtener los valores finales del deployment creado
	deployment, err := r.client.GetDeployment(deploymentID)
	if err == nil {
		if deployment.Description != "" {
			plan.Description = types.StringValue(deployment.Description)
		}
		plan.Visible = types.BoolValue(deployment.Visible)
	}

	// Los valores de hardware (vcpus, memory, network_interfaces) ya están en el plan
	// por los defaults del schema o por los valores especificados por el usuario
	// No necesitamos leerlos de la API en Create ya que los enviamos nosotros

	// Escribir el estado
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *deploymentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state deploymentResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Obtener el deployment de la API
	deployment, err := r.client.GetDeployment(state.ID.ValueString())
	if err != nil {
		// Si el deployment no existe (404), eliminarlo del estado
		if err.Error() == "deployment not found" {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error leyendo el deployment",
			fmt.Sprintf("No se pudo leer el deployment (ID: %s): %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	// Actualizar el estado con los valores de la API
	state.Name = types.StringValue(deployment.Name)
	state.Description = types.StringValue(deployment.Description)
	state.DesktopName = types.StringValue(deployment.DesktopName)
	state.TemplateID = types.StringValue(deployment.TemplateID)
	state.Visible = types.BoolValue(deployment.Visible)

	// Actualizar allowed
	if deployment.Allowed != nil {
		allowedBytes, err := json.Marshal(deployment.Allowed)
		if err == nil {
			var allowedMap map[string]interface{}
			if err := json.Unmarshal(allowedBytes, &allowedMap); err == nil {
				// Construir el objeto allowed para terraform
				allowedAttrs := make(map[string]attr.Value)
				
				if roles, ok := allowedMap["roles"].([]interface{}); ok {
					rolesList := make([]attr.Value, len(roles))
					for i, role := range roles {
						rolesList[i] = types.StringValue(fmt.Sprintf("%v", role))
					}
					allowedAttrs["roles"] = types.ListValueMust(types.StringType, rolesList)
				} else {
					allowedAttrs["roles"] = types.ListNull(types.StringType)
				}
				
				if categories, ok := allowedMap["categories"].([]interface{}); ok {
					categoriesList := make([]attr.Value, len(categories))
					for i, cat := range categories {
						categoriesList[i] = types.StringValue(fmt.Sprintf("%v", cat))
					}
					allowedAttrs["categories"] = types.ListValueMust(types.StringType, categoriesList)
				} else {
					allowedAttrs["categories"] = types.ListNull(types.StringType)
				}
				
				if groups, ok := allowedMap["groups"].([]interface{}); ok {
					groupsList := make([]attr.Value, len(groups))
					for i, group := range groups {
						groupsList[i] = types.StringValue(fmt.Sprintf("%v", group))
					}
					allowedAttrs["groups"] = types.ListValueMust(types.StringType, groupsList)
				} else {
					allowedAttrs["groups"] = types.ListNull(types.StringType)
				}
				
				if users, ok := allowedMap["users"].([]interface{}); ok {
					usersList := make([]attr.Value, len(users))
					for i, user := range users {
						usersList[i] = types.StringValue(fmt.Sprintf("%v", user))
					}
					allowedAttrs["users"] = types.ListValueMust(types.StringType, usersList)
				} else {
					allowedAttrs["users"] = types.ListNull(types.StringType)
				}
				
				allowedType := types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"roles":      types.ListType{ElemType: types.StringType},
						"categories": types.ListType{ElemType: types.StringType},
						"groups":     types.ListType{ElemType: types.StringType},
						"users":      types.ListType{ElemType: types.StringType},
					},
				}
				state.Allowed = types.ObjectValueMust(allowedType.AttrTypes, allowedAttrs)
			}
		}
	}

	// Nota: La API devuelve null para hardware en deployments
	// Los valores de vcpus, memory y network_interfaces se mantienen del state
	// ya que son los que se enviaron en la creación

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *deploymentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan deploymentResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Construir los datos de actualización
	updateData := make(map[string]interface{})
	updateData["name"] = plan.Name.ValueString()
	updateData["description"] = plan.Description.ValueString()
	updateData["desktop_name"] = plan.DesktopName.ValueString()

	// Extraer el objeto allowed
	allowedAttrs := plan.Allowed.Attributes()

	// Construir el mapa allowed para la API
	allowed := make(map[string]interface{})
	
	if rolesAttr, ok := allowedAttrs["roles"]; ok {
		if roles, ok := rolesAttr.(types.List); ok && !roles.IsNull() {
			var rolesList []string
			diags := roles.ElementsAs(ctx, &rolesList, false)
			resp.Diagnostics.Append(diags...)
			if !resp.Diagnostics.HasError() && len(rolesList) > 0 {
				allowed["roles"] = rolesList
			}
		}
	}
	
	if categoriesAttr, ok := allowedAttrs["categories"]; ok {
		if categories, ok := categoriesAttr.(types.List); ok && !categories.IsNull() {
			var categoriesList []string
			diags := categories.ElementsAs(ctx, &categoriesList, false)
			resp.Diagnostics.Append(diags...)
			if !resp.Diagnostics.HasError() && len(categoriesList) > 0 {
				allowed["categories"] = categoriesList
			}
		}
	}
	
	if groupsAttr, ok := allowedAttrs["groups"]; ok {
		if groups, ok := groupsAttr.(types.List); ok && !groups.IsNull() {
			var groupsList []string
			diags := groups.ElementsAs(ctx, &groupsList, false)
			resp.Diagnostics.Append(diags...)
			if !resp.Diagnostics.HasError() && len(groupsList) > 0 {
				allowed["groups"] = groupsList
			}
		}
	}
	
	if usersAttr, ok := allowedAttrs["users"]; ok {
		if users, ok := usersAttr.(types.List); ok && !users.IsNull() {
			var usersList []string
			diags := users.ElementsAs(ctx, &usersList, false)
			resp.Diagnostics.Append(diags...)
			if !resp.Diagnostics.HasError() && len(usersList) > 0 {
				allowed["users"] = usersList
			}
		}
	}

	updateData["allowed"] = allowed

	// Actualizar guest_properties si se especifican viewers
	if !plan.Viewers.IsNull() && !plan.Viewers.IsUnknown() {
		var viewers []string
		diags := plan.Viewers.ElementsAs(ctx, &viewers, false)
		resp.Diagnostics.Append(diags...)
		if !resp.Diagnostics.HasError() && len(viewers) > 0 {
			viewersMap := make(map[string]interface{})
			for _, viewer := range viewers {
				viewersMap[viewer] = map[string]interface{}{"options": nil}
			}
			updateData["guest_properties"] = map[string]interface{}{
				"viewers": viewersMap,
			}
		}
	}

	// Actualizar hardware si se especifica
	if !plan.VCPUs.IsNull() || !plan.Memory.IsNull() || !plan.NetworkInterfaces.IsNull() {
		hardware := make(map[string]interface{})
		
		if !plan.VCPUs.IsNull() && !plan.VCPUs.IsUnknown() {
			hardware["vcpus"] = plan.VCPUs.ValueInt64()
		}
		
		if !plan.Memory.IsNull() && !plan.Memory.IsUnknown() {
			hardware["memory"] = plan.Memory.ValueFloat64()
		}
		
		if !plan.NetworkInterfaces.IsNull() && !plan.NetworkInterfaces.IsUnknown() {
			var interfaces []string
			diags := plan.NetworkInterfaces.ElementsAs(ctx, &interfaces, false)
			resp.Diagnostics.Append(diags...)
			if !resp.Diagnostics.HasError() && len(interfaces) > 0 {
				interfacesList := make([]map[string]interface{}, len(interfaces))
				for i, iface := range interfaces {
					interfacesList[i] = map[string]interface{}{"id": iface}
				}
				hardware["interfaces"] = interfacesList
			}
		}
		
		updateData["hardware"] = hardware
	}

	// Actualizar el deployment usando la API
	err := r.client.UpdateDeployment(plan.ID.ValueString(), updateData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error actualizando el deployment",
			fmt.Sprintf("No se pudo actualizar el deployment (ID: %s): %s", plan.ID.ValueString(), err.Error()),
		)
		return
	}

	// Escribir el estado
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *deploymentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state deploymentResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Si force_stop_on_destroy es true, detener todas las VMs del deployment primero
	if !state.ForceStopOnDestroy.IsNull() && state.ForceStopOnDestroy.ValueBool() {
		err := r.client.StopDeployment(state.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddWarning(
				"Advertencia al detener desktops del deployment",
				fmt.Sprintf("No se pudieron detener todos los desktops del deployment (ID: %s): %s. Se procederá con la eliminación.", state.ID.ValueString(), err.Error()),
			)
		} else {
			// Esperar a que las VMs se detengan completamente (máximo 120 segundos)
			err = r.client.WaitForDeploymentStopped(state.ID.ValueString(), 120)
			if err != nil {
				resp.Diagnostics.AddWarning(
					"Advertencia al esperar el stop de desktops",
					fmt.Sprintf("Timeout o error esperando a que se detengan las VMs del deployment (ID: %s): %s. Se procederá con la eliminación.", state.ID.ValueString(), err.Error()),
				)
			}
		}
	}

	// Eliminar el deployment usando la API (permanent=true)
	err := r.client.DeleteDeployment(state.ID.ValueString(), true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error eliminando el deployment",
			fmt.Sprintf("No se pudo eliminar el deployment (ID: %s): %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	// El estado se elimina automáticamente si la función termina sin errores
}
