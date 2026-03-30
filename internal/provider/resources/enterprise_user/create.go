// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package enterpriseuser

import (
	"context"
	"fmt"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *EnterpriseUserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var data EnterpriseUserResourceModel

	// Get planned data from Terraform
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate ApiManager is configured
	if err := r.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError(
			"Provider Configuration Error",
			err.Error(),
		)
		return
	}

	if err := utils.RunWithManagedCompanyContext(ctx, r.ApiManager, data.ManagedCompany, func() error {
		if !data.Teams.IsNull() && len(data.Teams.Elements()) > 0 {
			return fmt.Errorf("teams cannot be set when creating an enterprise user; add teams after the user is created via update")
		}
		roles, err := utils.FetchAndProcessRoles(ctx, r.ApiManager, types.SetNull(types.StringType), data.Roles, "--add-role", "--remove-role")
		if err != nil {
			return err
		}
		if err := addUserBasicAttributes(ctx, r.ApiManager, &data); err != nil {
			return err
		}
		if roles != "" {
			command := fmt.Sprintf("enterprise-user '%s' %s -v", data.Id.ValueString(), roles)
			if _, err = r.ApiManager.ExecuteCommand(ctx, command, "Unable to add roles to the enterprise team"); err != nil {
				return err
			}
		}
		data.Status = types.StringValue(UserInvitedStatus)
		return nil
	}, "Create Enterprise User Failed", &resp.Diagnostics); err != nil {
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the ID in the state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func addUserBasicAttributes(ctx context.Context, apiManager *api.ApiManager, data *EnterpriseUserResourceModel) error {
	var parts []string

	parts = append(parts, "enterprise-user")

	parts = append(parts, fmt.Sprintf("--add '%s'", data.Email.ValueString()))

	if !data.Node.IsNull() {
		parts = append(parts, fmt.Sprintf("--node '%s'", data.Node.ValueString()))
	}

	if !data.Name.IsNull() {
		parts = append(parts, fmt.Sprintf("--name '%s'", data.Name.ValueString()))
	}

	if !data.JobTitle.IsNull() {
		parts = append(parts, fmt.Sprintf("--job-title '%s'", data.JobTitle.ValueString()))
	}

	command := strings.Join(parts, " ")

	createUserResponse, err := apiManager.ExecuteCommand(ctx, command, "Unable to add user")
	if err != nil {
		return err
	}

	createdUserId, isCreatedUserIdExtracted := utils.ExtractUserIDFromCreateUserResponse(string(createUserResponse.Message))
	if isCreatedUserIdExtracted {
		data.Id = types.StringValue(createdUserId)
	} else {
		// Fallback: some API versions may not return "User ID : <id>" in the message; use email as id for compatibility
		data.Id = data.Email
	}

	return nil
}
