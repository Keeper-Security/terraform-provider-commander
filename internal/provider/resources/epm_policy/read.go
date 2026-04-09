// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package epmpolicy

import (
	"context"
	"errors"
	"fmt"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	commonepm "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/epm_policy"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *EpmPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state EpmPolicyResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
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

	if state.Id.IsNull() || state.Id.IsUnknown() || state.Id.ValueString() == "" {
		resp.Diagnostics.AddError(
			commonepm.ErrSummaryReadFailed,
			"Cannot read EPM policy: state has no policy id.",
		)
		return
	}

	err := utils.RunWithManagedCompanyContext(ctx, r.ApiManager, state.ManagedCompany, func() error {
		if err := utils.EpmSyncDown(ctx, r.ApiManager); err != nil {
			return err
		}

		apiResp, err := r.ApiManager.ExecuteCommand(ctx, buildViewCommand(state.Id.ValueString()), commonepm.ErrOpReadEpmPolicy)
		if err != nil {
			if errors.Is(err, api.ErrResourceNotFound) {
				resp.State.RemoveResource(ctx)
				return utils.ErrResourceRemoved
			}
			return err
		}

		var view utils.EpmPolicyResponse
		if err := utils.UnmarshalApiResponse(apiResp.Data, &view); err != nil {
			return err
		}

		if err := mapEpmPolicyResponseToModel(&view, &state); err != nil {
			return err
		}
		return nil
	}, commonepm.ErrSummaryReadFailed, &resp.Diagnostics)
	if err != nil && errors.Is(err, utils.ErrResourceRemoved) {
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// mapEpmPolicyResponseToModel maps utils.EpmPolicyResponse (epm policy view --format json) into state.
// Same pattern as enterprise role mapRoleReadResponseToModel: API types live in utils, domain mapping in common.
// managed_company is left unchanged.
func mapEpmPolicyResponseToModel(view *utils.EpmPolicyResponse, state *EpmPolicyResourceModel) error {
	prior := &commonepm.PolicyViewPriorSets{
		Control:            setToStrings(state.Control),
		UserGroups:         setToStrings(state.UserGroups),
		MachineCollections: setToStrings(state.MachineCollections),
		Applications:       setToStrings(state.Applications),
		DayFilter:          setToStrings(state.DayFilter),
		DateFilter:         setToStrings(state.DateFilter),
		TimeFilter:         setToStrings(state.TimeFilter),
	}
	mapped, err := commonepm.MapPolicyViewToAttributes(view, prior)
	if err != nil {
		return err
	}

	state.Id = types.StringValue(mapped.ID)
	state.PolicyName = types.StringValue(mapped.PolicyName)
	state.PolicyType = types.StringValue(mapped.PolicyType)
	state.Status = types.StringValue(mapped.Status)

	var setErr error
	state.Control, setErr = commonepm.StringSliceToStringSet(mapped.Control)
	if setErr != nil {
		return fmt.Errorf("control: %w", setErr)
	}
	state.UserGroups, setErr = commonepm.StringSliceToStringSet(mapped.UserGroups)
	if setErr != nil {
		return fmt.Errorf("user_groups: %w", setErr)
	}
	state.MachineCollections, setErr = commonepm.StringSliceToStringSet(mapped.MachineCollections)
	if setErr != nil {
		return fmt.Errorf("machine_collections: %w", setErr)
	}
	state.Applications, setErr = commonepm.StringSliceToStringSet(mapped.Applications)
	if setErr != nil {
		return fmt.Errorf("applications: %w", setErr)
	}
	state.DayFilter, setErr = commonepm.StringSliceToStringSet(mapped.DayFilter)
	if setErr != nil {
		return fmt.Errorf("day_filter: %w", setErr)
	}
	state.DateFilter, setErr = commonepm.StringSliceToStringSet(mapped.DateFilter)
	if setErr != nil {
		return fmt.Errorf("date_filter: %w", setErr)
	}
	state.TimeFilter, setErr = commonepm.StringSliceToStringSet(mapped.TimeFilter)
	if setErr != nil {
		return fmt.Errorf("time_filter: %w", setErr)
	}

	state.Message, state.RequirePolicyAcknowledgement = commonepm.NotificationAttributesFromMapped(
		mapped.Status, mapped.Message, mapped.RequirePolicyAcknowledgement,
	)

	return nil
}
