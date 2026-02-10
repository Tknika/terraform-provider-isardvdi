package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tknika/terraform-provider-isard/internal/client"
)

var (
	_ resource.Resource                = &mediaResource{}
	_ resource.ResourceWithConfigure   = &mediaResource{}
	_ resource.ResourceWithImportState = &mediaResource{}
)

func NewMediaResource() resource.Resource {
	return &mediaResource{}
}

type mediaResource struct {
	client *client.Client
}

type mediaResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	URL         types.String `tfsdk:"url"`
	Kind        types.String `tfsdk:"kind"`
	Allowed     types.Object `tfsdk:"allowed"`
}

func (r *mediaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_media"
}

func (r *mediaResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Gestiona un media (ISO, disk image, etc.) en Isard VDI. Los medias son archivos que se descargan de URLs especificadas.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Identificador único del media",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Nombre del media (máximo 50 caracteres)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Descripción del media (máximo 255 caracteres)",
			},
			"url": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "URL HTTPS del archivo a descargar. Debe ser una URL válida y accesible",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"kind": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Tipo de media. Valores comunes: 'iso' (imagen ISO), 'disk' (imagen de disco), 'floppy' (imagen de disquete)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"allowed": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Configuración de usuarios, grupos y categorías permitidos para usar este media",
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
		},
	}
}

func (r *mediaResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *mediaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan mediaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Construir el mapa allowed si se especifica
	var allowed map[string]interface{}
	if !plan.Allowed.IsNull() {
		allowedAttrs := plan.Allowed.Attributes()
		allowed = make(map[string]interface{})

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
	}

	// Crear el media usando la API
	mediaID, err := r.client.CreateMedia(
		plan.Name.ValueString(),
		plan.Description.ValueString(),
		plan.URL.ValueString(),
		plan.Kind.ValueString(),
		allowed,
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creando el media",
			fmt.Sprintf("No se pudo crear el media: %s", err.Error()),
		)
		return
	}

	plan.ID = types.StringValue(mediaID)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *mediaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state mediaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	media, err := r.client.GetMedia(state.ID.ValueString())
	if err != nil {
		// Si el media no se encuentra, eliminarlo del state
		resp.State.RemoveResource(ctx)
		return
	}

	// Si el media está en estado "deleted", eliminarlo del state
	if media.Status == "deleted" {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(media.Name)
	state.Description = types.StringValue(media.Description)
	// No actualizamos URL y Kind porque son inmutables

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *mediaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan mediaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Los medias en Isard VDI no soportan actualización de los campos principales
	// Solo se puede recrear, por lo que name, url y kind tienen RequiresReplace
	// La descripción tampoco se puede actualizar directamente vía API

	resp.Diagnostics.AddWarning(
		"Actualización limitada",
		"Los medias no soportan actualización de sus propiedades. Para modificar un media, debe eliminarse y recrearse.",
	)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *mediaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state mediaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Primero verificar el estado actual del media
	media, err := r.client.GetMedia(state.ID.ValueString())
	if err != nil {
		// Si el media no existe, no hay nada que eliminar
		resp.Diagnostics.AddWarning(
			"Media no encontrado",
			fmt.Sprintf("El media (ID: %s) no se encontró, puede haber sido eliminado manualmente", state.ID.ValueString()),
		)
		return
	}

	// Si ya está en estado deleted, no intentar eliminarlo
	if media.Status == "deleted" {
		resp.Diagnostics.AddWarning(
			"Media ya eliminado",
			fmt.Sprintf("El media (ID: %s) ya está en estado 'deleted'", state.ID.ValueString()),
		)
		return
	}

	err = r.client.DeleteMedia(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error eliminando el media",
			fmt.Sprintf("No se pudo eliminar el media (ID: %s): %s", state.ID.ValueString(), err.Error()),
		)
		return
	}
}

func (r *mediaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
