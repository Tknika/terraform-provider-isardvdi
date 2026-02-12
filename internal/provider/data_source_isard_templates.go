package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tknika/terraform-provider-isardvdi/internal/client"
)

var _ datasource.DataSource = &templatesDataSource{}

func NewTemplatesDataSource() datasource.DataSource {
	return &templatesDataSource{}
}

type templatesDataSource struct {
	client *client.Client
}

type templatesDataSourceModel struct {
	ID          types.String   `tfsdk:"id"`
	NameFilter  types.String   `tfsdk:"name_filter"`
	Templates   []templateModel `tfsdk:"templates"`
}

type templateModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Category    types.String `tfsdk:"category"`
	Group       types.String `tfsdk:"group"`
	UserID      types.String `tfsdk:"user_id"`
	Icon        types.String `tfsdk:"icon"`
	Description types.String `tfsdk:"description"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	Status      types.String `tfsdk:"status"`
	DesktopSize types.Int64  `tfsdk:"desktop_size"`
}

func (d *templatesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_templates"
}

func (d *templatesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of available templates from Isard VDI.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Placeholder identifier for the data source.",
				Computed:    true,
			},
			"name_filter": schema.StringAttribute{
				Description: "Optional filter to match template names (case-insensitive substring match).",
				Optional:    true,
			},
			"templates": schema.ListNestedAttribute{
				Description: "List of templates available to the user.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Template ID.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Template name.",
							Computed:    true,
						},
						"category": schema.StringAttribute{
							Description: "Category ID.",
							Computed:    true,
						},
						"group": schema.StringAttribute{
							Description: "Group ID.",
							Computed:    true,
						},
						"user_id": schema.StringAttribute{
							Description: "User ID who owns the template.",
							Computed:    true,
						},
						"icon": schema.StringAttribute{
							Description: "Icon name.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "Template description.",
							Computed:    true,
						},
						"enabled": schema.BoolAttribute{
							Description: "Whether the template is enabled.",
							Computed:    true,
						},
						"status": schema.StringAttribute{
							Description: "Template status.",
							Computed:    true,
						},
						"desktop_size": schema.Int64Attribute{
							Description: "Desktop size in bytes.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *templatesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *templatesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data templatesDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	templates, err := d.client.GetTemplates()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read templates, got error: %s", err))
		return
	}

	// Apply name filter if provided
	nameFilter := data.NameFilter.ValueString()
	var filteredTemplates []client.Template
	for _, template := range templates {
		if nameFilter == "" || containsIgnoreCase(template.Name, nameFilter) {
			filteredTemplates = append(filteredTemplates, template)
		}
	}

	// Map filtered templates to model
	data.Templates = make([]templateModel, len(filteredTemplates))
	for i, template := range filteredTemplates {
		data.Templates[i] = templateModel{
			ID:          types.StringValue(template.ID),
			Name:        types.StringValue(template.Name),
			Category:    types.StringValue(template.Category),
			Group:       types.StringValue(template.Group),
			UserID:      types.StringValue(template.UserID),
			Icon:        types.StringValue(template.Icon),
			Description: types.StringValue(template.Description),
			Enabled:     types.BoolValue(template.Enabled),
			Status:      types.StringValue(template.Status),
			DesktopSize: types.Int64Value(template.DesktopSize),
		}
	}

	data.ID = types.StringValue("templates")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func containsIgnoreCase(s, substr string) bool {
	s = toLower(s)
	substr = toLower(substr)
	return contains(s, substr)
}

func toLower(s string) string {
	result := make([]rune, len(s))
	for i, r := range s {
		if r >= 'A' && r <= 'Z' {
			result[i] = r + 32
		} else {
			result[i] = r
		}
	}
	return string(result)
}

func contains(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(substr) > len(s) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
