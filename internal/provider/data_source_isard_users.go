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

var _ datasource.DataSource = &usersDataSource{}

func NewUsersDataSource() datasource.DataSource {
	return &usersDataSource{}
}

type usersDataSource struct {
	client *client.Client
}

type usersDataSourceModel struct {
	ID           types.String   `tfsdk:"id"`
	NameFilter   types.String   `tfsdk:"name_filter"`
	CategoryID   types.String   `tfsdk:"category_id"`
	GroupID      types.String   `tfsdk:"group_id"`
	Role         types.String   `tfsdk:"role"`
	Active       types.Bool     `tfsdk:"active"`
	Users        []userModel    `tfsdk:"users"`
}

type userModel struct {
	ID                     types.String   `tfsdk:"id"`
	Name                   types.String   `tfsdk:"name"`
	Username               types.String   `tfsdk:"username"`
	UID                    types.String   `tfsdk:"uid"`
	Email                  types.String   `tfsdk:"email"`
	Active                 types.Bool     `tfsdk:"active"`
	Role                   types.String   `tfsdk:"role"`
	Category               types.String   `tfsdk:"category"`
	Group                  types.String   `tfsdk:"group"`
	SecondaryGroups        []types.String `tfsdk:"secondary_groups"`
	Provider               types.String   `tfsdk:"provider"`
	EmailVerified          types.Bool     `tfsdk:"email_verified"`
	DisclaimerAcknowledged types.Bool     `tfsdk:"disclaimer_acknowledged"`
	RoleName               types.String   `tfsdk:"role_name"`
	CategoryName           types.String   `tfsdk:"category_name"`
	GroupName              types.String   `tfsdk:"group_name"`
}

func (d *usersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_users"
}

func (d *usersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of users from Isard VDI. Allows filtering by name, category, group, role and active status.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Placeholder identifier for the data source.",
				Computed:    true,
			},
			"name_filter": schema.StringAttribute{
				Description: "Optional filter to match user names (case-insensitive substring match).",
				Optional:    true,
			},
			"category_id": schema.StringAttribute{
				Description: "Optional filter to match users by category ID.",
				Optional:    true,
			},
			"group_id": schema.StringAttribute{
				Description: "Optional filter to match users by group ID.",
				Optional:    true,
			},
			"role": schema.StringAttribute{
				Description: "Optional filter to match users by role (admin, manager, advanced, user).",
				Optional:    true,
			},
			"active": schema.BoolAttribute{
				Description: "Optional filter to match users by active status.",
				Optional:    true,
			},
			"users": schema.ListNestedAttribute{
				Description: "List of users that match the filters.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "User ID.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "User name.",
							Computed:    true,
						},
						"username": schema.StringAttribute{
							Description: "Username.",
							Computed:    true,
						},
						"uid": schema.StringAttribute{
							Description: "User UID.",
							Computed:    true,
						},
						"email": schema.StringAttribute{
							Description: "User email.",
							Computed:    true,
						},
						"active": schema.BoolAttribute{
							Description: "Whether the user is active.",
							Computed:    true,
						},
						"role": schema.StringAttribute{
							Description: "User role.",
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
						"secondary_groups": schema.ListAttribute{
							Description: "Secondary groups IDs.",
							Computed:    true,
							ElementType: types.StringType,
						},
						"provider": schema.StringAttribute{
							Description: "Authentication provider.",
							Computed:    true,
						},
						"email_verified": schema.BoolAttribute{
							Description: "Whether the email is verified.",
							Computed:    true,
						},
						"disclaimer_acknowledged": schema.BoolAttribute{
							Description: "Whether the disclaimer has been acknowledged.",
							Computed:    true,
						},
						"role_name": schema.StringAttribute{
							Description: "Role name (human readable).",
							Computed:    true,
						},
						"category_name": schema.StringAttribute{
							Description: "Category name (human readable).",
							Computed:    true,
						},
						"group_name": schema.StringAttribute{
							Description: "Group name (human readable).",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *usersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *usersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data usersDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Obtener todos los usuarios y aplicar filtros localmente
	users, err := d.client.GetUsers()
	
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read users, got error: %s", err))
		return
	}

	// Apply filters if provided
	nameFilter := data.NameFilter.ValueString()
	categoryFilter := data.CategoryID.ValueString()
	groupFilter := data.GroupID.ValueString()
	roleFilter := data.Role.ValueString()
	activeFilter := data.Active
	
	var filteredUsers []client.User
	for _, user := range users {
		// Apply name filter (case-insensitive substring match)
		if nameFilter != "" && !containsIgnoreCaseUsers(user.Name, nameFilter) {
			continue
		}
		
		// Apply category filter
		if categoryFilter != "" && user.Category != categoryFilter {
			continue
		}
		
		// Apply group filter
		if groupFilter != "" && user.Group != groupFilter {
			continue
		}
		
		// Apply role filter
		if roleFilter != "" && user.Role != roleFilter {
			continue
		}
		
		// Apply active filter
		if !activeFilter.IsNull() && user.Active != activeFilter.ValueBool() {
			continue
		}
		
		filteredUsers = append(filteredUsers, user)
	}

	// Map filtered users to model
	data.Users = make([]userModel, len(filteredUsers))
	for i, user := range filteredUsers {
		secondaryGroups := make([]types.String, len(user.SecondaryGroups))
		for j, sg := range user.SecondaryGroups {
			secondaryGroups[j] = types.StringValue(sg)
		}
		
		data.Users[i] = userModel{
			ID:                     types.StringValue(user.ID),
			Name:                   types.StringValue(user.Name),
			Username:               types.StringValue(user.Username),
			UID:                    types.StringValue(user.UID),
			Email:                  types.StringValue(user.Email),
			Active:                 types.BoolValue(user.Active),
			Role:                   types.StringValue(user.Role),
			Category:               types.StringValue(user.Category),
			Group:                  types.StringValue(user.Group),
			SecondaryGroups:        secondaryGroups,
			Provider:               types.StringValue(user.Provider),
			EmailVerified:          types.BoolValue(user.GetEmailVerified()),
			DisclaimerAcknowledged: types.BoolValue(user.GetDisclaimerAcknowledged()),
			RoleName:               types.StringValue(user.RoleName),
			CategoryName:           types.StringValue(user.CategoryName),
			GroupName:              types.StringValue(user.GroupName),
		}
	}

	data.ID = types.StringValue("users")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Helper function for case-insensitive substring matching
func containsIgnoreCaseUsers(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
