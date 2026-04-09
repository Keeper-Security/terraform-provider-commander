// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package epmpolicy

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	commonepm "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/epm_policy"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (d *EpmPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data EpmPolicyDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := d.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError(
			utils.ERR_MSG_PROVIDER_CONFIGURATION_ERROR,
			err.Error(),
		)
		return
	}

	policyID := strings.TrimSpace(data.Policy.ValueString())
	if policyID == "" {
		resp.Diagnostics.AddError(
			commonepm.ErrSummaryReadFailed,
			"The policy argument must be a non-empty EPM policy ID.",
		)
		return
	}

	err := utils.RunWithManagedCompanyContext(ctx, d.ApiManager, data.ManagedCompany, func() error {
		if err := utils.EpmSyncDown(ctx, d.ApiManager); err != nil {
			return err
		}

		apiResp, err := d.ApiManager.ExecuteCommand(ctx, commonepm.BuildPolicyViewCommand(policyID), commonepm.ErrOpReadEpmPolicy)
		if err != nil {
			if errors.Is(err, api.ErrResourceNotFound) {
				return fmt.Errorf("EPM policy %q not found", policyID)
			}
			return err
		}

		var view utils.EpmPolicyResponse
		if err := utils.UnmarshalApiResponse(apiResp.Data, &view); err != nil {
			return err
		}

		mapped, err := commonepm.MapPolicyViewToAttributes(&view, nil)
		if err != nil {
			return err
		}

		data.Id = types.StringValue(mapped.ID)
		data.PolicyName = types.StringValue(mapped.PolicyName)
		data.PolicyType = types.StringValue(mapped.PolicyType)
		data.Status = types.StringValue(mapped.Status)

		var setErr error
		data.Control, setErr = commonepm.StringSliceToStringSet(mapped.Control)
		if setErr != nil {
			return fmt.Errorf("control: %w", setErr)
		}
		data.UserGroups, setErr = commonepm.StringSliceToStringSet(mapped.UserGroups)
		if setErr != nil {
			return fmt.Errorf("user_groups: %w", setErr)
		}
		data.MachineCollections, setErr = commonepm.StringSliceToStringSet(mapped.MachineCollections)
		if setErr != nil {
			return fmt.Errorf("machine_collections: %w", setErr)
		}
		data.Applications, setErr = commonepm.StringSliceToStringSet(mapped.Applications)
		if setErr != nil {
			return fmt.Errorf("applications: %w", setErr)
		}
		data.DayFilter, setErr = commonepm.StringSliceToStringSet(mapped.DayFilter)
		if setErr != nil {
			return fmt.Errorf("day_filter: %w", setErr)
		}
		data.DateFilter, setErr = commonepm.StringSliceToStringSet(mapped.DateFilter)
		if setErr != nil {
			return fmt.Errorf("date_filter: %w", setErr)
		}
		data.TimeFilter, setErr = commonepm.StringSliceToStringSet(mapped.TimeFilter)
		if setErr != nil {
			return fmt.Errorf("time_filter: %w", setErr)
		}

		data.Message, data.RequirePolicyAcknowledgement = commonepm.NotificationAttributesFromMapped(
			mapped.Status, mapped.Message, mapped.RequirePolicyAcknowledgement,
		)

		return nil
	}, commonepm.ErrSummaryReadFailed, &resp.Diagnostics)
	if err != nil {
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
