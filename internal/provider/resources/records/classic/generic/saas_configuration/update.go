// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package saasconfiguration

import (
	"context"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonrecordsaasconfiguration "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/saas_configuration"
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *SaasConfigurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state SaasConfigurationResourceModel
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
		resp.Diagnostics.AddError(ErrSummaryUpdateFailed, "SaaS configuration record id is empty")
		return
	}

	// Restrict changing gateway or configuration attributes (same behaviour as UI) - but not
	// when state's value is empty, e.g. right after import, and the plan
	// is providing a real value for the first time rather than changing an existing one.
	gatewayChanged := state.Gateway.ValueString() != "" && plan.Gateway.ValueString() != state.Gateway.ValueString()
	configurationChanged := state.Configuration.ValueString() != "" && plan.Configuration.ValueString() != state.Configuration.ValueString()
	if gatewayChanged || configurationChanged {
		resp.Diagnostics.AddError(ErrSummaryUpdateFailed, commonrecordsaasconfiguration.ErrOpChangeGatewayOrConfiguration)
		return
	}

	if err := commonrecordsutils.MoveRecordFromSourceToDestination(ctx, r.ApiManager, state.Id.ValueString(), plan.FolderLocation.ValueString(), state.FolderLocation.ValueString()); err != nil {
		resp.Diagnostics.AddError(utils.ErrSummaryMoveRecordFailed, err.Error())
		return
	}

	// SaaS Type/fields need (re)validating, and the record needs (re)linking to its gateway,
	// when SaaS Type itself changed, or when gateway/configuration are being provided for the
	// first time (e.g. right after import).
	// Otherwise it was already validated (at create, or a prior update) and redoing it
	// here would just be extra live Commander calls on every unrelated update.
	saasTypeChanged := commonrecordsaasconfiguration.SaasTypeFromCustom(plan.Custom) != commonrecordsaasconfiguration.SaasTypeFromCustom(state.Custom)
	firstTimeLink := (state.Gateway.ValueString() == "" || state.Configuration.ValueString() == "") &&
		plan.Gateway.ValueString() != "" && plan.Configuration.ValueString() != ""

	needsSaasUpdate := saasTypeChanged || firstTimeLink
	if saasType := commonrecordsaasconfiguration.SaasTypeFromCustom(plan.Custom); saasType != "" && needsSaasUpdate {
		// validate the saas type
		entries, err := commonrecordsaasconfiguration.FetchSaasPluginList(ctx, r.ApiManager, plan.Gateway.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(commonrecordsaasconfiguration.ErrOpListSaasPluginsForGateway(plan.Gateway.ValueString()), err.Error())
			return
		}
		if err := commonrecordsaasconfiguration.ValidateSaasType(saasType, entries); err != nil {
			resp.Diagnostics.AddError(commonrecordsaasconfiguration.ErrSummaryInvalidSaasType, err.Error())
			return
		}

		// validate the saas type (plugin) fields
		specs, err := commonrecordsaasconfiguration.FetchSaasFieldSpecs(ctx, r.ApiManager, plan.Gateway.ValueString(), saasType)
		if err != nil {
			resp.Diagnostics.AddError(commonrecordsaasconfiguration.ErrOpGetSaasPluginInfo, err.Error())
			return
		}
		if err := commonrecordsaasconfiguration.ValidateSaasFields(saasType, plan.Custom, specs); err != nil {
			resp.Diagnostics.AddError(commonrecordsaasconfiguration.ErrSummaryInvalidSaasFields, err.Error())
			return
		}
	}

	if commonrecordsaasconfiguration.UpdateHasMutations(plan.SaasConfigurationModel, state.SaasConfigurationModel) {
		cmd := commonrecordsaasconfiguration.BuildUpdateCommand(utils.CmdRecordUpdate, uid, plan.SaasConfigurationModel, state.SaasConfigurationModel)
		if _, err := r.ApiManager.ExecuteCommand(ctx, cmd, ErrDetailUpdateFailed); err != nil {
			resp.Diagnostics.AddError(ErrSummaryUpdateFailed, err.Error())
			return
		}
	}

	if needsSaasUpdate {
		if err := commonrecordsaasconfiguration.LinkToGateway(ctx, r.ApiManager, plan.Gateway.ValueString(), plan.Configuration.ValueString(), uid); err != nil {
			resp.Diagnostics.AddError(ErrSummaryUpdateFailed, err.Error())
			return
		}
	}

	if err := classic_share.SyncSharePermissions(ctx, r.ApiManager, uid, plan.Share, state.Share); err != nil {
		resp.Diagnostics.AddError(ErrSummaryUpdateFailed, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
