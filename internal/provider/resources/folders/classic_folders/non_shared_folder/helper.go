// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package nonsharedfolder

import (
	"context"
	"fmt"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	folderutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/folders/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// setElements returns the string elements of a types.Set; empty slice if null/unknown.
func setElements(s types.Set) []string {
	if s.IsNull() || s.IsUnknown() {
		return nil
	}
	var out []string
	for _, v := range s.Elements() {
		if sv, ok := v.(types.String); ok && !sv.IsNull() && !sv.IsUnknown() {
			out = append(out, sv.ValueString())
		}
	}
	return out
}

// setToMap converts a string set to a map for fast lookup.
func setToMap(s types.Set) map[string]bool {
	m := make(map[string]bool)
	for _, v := range setElements(s) {
		m[v] = true
	}
	return m
}

// LinkRecords links all records in the set into the folder via `ln '<record>' '<folderUID>'`.
func LinkRecords(ctx context.Context, apiManager *api.ApiManager, folderUID string, records types.Set) error {
	for _, record := range setElements(records) {
		cmd := fmt.Sprintf("%s '%s' '%s'", CmdLn, record, folderUID)
		if _, err := apiManager.ExecuteCommand(ctx, cmd, ErrOpLinkRecord); err != nil {
			return err
		}
	}
	return nil
}

// UnlinkRecord removes a record from a folder via `rm -f "<folderName>/<record>"`.
func UnlinkRecord(ctx context.Context, apiManager *api.ApiManager, folderName, record string) error {
	path := folderutils.EscapeDoubleQuotesForCLI(folderName + "/" + record)
	cmd := fmt.Sprintf(`%s %s "%s"`, CmdRm, FlagForce, path)
	if _, err := apiManager.ExecuteCommand(ctx, cmd, ErrOpUnlinkRecord); err != nil {
		return err
	}
	return nil
}

// SyncFolderRecords syncs records: links added records, unlinks removed records.
func SyncFolderRecords(ctx context.Context, apiManager *api.ApiManager, folderUID, folderName string, planRecords, stateRecords types.Set) error {
	planSet := setToMap(planRecords)
	stateSet := setToMap(stateRecords)

	for record := range stateSet {
		if !planSet[record] {
			if err := UnlinkRecord(ctx, apiManager, folderName, record); err != nil {
				return err
			}
		}
	}

	for record := range planSet {
		if !stateSet[record] {
			cmd := fmt.Sprintf("%s '%s' '%s'", CmdLn, record, folderUID)
			if _, err := apiManager.ExecuteCommand(ctx, cmd, ErrOpLinkRecord); err != nil {
				return err
			}
		}
	}

	return nil
}
