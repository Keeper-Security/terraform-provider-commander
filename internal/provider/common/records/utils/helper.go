// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package utils

import (
	"context"
	"fmt"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// FetchVaultRecord runs `get <recordUID> --format json --include-dag`.
func FetchVaultRecord(ctx context.Context, apiManager *api.ApiManager, recordUID string) (*api.RequestResultResponse, error) {
	command := fmt.Sprintf("%s '%s' %s %s", utils.CmdGet, recordUID, utils.FlagFormatJSON, utils.FlagIncludeDag)
	return apiManager.ExecuteCommand(ctx, command, utils.ErrSummaryFetchVaultRecordFailed)
}

// MoveRecordFromSourceToDestination moves a record when plan and state folder paths differ.
func MoveRecordFromSourceToDestination(ctx context.Context, apiManager *api.ApiManager, recordUID string, planFolderData string, stateFolderData string) error {
	if planFolderData == stateFolderData {
		return nil
	}

	dest := planFolderData
	if dest == "" {
		dest = "/"
	}

	command := fmt.Sprintf("%s '%s' '%s' %s", utils.CmdMv, recordUID, dest, utils.FlagForce)
	_, err := apiManager.ExecuteCommand(ctx, command, utils.ErrSummaryMoveRecordFailed)
	return err
}

// MapBaseVaultRecord maps record_uid, title, notes, and folder from API onto base.
func MapBaseVaultRecord(rec *utils.VaultRecordGetResponse, stateFolderLocation types.String, base *BaseVaultRecordModel) {
	if rec == nil || base == nil {
		return
	}
	if strings.TrimSpace(rec.RecordUID) != "" {
		base.Id = types.StringValue(strings.TrimSpace(rec.RecordUID))
	}
	base.Title = utils.StringOrNull(rec.Title)
	base.Notes = utils.StringOrNull(rec.Notes)
	base.FolderLocation = utils.ExtractFolderValue(rec.FolderLocation, stateFolderLocation)
}
