// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package folder

import (
	"context"
	"fmt"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SplitFolderPath splits a full vault path into parent and leaf name.
// Example: "Parent/MyFolder" -> parent "Parent", leaf "MyFolder".
// A path with no "/" returns empty parent and the whole string as leaf.
func SplitFolderPath(full string) (parent, leaf string) {
	full = strings.TrimSpace(full)
	if full == "" {
		return "", ""
	}
	i := strings.LastIndex(full, "/")
	if i < 0 {
		return "", full
	}
	return strings.TrimSpace(full[:i]), strings.TrimSpace(full[i+1:])
}

// EscapeDoubleQuotesForCLI escapes double quotes for use inside double-quoted shell arguments.
func EscapeDoubleQuotesForCLI(s string) string {
	return strings.ReplaceAll(s, `"`, `\"`)
}

// BuildFolderPath constructs the full folder path from name and optional folder_location.
// If folderLocation is non-empty, the result is "folderLocation/name"; otherwise just "name".
func BuildFolderPath(name, folderLocation string) string {
	name = strings.TrimSpace(name)
	folderLocation = strings.TrimSpace(folderLocation)
	if folderLocation == "" {
		return name
	}
	return folderLocation + "/" + name
}

// MvPathForCommander normalizes a vault path for Commander `mv`. Paths with no parent
// (no `/` — at vault root) are prefixed with `/` so the CLI targets the root folder.
func MvPathForCommander(full string) string {
	full = strings.TrimSpace(full)
	if full == "" {
		return full
	}
	if strings.HasPrefix(full, "/") {
		return full
	}
	parent, leaf := SplitFolderPath(full)
	if parent == "" {
		return "/" + leaf
	}
	return full
}

// MvMoveTargetParent returns the destination parent folder for Commander `mv`.
// Example: "Templates/test4/MyFolder" -> "Templates/test4".
// "MyFolder" (vault root) -> "/".
func MvMoveTargetParent(planPath string) string {
	planPath = strings.TrimSpace(planPath)
	if planPath == "" {
		return planPath
	}
	trim := planPath
	if strings.HasPrefix(trim, "/") {
		trim = strings.TrimSpace(trim[1:])
	}
	parent, _ := SplitFolderPath(trim)
	parent = strings.TrimSpace(parent)
	if parent == "" {
		return "/"
	}
	return parent
}

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
	path := EscapeDoubleQuotesForCLI(folderName + "/" + record)
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
