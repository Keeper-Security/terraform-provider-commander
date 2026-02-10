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
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *EnterpriseTeamResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data EnterpriseTeamResourceModel

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
		users, err := utils.FetchAndProcessUsers(ctx, r.ApiManager, types.SetNull(types.StringType), data.Users)
		if err != nil {
			return err
		}
		roles, err := utils.FetchAndProcessRoles(ctx, r.ApiManager, types.SetNull(types.StringType), data.Roles)
		if err != nil {
			return err
		}
		if err := addTeamBasicAttributes(ctx, r.ApiManager, &data); err != nil {
			return err
		}
		if users != "" {
			command := fmt.Sprintf("enterprise-team '%s' %s -v", data.Id.ValueString(), users)
			if _, err = r.ApiManager.ExecuteCommand(ctx, command, "Unable to add users to the enterprise team"); err != nil {
				return err
			}
		}
		if roles != "" {
			command := fmt.Sprintf("enterprise-team '%s' %s -v", data.Id.ValueString(), roles)
			if _, err = r.ApiManager.ExecuteCommand(ctx, command, "Unable to add roles to the enterprise team"); err != nil {
				return err
			}
		}
		return nil
	}, "Create Enterprise Team Failed", &resp.Diagnostics); err != nil {
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the ID in the state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func addTeamBasicAttributes(ctx context.Context, apiManager *api.ApiManager, data *EnterpriseTeamResourceModel) error {
	var parts []string

	parts = append(parts, "enterprise-team")

	// Required parameters
	parts = append(parts, fmt.Sprintf("--add '%s'", data.Name.ValueString()))

	// Optional parameters
	if !data.RestrictEdit.IsNull() && data.RestrictEdit.ValueBool() {
		parts = append(parts, "--restrict-edit on")
	}

	if !data.RestrictShare.IsNull() && data.RestrictShare.ValueBool() {
		parts = append(parts, "--restrict-share on")
	}

	if !data.RestrictView.IsNull() && data.RestrictView.ValueBool() {
		parts = append(parts, "--restrict-view on")
	}

	if !data.Node.IsNull() {
		parts = append(parts, fmt.Sprintf("--node '%s'", data.Node.ValueString()))
	}

	command := strings.Join(parts, " ")

	createdTeamResponse, err := apiManager.ExecuteCommand(ctx, command, "Unable to create enterprise team")
	if err != nil {
		return err
	}

	createdTeamUid, isTeamIdExtracted := extractTeamIdFromCreateTeamResponse(string(createdTeamResponse.Message))

	if isTeamIdExtracted {
		data.Id = types.StringValue(createdTeamUid)
	} else {
		return fmt.Errorf("failed to extract team id from create team response: %s", string(createdTeamResponse.Message))
	}

	return nil
}
