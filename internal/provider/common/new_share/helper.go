// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package new_share

import (
	"context"
	"fmt"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// BuildGrantCommand builds the CLI command for granting or updating a share:
//
//	<command> "<id>" --email=<email> --action=grant --role=<role>
//
// command is one of CmdNsfShareFolder / CmdNsfShareRecord. id is the folder UID or
// record UID. email and role are passed through; callers are expected to have
// already validated them via the schema validators.
func BuildGrantCommand(command, id, email, role string) string {
	return fmt.Sprintf(`%s "%s" %s=%s %s=%s %s=%s`,
		command,
		escapeDoubleQuotes(id),
		FlagEmail, quote(email),
		FlagAction, ActionGrant,
		FlagRole, quote(role),
	)
}

// BuildRevokeCommand builds the CLI command for revoking a share:
//
//	<command> "<id>" --email=<email> --action=<revoke|remove>
//
// The --action value is derived from command: nsf-share-folder uses
// ActionRemove, nsf-share-record uses ActionRevoke.
func BuildRevokeCommand(command, id, email string) string {
	return fmt.Sprintf(`%s "%s" %s=%s %s=%s`,
		command,
		escapeDoubleQuotes(id),
		FlagEmail, quote(email),
		FlagAction, revokeActionFor(command),
	)
}

// revokeActionFor returns the --action value used to remove a user from a
// share. Folders accept "remove"; records accept "revoke". Unknown commands
// fall back to ActionRevoke for forward compatibility.
func revokeActionFor(command string) string {
	if command == CmdNsfShareFolder {
		return ActionRemove
	}
	return ActionRevoke
}

// SyncSharePermissions reconciles plan vs state for the `share` attribute
// against the API:
//
//   - emails in state but not in plan -> revoke
//   - emails in plan but not in state -> grant
//   - emails in both, role changed -> grant (CLI treats grant as upsert)
//   - emails in both, role unchanged -> skip (no API call)
//
// command is CmdNsfShareFolder for folder resources or CmdNsfShareRecord for record
// resources. id is the folder/record UID. If planShare or stateShare is null
// or unknown it is treated as an empty map.
//
// SyncSharePermissions stops at the first API error and returns it.
func SyncSharePermissions(ctx context.Context, apiManager *api.ApiManager, command, id string, planShare, stateShare types.Map) error {
	if id == "" {
		return fmt.Errorf("share permissions sync: id is empty")
	}

	planEntries := mapToStringMap(planShare)
	stateEntries := mapToStringMap(stateShare)

	for email := range stateEntries {
		if _, stillPresent := planEntries[email]; !stillPresent {
			cmd := BuildRevokeCommand(command, id, email)
			if _, err := apiManager.ExecuteCommand(ctx, cmd, ErrOpShareRevoke); err != nil {
				return err
			}
		}
	}

	for email, planRole := range planEntries {
		if stateRole, inState := stateEntries[email]; inState && stateRole == planRole {
			continue
		}
		cmd := BuildGrantCommand(command, id, email, planRole)
		if _, err := apiManager.ExecuteCommand(ctx, cmd, ErrOpShareGrant); err != nil {
			return err
		}
	}

	return nil
}

// MapResponseToModel populates m.Share from the API's user_permissions array.
// Owner entries are silently dropped (owner is managed by Keeper and is not
// tracked in Terraform state). Owner is detected via role == RoleOwner
// (for NSF folder shape) or the owner boolean (for NSF/classic record shape). Entries
// with empty principal or empty role are also dropped.
//
// Principal resolution matches both NSF response shapes documented on
// utils.UserPermissionEntry:
//   - folders: accessor + role
//   - records: username + role (accessor is empty)
//
// When the filtered set is empty (e.g. the API returned only the owner row),
// m.Share is set to null rather than an empty map. The schema's
// MapNonEmptyValidator rejects `share = {}` in config, so null is the only
// way "no managed shares" can be expressed; producing null here keeps the
// config-vs-state diff clean (null == null).
func MapResponseToModel(permissions []UserPermissionEntry, m *ShareModel) error {
	if m == nil {
		return fmt.Errorf("share model is nil")
	}
	elements := make(map[string]attr.Value, len(permissions))
	for _, p := range permissions {

		// Skip owner entries as they are managed by Keeper and are not tracked in Terraform state.
		// Will not add them to the share map.
		if isOwnerPermission(p) {
			continue
		}
		principal := sharePrincipal(p)
		if principal == "" || strings.TrimSpace(p.Role) == "" {
			continue
		}
		elements[principal] = types.StringValue(p.Role)
	}
	if len(elements) == 0 {
		m.Share = types.MapNull(types.StringType)
		return nil
	}
	mv, diags := types.MapValue(types.StringType, elements)
	if diags.HasError() {
		return fmt.Errorf("unable to build share map from API response: %s", diags)
	}
	m.Share = mv
	return nil
}

// isOwnerPermission reports whether the entry represents the Keeper owner.
// NSF folders use role "owner"; NSF/classic records set the owner boolean
// (often with a non-owner role such as full-manager).
func isOwnerPermission(p UserPermissionEntry) bool {
	return p.Owner || strings.EqualFold(p.Role, RoleOwner)
}

// sharePrincipal returns the map key for a permission entry. NSF folders use
// accessor; NSF records use username. Accessor wins when both are set.
func sharePrincipal(p UserPermissionEntry) string {
	if accessor := strings.TrimSpace(p.Accessor); accessor != "" {
		return accessor
	}
	return strings.TrimSpace(p.Username)
}

// mapToStringMap converts a types.Map of StringType into a Go map[string]string.
// Null/unknown values and unknown string elements are ignored; null string
// values are treated as empty strings.
func mapToStringMap(m types.Map) map[string]string {
	out := map[string]string{}
	if m.IsNull() || m.IsUnknown() {
		return out
	}
	for k, v := range m.Elements() {
		s, ok := v.(types.String)
		if !ok || s.IsUnknown() {
			continue
		}
		if s.IsNull() {
			out[k] = ""
			continue
		}
		out[k] = s.ValueString()
	}
	return out
}

// quote wraps s in single quotes so spaces and CLI-special characters are
// passed through unchanged. Any single quote inside s is doubled, matching
// normalizeCommandForShell in the api package.
func quote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}

// escapeDoubleQuotes escapes embedded double quotes for use inside a
// double-quoted shell argument.
func escapeDoubleQuotes(s string) string {
	return strings.ReplaceAll(s, `"`, `\"`)
}
