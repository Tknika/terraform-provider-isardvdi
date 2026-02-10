package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tknika/terraform-provider-isard/internal/client"
)

var _ datasource.DataSource = &mediasDataSource{}

func NewMediasDataSource() datasource.DataSource {
	return &mediasDataSource{}
}

type mediasDataSource struct {
	client *client.Client
}

type mediasDataSourceModel struct {
	ID         types.String  `tfsdk:"id"`
	NameFilter types.String  `tfsdk:"name_filter"`
	Kind       types.String  `tfsdk:"kind"`
	Status     types.String  `tfsdk:"status"`
	CategoryID types.String  `tfsdk:"category_id"`
	GroupID    types.String  `tfsdk:"group_id"`
	UserID     types.String  `tfsdk:"user_id"`
	Medias     []mediaModel  `tfsdk:"medias"`
}

type mediaModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	URL         types.String `tfsdk:"url"`
	Kind        types.String `tfsdk:"kind"`
	Status      types.String `tfsdk:"status"`
	User        types.String `tfsdk:"user"`
	Category    types.String `tfsdk:"category"`
	Group       types.String `tfsdk:"group"`
	Icon        types.String `tfsdk:"icon"`
	Path        types.String `tfsdk:"path"`
}

func (d *mediasDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_medias"
}

func (d *mediasDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Obtiene la lista de medios (ISOs, discos) de Isard VDI. Permite filtrar por nombre, tipo, estado, categoría, grupo y usuario.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identificador del data source.",
				Computed:    true,
			},
			"name_filter": schema.StringAttribute{
				Description: "Filtro opcional para buscar medios por nombre (búsqueda insensible a mayúsculas, coincidencia parcial).",
				Optional:    true,
			},
			"kind": schema.StringAttribute{
				Description: "Filtro opcional por tipo de medio (iso, disk, floppy).",
				Optional:    true,
			},
			"status": schema.StringAttribute{
				Description: "Filtro opcional por estado del medio (DownloadStarting, Downloading, Downloaded, Failed, etc.).",
				Optional:    true,
			},
			"category_id": schema.StringAttribute{
				Description: "Filtro opcional por ID de categoría.",
				Optional:    true,
			},
			"group_id": schema.StringAttribute{
				Description: "Filtro opcional por ID de grupo.",
				Optional:    true,
			},
			"user_id": schema.StringAttribute{
				Description: "Filtro opcional por ID de usuario propietario.",
				Optional:    true,
			},
			"medias": schema.ListNestedAttribute{
				Description: "Lista de medios que coinciden con los filtros.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "ID del medio.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Nombre del medio.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "Descripción del medio.",
							Computed:    true,
						},
						"url": schema.StringAttribute{
							Description: "URL de descarga del medio.",
							Computed:    true,
						},
						"kind": schema.StringAttribute{
							Description: "Tipo de medio (iso, disk, floppy).",
							Computed:    true,
						},
						"status": schema.StringAttribute{
							Description: "Estado del medio.",
							Computed:    true,
						},
						"user": schema.StringAttribute{
							Description: "ID del usuario propietario.",
							Computed:    true,
						},
						"category": schema.StringAttribute{
							Description: "ID de la categoría.",
							Computed:    true,
						},
						"group": schema.StringAttribute{
							Description: "ID del grupo.",
							Computed:    true,
						},
						"icon": schema.StringAttribute{
							Description: "Icono del medio.",
							Computed:    true,
						},
						"path": schema.StringAttribute{
							Description: "Ruta del archivo en el servidor.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *mediasDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*client.Client)
}

func (d *mediasDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data mediasDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Obtener todos los medios
	medias, err := d.client.GetMedias()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error obteniendo medios",
			fmt.Sprintf("No se pudo obtener la lista de medios: %s", err.Error()),
		)
		return
	}

	// Aplicar filtros
	var filteredMedias []client.Media
	for _, media := range medias {
		// Filtrar por nombre (case-insensitive substring match)
		if !data.NameFilter.IsNull() && !data.NameFilter.IsUnknown() {
			if !containsIgnoreCaseMedias(media.Name, data.NameFilter.ValueString()) {
				continue
			}
		}

		// Filtrar por tipo
		if !data.Kind.IsNull() && !data.Kind.IsUnknown() {
			if media.Kind != data.Kind.ValueString() {
				continue
			}
		}

		// Filtrar por estado
		if !data.Status.IsNull() && !data.Status.IsUnknown() {
			if media.Status != data.Status.ValueString() {
				continue
			}
		}

		// Filtrar por categoría
		if !data.CategoryID.IsNull() && !data.CategoryID.IsUnknown() {
			if media.Category != data.CategoryID.ValueString() {
				continue
			}
		}

		// Filtrar por grupo
		if !data.GroupID.IsNull() && !data.GroupID.IsUnknown() {
			if media.Group != data.GroupID.ValueString() {
				continue
			}
		}

		// Filtrar por usuario
		if !data.UserID.IsNull() && !data.UserID.IsUnknown() {
			if media.User != data.UserID.ValueString() {
				continue
			}
		}

		filteredMedias = append(filteredMedias, media)
	}

	// Convertir a modelo de Terraform
	data.Medias = make([]mediaModel, 0, len(filteredMedias))
	for _, media := range filteredMedias {
		mediaItem := mediaModel{
			ID:          types.StringValue(media.ID),
			Name:        types.StringValue(media.Name),
			Description: types.StringValue(media.Description),
			URL:         types.StringValue(media.URL),
			Kind:        types.StringValue(media.Kind),
			Status:      types.StringValue(media.Status),
			User:        types.StringValue(media.User),
			Category:    types.StringValue(media.Category),
			Group:       types.StringValue(media.Group),
			Icon:        types.StringValue(media.Icon),
			Path:        types.StringValue(media.Path),
		}
		data.Medias = append(data.Medias, mediaItem)
	}

	// Establecer un ID para el data source
	data.ID = types.StringValue("medias")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// containsIgnoreCaseMedias verifica si str contiene substr (case-insensitive)
func containsIgnoreCaseMedias(str, substr string) bool {
	return strings.Contains(strings.ToLower(str), strings.ToLower(substr))
}
