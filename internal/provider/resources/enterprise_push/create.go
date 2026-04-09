// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package enterprisepush

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *EnterprisePushResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data EnterprisePushResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError(
			utils.ERR_MSG_PROVIDER_CONFIGURATION_ERROR,
			err.Error(),
		)
		return
	}

	emails := sortedSetStrings(data.Email)
	teams := sortedSetStrings(data.Team)
	if len(emails) == 0 && len(teams) == 0 {
		resp.Diagnostics.AddError(
			"Invalid Configuration",
			"At least one of email or team must be specified. Provide at least one email address or one team to push records to.",
		)
		return
	}

	filePath := data.FilePath.ValueString()
	content, fileData, err := readFileAndParseJSON(filePath)
	if err != nil {
		resp.Diagnostics.AddError("Read File Failed", err.Error())
		return
	}

	if err := utils.RunWithManagedCompanyContext(ctx, r.ApiManager, data.ManagedCompany, func() error {
		command := buildEnterprisePushCommand(&data)
		_, err := r.ApiManager.ExecuteCommand(ctx, command, "Enterprise push failed", fileData)
		return err
	}, "Enterprise Push Failed", &resp.Diagnostics); err != nil {
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Deterministic ID: sha256(content + sorted emails + sorted teams + managed_company).
	// Any config change produces a new ID → Terraform replaces → push runs again.
	data.Id = types.StringValue(computeID(content, &data))
	data.FileContentSha256 = types.StringValue(contentSHA256Hex(content))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
