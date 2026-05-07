// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package shared_folder

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// RecordEntryMapElemType is the object type for each entry in the records map.
var RecordEntryMapElemType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		AttrCanShare: types.BoolType,
		AttrCanEdit:  types.BoolType,
	},
}

// UserEntryMapElemType is the object type for each entry in the users map.
var UserEntryMapElemType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		AttrManageUsers:   types.BoolType,
		AttrManageRecords: types.BoolType,
		AttrExpiration:    types.StringType,
	},
}
