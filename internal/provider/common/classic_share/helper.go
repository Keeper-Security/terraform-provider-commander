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

// BuildGrantCommand builds a share-record grant that turns on only the listed
// permission flags. Grant is merge/update: flags passed are set to true and
// existing permissions are left unchanged.
//
//	share-record --email '<email>' '<recordUID>' [--share] [--write]
//
// Pass neither flag for view-only on a new share.
func BuildGrantCommand(recordUID, email string, canShare, canEdit bool) string {
	parts := shareCommandBase(recordUID, email)
	if canShare {
		parts = append(parts, FlagShare)
	}
	if canEdit {
		parts = append(parts, FlagWrite)
	}
	return strings.Join(parts, " ")
}

// BuildRevokeFlagsCommand strips the listed permission flags while keeping
// view access. Revoke is selective: only the flags passed are removed.
//
//	share-record --email '<email>' '<recordUID>' --action revoke [--share] [--write]
func BuildRevokeFlagsCommand(recordUID, email string, stripShare, stripEdit bool) string {
	parts := shareCommandBase(recordUID, email)
	parts = append(parts, FlagAction, ActionRevoke)
	if stripShare {
		parts = append(parts, FlagShare)
	}
	if stripEdit {
		parts = append(parts, FlagWrite)
	}
	return strings.Join(parts, " ")
}

// BuildRevokeCommand builds the share-record command that fully removes a
// user from the classic vault record share:
//
//	share-record --email '<email>' '<recordUID>' --action revoke
func BuildRevokeCommand(recordUID, email string) string {
	parts := shareCommandBase(recordUID, email)
	parts = append(parts, FlagAction, ActionRevoke)
	return strings.Join(parts, " ")
}

// BuildSharePermissionSyncCommands returns the minimal share-record commands
// needed to move from state to plan for one email. Grant only adds flags;
// revoke only removes flags — so both state and plan must be considered.
func BuildSharePermissionSyncCommands(recordUID, email string, state, plan PermEntry, inState bool) []string {
	// Permissions already match Terraform state — nothing to send to the API.
	if inState && state == plan {
		return nil
	}

	// User is new in plan: issue a single grant with the desired can_share/can_edit
	// flags. Omitting both flags gives view-only access.
	if !inState {
		return []string{BuildGrantCommand(recordUID, email, plan.CanShare, plan.CanEdit)}
	}

	var cmds []string

	// Existing user with changed permissions. The share-record CLI does not replace
	// permissions in one shot: grant only turns flags on, revoke only turns flags off.
	// Compare state vs plan and revoke any permission that was removed.
	stripShare := state.CanShare && !plan.CanShare
	stripEdit := state.CanEdit && !plan.CanEdit
	if stripShare || stripEdit {
		cmds = append(cmds, BuildRevokeFlagsCommand(recordUID, email, stripShare, stripEdit))
	}

	// Then grant any permission that was added in plan but was not in state.
	addShare := !state.CanShare && plan.CanShare
	addEdit := !state.CanEdit && plan.CanEdit
	if addShare || addEdit {
		cmds = append(cmds, BuildGrantCommand(recordUID, email, addShare, addEdit))
	}
	return cmds
}

func shareCommandBase(recordUID, email string) []string {
	return []string{
		CmdShareRecord,
		FlagEmail, quote(email),
		quote(recordUID),
	}
}

// SyncSharePermissions reconciles plan vs state for the `share` attribute
// against the API:
//
//   - emails in state but not in plan -> full revoke
//   - emails in plan but not in state -> grant desired flags
//   - emails in both, perms changed -> revoke dropped flags, then grant new ones
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

	// Revoke permissions for emails that are in state but not in plan.
	for email := range stateEntries {
		if _, stillPresent := planEntries[email]; !stillPresent {
			cmd := BuildRevokeCommand(recordUID, email)
			if _, err := apiManager.ExecuteCommand(ctx, cmd, ErrOpShareRevoke); err != nil {
				return err
			}
		}
	}

	// Apply adds and permission changes for every email still in plan. Users removed
	// from plan were fully revoked in the loop above.
	for email, planPerm := range planEntries {
		statePerm, inState := stateEntries[email]
		for _, cmd := range BuildSharePermissionSyncCommands(recordUID, email, statePerm, planPerm, inState) {
			// Run each generated command. Flag-stripping commands use --action revoke;
			// grant commands have no --action. Pick the matching error label for logs.
			errOp := ErrOpShareGrant
			if strings.Contains(cmd, FlagAction+" "+ActionRevoke) {
				errOp = ErrOpShareRevoke
			}
			if _, err := apiManager.ExecuteCommand(ctx, cmd, errOp); err != nil {
				return err
			}
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
		// Skip owner entries as they are managed by Keeper and are not tracked in Terraform state.
		// Will not add them to the share map.
		if p.Owner {
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

// PermEntry is the internal Go representation of one share map element.
type PermEntry struct {
	CanShare bool
	CanEdit  bool
}

// mapToPermEntries converts a types.Map of nested objects into a Go
// map[email]permEntry. Null/unknown maps and unknown element objects are
// skipped; null bool attributes default to false.
func mapToPermEntries(m types.Map) map[string]PermEntry {
	out := map[string]PermEntry{}
	if m.IsNull() || m.IsUnknown() {
		return out
	}
	for email, v := range m.Elements() {
		obj, ok := v.(types.Object)
		if !ok || obj.IsNull() || obj.IsUnknown() {
			continue
		}
		attrs := obj.Attributes()
		out[email] = PermEntry{
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
