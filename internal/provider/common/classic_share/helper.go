// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package classic_share

import (
	"context"
	"fmt"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// BuildGrantCommand builds the share-record command that grants/updates the
// per-user share flags for a classic vault record:
//
//	share-record --email '<email>' '<recordUID>' [--share] [--write]
//
// When both canShare and canEdit are false (downgrade-to-viewer case) the
// command explicitly strips the share/write flags while keeping view access:
//
//	share-record --email '<email>' '<recordUID>' --action revoke --share --write
//
// Use this for emails present in the plan (whether new or with changed
// permissions). Use BuildRevokeCommand for emails removed from the plan.
func BuildGrantCommand(recordUID, email string, canShare, canEdit bool) string {
	parts := []string{
		CmdShareRecord,
		FlagEmail, quote(email),
		quote(recordUID),
	}
	switch {
	case !canShare && !canEdit:
		parts = append(parts, FlagAction, ActionRevoke, FlagShare, FlagWrite)
	default:
		if canShare {
			parts = append(parts, FlagShare)
		}
		if canEdit {
			parts = append(parts, FlagWrite)
		}
	}
	return strings.Join(parts, " ")
}

// BuildRevokeCommand builds the share-record command that fully removes a
// user from the classic vault record share:
//
//	share-record --email '<email>' '<recordUID>' --action revoke
func BuildRevokeCommand(recordUID, email string) string {
	parts := []string{
		CmdShareRecord,
		FlagEmail, quote(email),
		quote(recordUID),
		FlagAction, ActionRevoke,
	}
	return strings.Join(parts, " ")
}

// SyncSharePermissions reconciles plan vs state for the `share` attribute
// against the API:
//
//   - emails in state but not in plan -> revoke
//   - emails in plan but not in state -> grant
//   - emails in both, perms changed -> grant (CLI treats grant as upsert)
//   - emails in both, perms unchanged -> skip (no API call)
//
// recordUID is the classic vault record UID. If planShare or stateShare is
// null or unknown it is treated as an empty map.
//
// SyncSharePermissions stops at the first API error and returns it.
func SyncSharePermissions(ctx context.Context, apiManager *api.ApiManager, recordUID string, planShare, stateShare types.Map) error {
	if recordUID == "" {
		return fmt.Errorf("share permissions sync: recordUID is empty")
	}

	planEntries := mapToPermEntries(planShare)
	stateEntries := mapToPermEntries(stateShare)

	for email := range stateEntries {
		if _, stillPresent := planEntries[email]; !stillPresent {
			cmd := BuildRevokeCommand(recordUID, email)
			if _, err := apiManager.ExecuteCommand(ctx, cmd, ErrOpShareRevoke); err != nil {
				return err
			}
		}
	}

	for email, planPerm := range planEntries {
		if statePerm, inState := stateEntries[email]; inState && statePerm == planPerm {
			continue
		}
		cmd := BuildGrantCommand(recordUID, email, planPerm.CanShare, planPerm.CanEdit)
		if _, err := apiManager.ExecuteCommand(ctx, cmd, ErrOpShareGrant); err != nil {
			return err
		}
	}

	return nil
}

// MapResponseToModel populates m.Share from the API's user_permissions array.
// Entries with empty username are dropped silently.
//
// When the filtered set is empty (no managed shares from the API), m.Share is
// set to null rather than an empty map. The schema's MapNonEmptyValidator
// rejects `share = {}` in config, so null is the only way "no managed shares"
// can be expressed; producing null here keeps the config-vs-state diff clean
// (null == null).
func MapResponseToModel(permissions []UserPermissionEntry, m *ShareModel) error {
	if m == nil {
		return fmt.Errorf("share model is nil")
	}
	objectType := types.ObjectType{AttrTypes: SharePermissionsObjectType()}
	elements := make(map[string]attr.Value, len(permissions))
	for _, p := range permissions {
		if strings.TrimSpace(p.Username) == "" {
			continue
		}
		obj, diags := types.ObjectValue(SharePermissionsObjectType(), map[string]attr.Value{
			AttrCanShare: types.BoolValue(p.Shareable),
			AttrCanEdit:  types.BoolValue(p.Editable),
		})
		if diags.HasError() {
			return fmt.Errorf("unable to build share entry for %q: %s", p.Username, diags)
		}
		elements[p.Username] = obj
	}
	if len(elements) == 0 {
		m.Share = types.MapNull(objectType)
		return nil
	}
	mv, diags := types.MapValue(objectType, elements)
	if diags.HasError() {
		return fmt.Errorf("unable to build share map from API response: %s", diags)
	}
	m.Share = mv
	return nil
}

// permEntry is the internal Go representation of one share map element.
type permEntry struct {
	CanShare bool
	CanEdit  bool
}

// mapToPermEntries converts a types.Map of nested objects into a Go
// map[email]permEntry. Null/unknown maps and unknown element objects are
// skipped; null bool attributes default to false.
func mapToPermEntries(m types.Map) map[string]permEntry {
	out := map[string]permEntry{}
	if m.IsNull() || m.IsUnknown() {
		return out
	}
	for email, v := range m.Elements() {
		obj, ok := v.(types.Object)
		if !ok || obj.IsNull() || obj.IsUnknown() {
			continue
		}
		attrs := obj.Attributes()
		out[email] = permEntry{
			CanShare: attrBool(attrs, AttrCanShare),
			CanEdit:  attrBool(attrs, AttrCanEdit),
		}
	}
	return out
}

// attrBool returns the bool value at key, or false if missing/null/unknown.
func attrBool(attrs map[string]attr.Value, key string) bool {
	v, ok := attrs[key].(types.Bool)
	if !ok || v.IsNull() || v.IsUnknown() {
		return false
	}
	return v.ValueBool()
}

// quote wraps s in single quotes so spaces and CLI-special characters are
// passed through unchanged. Any single quote inside s is doubled, matching
// normalizeCommandForShell in the api package.
func quote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}
