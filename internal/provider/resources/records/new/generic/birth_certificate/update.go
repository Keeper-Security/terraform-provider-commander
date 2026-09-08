// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package birthcertificate

import (
	"context"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonrecordbirthcertificate "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/birth_certificate"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *BirthCertificateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state BirthCertificateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError(utils.ERR_MSG_PROVIDER_CONFIGURATION_ERROR, err.Error())
		return
	}

	if err := utils.SyncDown(ctx, r.ApiManager); err != nil {
		resp.Diagnostics.AddError(utils.ErrSummarySyncDownFailed, err.Error())
		return
	}

	plan.Id = state.Id

	uid := strings.TrimSpace(plan.Id.ValueString())
	if uid == "" {
		resp.Diagnostics.AddError(ErrSummaryUpdateFailed, "Birth Certificate record id is empty")
		return
	}

	if err := utils.MoveNsfFromSourceToDestination(ctx, r.ApiManager, state.Id.ValueString(), plan.FolderLocation.ValueString(), state.FolderLocation.ValueString()); err != nil {
		resp.Diagnostics.AddError(utils.ErrSummaryNsfMoveRecordFailed, err.Error())
		return
	}

	if commonrecordbirthcertificate.UpdateHasMutations(plan.BirthCertificateModel, state.BirthCertificateModel) {
		cmd := commonrecordbirthcertificate.BuildUpdateCommand(utils.CmdNsfRecordUpdate, uid, plan.BirthCertificateModel, state.BirthCertificateModel)
		if _, err := r.ApiManager.ExecuteCommand(ctx, cmd, ErrDetailUpdateFailed); err != nil {
			resp.Diagnostics.AddError(ErrSummaryUpdateFailed, err.Error())
			return
		}
	}

	if err := new_share.SyncSharePermissions(ctx, r.ApiManager, new_share.CmdNsfShareRecord, uid, plan.Share, state.Share); err != nil {
		resp.Diagnostics.AddError(ErrSummaryUpdateFailed, err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
