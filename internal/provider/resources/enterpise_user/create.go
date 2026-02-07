// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpiseuser

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
	if err := r.ensureApiManager(); err != nil {
		resp.Diagnostics.AddError(
			"Provider Configuration Error",
			err.Error(),
		)
		return
	}

	// Execute with managed company context if provided
	// ExecuteWithManagedCompanyContext handles context switching and enterprise-down internally
	err := utils.ExecuteWithManagedCompanyContext(ctx, r.apiManager, data.ManagedCompany, func() error {

		if !data.Teams.IsNull() && len(data.Teams.Elements()) > 0 {
			return fmt.Errorf("teams cannot be set when creating an enterprise user; add teams after the user is created via update")
		}

		roles, err := utils.FetchAndProcessRoles(ctx, r.apiManager, types.SetNull(types.StringType), data.Roles, "--add-role", "--remove-role")
		if err != nil {
			return err
		}
		if err := addUserBasicAttributes(ctx, r.apiManager, &data); err != nil {
			return err
		}

		if roles != "" {
			// Add Users and Roles to the recently created team
			command := fmt.Sprintf("enterprise-user '%s' %s -v", data.Id.ValueString(), roles)

			_, err = r.apiManager.ExecuteCommand(ctx, command, "Unable to add roles to the enterprise team")
			if err != nil {
				return err
			}
		}

		data.Status = types.StringValue(UserInvitedStatus)

		return nil
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Create Enterprise User Failed",
			err.Error(),
		)
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

	_, err := apiManager.ExecuteCommand(ctx, command, "Unable to add user")
	if err != nil {
		return err
	}

	// TODO: later user created users_id in date.id
	data.Id = data.Email

	return nil
}
