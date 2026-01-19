// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriseteam

import (
	"context"
	"fmt"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &EnterpriseTeamResource{}
var _ resource.ResourceWithConfigure = &EnterpriseTeamResource{}

type EnterpriseTeamResource struct {
	apiManager *api.ApiManager
}

type EnterpriseTeamResourceModel struct {
	Id             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	RestrictEdit   types.Bool   `tfsdk:"restrict_record_edit"`
	RestrictShare  types.Bool   `tfsdk:"restrict_record_re_share"`
	RestrictView   types.Bool   `tfsdk:"enable_privacy_screen"`
	Users          types.Set    `tfsdk:"users"`
	Roles          types.Set    `tfsdk:"roles"`
	Node           types.String `tfsdk:"node"`
	ManagedCompany types.String `tfsdk:"managed_company"`
}

func (r *EnterpriseTeamResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_enterprise_team"
}

func (r *EnterpriseTeamResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "The ID of the enterprise team.",
				MarkdownDescription: "The ID of the enterprise team.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "Enterprise Team Name.",
				MarkdownDescription: "Enterprise Team Name.",
				Validators: []validator.String{
					nameValidator{},
				},
			},
			"restrict_record_edit": schema.BoolAttribute{
				Optional:            true,
				Description:         "Restrict record editing.",
				MarkdownDescription: "Restrict record editing.",
			},
			"restrict_record_re_share": schema.BoolAttribute{
				Optional:            true,
				Description:         "Restrict record re-sharing.",
				MarkdownDescription: "Restrict record re-sharing.",
			},
			"enable_privacy_screen": schema.BoolAttribute{
				Optional:            true,
				Description:         "Enable privacy screen.",
				MarkdownDescription: "Enable privacy screen.",
			},
			"users": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.Set{
					usersValidator{},
				},
				Description:         "Set of users in the enterprise team. Duplicate values are automatically prevented. Empty strings are not allowed.",
				MarkdownDescription: "Set of users in the enterprise team. Duplicate values are automatically prevented. Empty strings are not allowed.",
			},
			"roles": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.Set{
					rolesValidator{},
				},
				Description:         "Set of roles in the enterprise team. Duplicate values are automatically prevented. Empty strings are not allowed.",
				MarkdownDescription: "Set of roles in the enterprise team. Duplicate values are automatically prevented. Empty strings are not allowed.",
			},
			"node": schema.StringAttribute{
				Optional:            true,
				Description:         "Managing Node name or ID.",
				MarkdownDescription: "Managing Node name or ID.",
				Validators: []validator.String{
					nodeValidator{},
				},
			},
			"managed_company": schema.StringAttribute{
				Optional:            true,
				Description:         "Managed Company name or ID.",
				MarkdownDescription: "Managed Company name or ID.",
				Validators: []validator.String{
					managedCompanyValidator{},
				},
			},
		},
	}
}

func (r *EnterpriseTeamResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	apiManager, ok := req.ProviderData.(*api.ApiManager)
	if !ok {
		resp.Diagnostics.AddError(
			"Provider Configuration Error",
			fmt.Sprintf("The provider was not configured correctly. Expected API manager, but got: %T. Please check your provider configuration.", req.ProviderData),
		)
		return
	}

	r.apiManager = apiManager
}

// ensureApiManager validates that apiManager is configured and returns an error if not
func (r *EnterpriseTeamResource) ensureApiManager() error {
	if r.apiManager == nil {
		return fmt.Errorf("The Keeper Commander provider is not properly configured. Please ensure the provider is set up with valid service_mode_url and service_mode_api_key")
	}
	return nil
}

func (r *EnterpriseTeamResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data EnterpriseTeamResourceModel

	// Get planned data from Terraform
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate ApiManager is configured
	if err := r.ensureApiManager(); err != nil {
		resp.Diagnostics.AddError(
			"Provider Configuration Error",
			err.Error(),
		)
		return
	}

	// Execute with managed company context if provided
	err := utils.ExecuteWithManagedCompanyContext(ctx, r.apiManager, data.ManagedCompany, func() error {

		// For create, stateUsers and stateRoles are null/empty, only planUsers and planRoles have items to add
		users, err := fetchAndProcessUsers(ctx, r.apiManager, types.SetNull(types.StringType), data.Users)
		if err != nil {
			return err
		}

		roles, err := fetchAndProcessRoles(ctx, r.apiManager, types.SetNull(types.StringType), data.Roles)
		if err != nil {
			return err
		}

		// Build create command
		command := buildEnterpriseTeamAddCommand(data)

		_, err = r.apiManager.ExecuteCommand(ctx, command, "Unable to create enterprise team")
		if err != nil {
			return err
		}

		// TODO: WE WILL REMOVE THIS fetchTeamUidByName FUNCTION AFTER WE ARE RECIEVING TEAM UID IN THE RESPONSE IN COMMANDER CLI
		// Fetch the team UID by name
		teamUid, err := fetchTeamUidByName(ctx, r.apiManager, data.Name.ValueString())
		if err != nil {
			return err
		}

		// Combine users and roles flags
		var userRoleFlags []string
		if users != "" {
			userRoleFlags = append(userRoleFlags, users)
		}
		if roles != "" {
			userRoleFlags = append(userRoleFlags, roles)
		}

		if len(userRoleFlags) > 0 {
			// Add Users and Roles to the recently created team
			command = fmt.Sprintf("enterprise-team '%s' %s -v", teamUid, strings.Join(userRoleFlags, " "))

			_, err = r.apiManager.ExecuteCommand(ctx, command, "Unable to add users/roles to the enterprise team")
			if err != nil {
				return err
			}
		}

		data.Id = types.StringValue(teamUid)

		// Set the ID in the state
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		return nil

	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Create Enterprise Team Failed",
			err.Error(),
		)
		return
	}

}

func (r *EnterpriseTeamResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state EnterpriseTeamResourceModel

	// Get current state (old values)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate ApiManager is configured
	if err := r.ensureApiManager(); err != nil {
		resp.Diagnostics.AddError(
			"Provider Configuration Error",
			err.Error(),
		)
		return
	}

	// Use managed company from state if provided
	managedCompany := state.ManagedCompany

	// Execute with managed company context if provided
	err := utils.ExecuteWithManagedCompanyContext(ctx, r.apiManager, managedCompany, func() error {
		// Build command to get enterprise team info
		command := fmt.Sprintf("enterprise-info '%s' -t --format json --columns='users,roles,restricts,node' -q", state.Id.ValueString())

		// Execute the command
		apiResp, err := r.apiManager.ExecuteCommand(ctx, command, "Unable to retrieve enterprise team information")
		if err != nil {
			return fmt.Errorf("Read Enterprise Team Failed: %w", err)
		}

		// Parse the response
		teams, err := parseEnterpriseTeamReadResponse(apiResp.Data)
		if err != nil {
			return fmt.Errorf("Failed to parse team response: %w", err)
		}

		// Find the team matching the ID
		var teamInfo *EnterpriseTeamReadResponse
		for i := range teams {
			if teams[i].TeamUid == state.Id.ValueString() {
				teamInfo = &teams[i]
				break
			}
		}

		if teamInfo == nil {
			// Resource not found - remove from state
			resp.State.RemoveResource(ctx)
			return nil
		}

		// Map the response to the model
		if err := mapTeamReadResponseToModel(ctx, r.apiManager, *teamInfo, &state); err != nil {
			return fmt.Errorf("Failed to map team response to model: %w", err)
		}

		return nil
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Read Enterprise Team Failed",
			err.Error(),
		)
		return
	}

	// Set the updated state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *EnterpriseTeamResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan EnterpriseTeamResourceModel
	var state EnterpriseTeamResourceModel

	// Get planned data (new values)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state (old values)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate ApiManager is configured
	if err := r.ensureApiManager(); err != nil {
		resp.Diagnostics.AddError(
			"Provider Configuration Error",
			err.Error(),
		)
		return
	}

	// NOTE: We should not allow user to update managed company, bec. once team is created in managed company, if we allow user to update managed company then switching to that MC we will not able to find that team, so command will fail.
	if !plan.ManagedCompany.Equal(state.ManagedCompany) {
		resp.Diagnostics.AddError(
			"Managed Company Cannot Be Updated",
			"Cannot update the managed_company field. Once an enterprise team is created in a managed company, the managed company cannot be changed.",
		)
		return
	}

	// Use managed company from plan (or state if plan doesn't have it)
	managedCompany := plan.ManagedCompany
	if managedCompany.IsNull() || managedCompany.IsUnknown() {
		managedCompany = state.ManagedCompany
	}

	// Execute with managed company context if provided
	err := utils.ExecuteWithManagedCompanyContext(ctx, r.apiManager, managedCompany, func() error {
		command, err := buildEnterpriseTeamUpdateCommand(ctx, r.apiManager, &plan, &state)
		if err != nil {
			return fmt.Errorf("failed to build update command: %w", err)
		}

		_, err = r.apiManager.ExecuteCommand(ctx, command, "Unable to update enterprise team")
		if err != nil {
			return err
		}

		plan.Id = state.Id
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
		return nil
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Update Enterprise Team Failed",
			err.Error(),
		)
		return
	}
}

func (r *EnterpriseTeamResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state EnterpriseTeamResourceModel

	// Get state from Terraform
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate ApiManager is configured
	if err := r.ensureApiManager(); err != nil {
		resp.Diagnostics.AddError(
			"Provider Configuration Error",
			err.Error(),
		)
		return
	}

	// Execute with managed company context if provided
	err := utils.ExecuteWithManagedCompanyContext(ctx, r.apiManager, state.ManagedCompany, func() error {
		// Build delete command
		command := fmt.Sprintf("enterprise-team --delete --force '%s'", state.Id.ValueString())

		_, err := r.apiManager.ExecuteCommand(ctx, command, "Unable to delete enterprise team")
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Delete Enterprise Team Failed",
			err.Error(),
		)
		return
	}
}

func NewEnterpriseTeamResource() resource.Resource {
	return &EnterpriseTeamResource{}
}
