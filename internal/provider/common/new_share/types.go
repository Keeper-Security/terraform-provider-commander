// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package new_share

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ShareEntryAttrType is the element type of the `share` map: a plain string
// holding the role value. Exposed for callers that need to construct an empty
// map (e.g. types.MapNull(new_share.ShareEntryAttrType)).
var ShareEntryAttrType attr.Type = types.StringType

// UserPermissionEntry aliases the canonical type in utils so all callers
// share the same struct.
//
// Example JSON:
//
//	{ "accessor": "user@example.com", "role": "viewer" }
type UserPermissionEntry = utils.UserPermissionEntry

// ShareResponseFragment is a tiny embeddable struct exposing only the
// user_permissions field. Resources/data sources whose API response does not
// already include user_permissions can compose this struct into their
// response type to inherit the share parsing.
type ShareResponseFragment struct {
	UserPermissions []UserPermissionEntry `json:"user_permissions"`
	TeamPermissions []UserPermissionEntry `json:"team_permissions"`
}
