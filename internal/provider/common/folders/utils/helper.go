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

// BuildFolderPath constructs the full vault folder path from a leaf name and
// optional parent folder_location. If folderLocation is non-empty the result is
// "folderLocation/name"; otherwise just "name". Shared by non_shared_folder,
// shared_folder and new_folder resources.
func BuildFolderPath(name, folderLocation string) string {
	name = strings.TrimSpace(name)
	folderLocation = strings.TrimSpace(folderLocation)
	if folderLocation == "" {
		return name
	}
	return folderLocation + "/" + name
}

// EscapeDoubleQuotesForCLI escapes double quotes for use inside double-quoted
// shell arguments passed to Commander CLI commands.
func EscapeDoubleQuotesForCLI(s string) string {
	return strings.ReplaceAll(s, `"`, `\"`)
}

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

// FolderRecordsToSet converts a slice of record UIDs into a Terraform Set of strings. Returns
// SetNull(StringType) for empty input. Shared by non_shared_folder and
// new_folder read-flow mappers; keeps null/empty semantics consistent.
func FolderRecordsToSet(ctx context.Context, uids []string) (types.Set, error) {
	if len(uids) == 0 {
		return types.SetNull(types.StringType), nil
	}
	set, diags := types.SetValueFrom(ctx, types.StringType, uids)
	if diags.HasError() {
		return types.SetNull(types.StringType), fmt.Errorf("build records set: %v", diags)
	}
	return set, nil
}

// setElements returns the string elements of a types.Set; nil if null/unknown.
func setElements(s types.Set) []string {
	if s.IsNull() || s.IsUnknown() {
		return nil
	}
	out := make([]string, 0, len(s.Elements()))
	for _, v := range s.Elements() {
		if sv, ok := v.(types.String); ok && !sv.IsNull() && !sv.IsUnknown() {
			out = append(out, sv.ValueString())
		}
	}
	return out
}

// setToMap converts a string set to a map for O(1) membership lookup.
func setToMap(s types.Set) map[string]struct{} {
	elems := setElements(s)
	m := make(map[string]struct{}, len(elems))
	for _, v := range elems {
		m[v] = struct{}{}
	}
	return m
}

// LinkRecords links each record UID in the set into the folder via
// `<lnCmd> '<record>' '<folderUID>'`. lnCmd is `ln` for classic non-shared
// folders and `nsf-ln` for Nested Shared Folders.
func LinkRecords(ctx context.Context, apiManager *api.ApiManager, lnCmd, folderUID string, records types.Set) error {
	for _, record := range setElements(records) {
		cmd := fmt.Sprintf("%s '%s' '%s'", lnCmd, record, folderUID)
		if _, err := apiManager.ExecuteCommand(ctx, cmd, ErrOpLinkRecord); err != nil {
			return err
		}
	}
	return nil
}

// UnlinkRecord removes a single record from a classic non-shared folder via
// `<rmCmd> --force "<folderName>/<record>"`. The classic CLI takes a vault
// path (folder name + record UID joined by `/`).
func UnlinkRecord(ctx context.Context, apiManager *api.ApiManager, rmCmd, folderName, record string) error {
	path := EscapeDoubleQuotesForCLI(folderName + "/" + record)
	cmd := fmt.Sprintf(`%s %s "%s"`, rmCmd, utils.FlagForce, path)
	if _, err := apiManager.ExecuteCommand(ctx, cmd, ErrOpUnlinkRecord); err != nil {
		return err
	}
	return nil
}

// UnlinkRecordNsf removes a single record from a Nested Shared Folder via
// `<rmCmd> '<record>' --folder '<folderUID>' --force --operation unlink`.
// Unlike the classic form, nsf-rm takes the record UID positionally and the
// folder UID via --folder (not a path), and requires --operation unlink to
// distinguish unlinking from deleting the record outright.
func UnlinkRecordNsf(ctx context.Context, apiManager *api.ApiManager, rmCmd, folderUID, record string) error {
	cmd := fmt.Sprintf("%s '%s' %s '%s' %s %s unlink",
		rmCmd, record, utils.FlagFolder, folderUID, utils.FlagForce, utils.FlagOperation)
	if _, err := apiManager.ExecuteCommand(ctx, cmd, ErrOpUnlinkRecord); err != nil {
		return err
	}
	return nil
}

// UnlinkRecordFromFolderFn is the strategy the caller supplies to SyncFolderRecords for
// the unlink half of the reconcile. Use UnlinkRecord (classic, path-based) or
// UnlinkRecordNsf (nsf-rm, flag-based) when binding the closure, or roll your
// own for new flavors.
type UnlinkRecordFromFolderFn func(record string) error

// SyncFolderRecords reconciles the record set on a folder against state:
// unlinks records present in state but missing from plan via unlinkFn, then
// links records added in plan via `<lnCmd> '<record>' '<folderUID>'`. The
// link CLI shape is the same for classic `ln` and `nsf-ln`; the unlink shape
// differs, hence the closure.
func SyncFolderRecords(ctx context.Context, apiManager *api.ApiManager, lnCmd, folderUID string, planRecords, stateRecords types.Set, unlinkFn UnlinkRecordFromFolderFn) error {
	planSet := setToMap(planRecords)
	stateSet := setToMap(stateRecords)

	for record := range stateSet {
		if _, keep := planSet[record]; keep {
			continue
		}
		if err := unlinkFn(record); err != nil {
			return err
		}
	}

	for record := range planSet {
		if _, existed := stateSet[record]; existed {
			continue
		}
		cmd := fmt.Sprintf("%s '%s' '%s'", lnCmd, record, folderUID)
		if _, err := apiManager.ExecuteCommand(ctx, cmd, ErrOpLinkRecord); err != nil {
			return err
		}
	}

	return nil
}

// ExtractFolderUIDFromCreateResponse pulls KeyFolderUID ("folder_uid") out of
// the API response Data (a JSON-decoded map) returned by Commander folder
// create commands (mkdir/mkdir --shared-folder). If folder_uid is missing,
// empty, or not a string (e.g. a number), an error or its stringified form is
// returned accordingly. Shared by non_shared_folder and classic shared_folder
// resource Create flows.
func ExtractFolderUIDFromCreateResponse(data any) (string, error) {
	m, _ := data.(map[string]interface{})
	v, ok := m[KeyFolderUID]
	if !ok || v == nil {
		return "", fmt.Errorf("API response missing %s", KeyFolderUID)
	}
	if s, ok := v.(string); ok {
		if s == "" {
			return "", fmt.Errorf("API response %s is empty", KeyFolderUID)
		}
		return s, nil
	}
	return fmt.Sprintf("%v", v), nil
}
