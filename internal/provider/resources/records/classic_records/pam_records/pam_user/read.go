// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamuser

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *PamUserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state PamUserResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError(utils.ERR_MSG_PROVIDER_CONFIGURATION_ERROR, err.Error())
		return
	}

	id := strings.TrimSpace(state.Id.ValueString())
	if id == "" {
		resp.Diagnostics.AddError(ErrSummaryReadFailed, "PAM User record id is empty")
		return
	}

	if err := utils.SyncDown(ctx, r.ApiManager); err != nil {
		resp.Diagnostics.AddError(utils.ErrSummarySyncDownFailed, err.Error())
		return
	}

	// Phase 1: read the vault record.
	command := fmt.Sprintf("%s '%s' %s", utils.CmdGet, id, utils.FlagFormatJSON)
	apiResp, err := r.ApiManager.ExecuteCommand(ctx, command, ErrDetailReadFailed)
	if err != nil {
		if errors.Is(err, api.ErrResourceNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(ErrSummaryReadFailed, err.Error())
		return
	}

	if apiResp == nil || apiResp.Data == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	var rec utils.VaultRecordGetResponse
	if err := utils.UnmarshalApiResponse(apiResp.Data, &rec); err != nil {
		resp.Diagnostics.AddError(ErrSummaryReadFailed, err.Error())
		return
	}

	if rec.Type != "" && rec.Type != utils.RecordTypePamUser {
		resp.Diagnostics.AddError(
			ErrSummaryReadFailed,
			fmt.Sprintf("vault record type is %q, expected %q", rec.Type, utils.RecordTypePamUser),
		)
		return
	}

	mapVaultRecordToState(&rec, &state)

	// Phase 2: always read rotation info — needed for import to auto-discover rotation state.
	rotCmd := fmt.Sprintf("%s %s '%s'", CmdPamRotationInfo, FlagRecordShort, id)
	rotResp, err2 := r.ApiManager.ExecuteCommand(ctx, rotCmd, ErrDetailRotationInfoFailed)
	if err2 == nil && rotResp != nil {
		messages := parseFlexibleMessageToLines(rotResp.Message.String())
		if hasRotationData(messages) {
			parseRotationInfoMessage(messages, state.RotationSettings, &state)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// parseFlexibleMessageToLines converts the FlexibleMessage string (may be a JSON array or plain string)
// into individual lines for parsing.
func parseFlexibleMessageToLines(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}

	// FlexibleMessage serializes arrays as JSON strings — try to unmarshal as []string first.
	var arr []string
	if err := json.Unmarshal([]byte(raw), &arr); err == nil {
		return arr
	}

	// Fallback: split on newlines for plain text.
	return strings.Split(raw, "\n")
}
