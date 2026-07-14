// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package classic_share

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// UserPermissionEntry aliases the canonical type in utils. The helper code
// in this package reads only the classic shape (Username / Shareable /
// Editable); NSF-shape entries with empty Username are filtered out by
// MapResponseToModel.
//
// Example JSON:
//
//	{ "username": "user@example.com", "shareable": true, "editable": false }
type UserPermissionEntry = utils.UserPermissionEntry

// ShareResponseFragment is a tiny embeddable struct exposing only the
// user_permissions field for classic vault record responses. Response types
// that don't already include user_permissions can compose this struct to
// inherit the share parsing.
type ShareResponseFragment struct {
	UserPermissions []UserPermissionEntry `json:"user_permissions"`
}

// SharePermissionsObjectType returns the AttrTypes used for the share map's
// nested value object (`can_share`, `can_edit`). Exposed for callers that
// need to construct an empty/null map (e.g. types.MapNull(...)) or build
// ObjectValue instances directly.
func SharePermissionsObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		AttrCanShare: types.BoolType,
		AttrCanEdit:  types.BoolType,
	}
}

// ShareEntryAttrType is the element type of the `share` map: an object with
// can_share and can_edit booleans. Exposed for callers that need to
// construct an empty map (e.g. types.MapNull(classic_share.ShareEntryAttrType)).
var ShareEntryAttrType attr.Type = types.ObjectType{
	AttrTypes: SharePermissionsObjectType(),
}
