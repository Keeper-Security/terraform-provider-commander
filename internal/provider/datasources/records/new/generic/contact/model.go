// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package contact

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonrecordcontact "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/contact"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ContactDataSourceModel adds a lookup key (`contact`) to the shared contact model.
type ContactDataSourceModel struct {
	Contact types.String `tfsdk:"contact"`
	commonrecordcontact.ContactModel
	new_share.ShareModel
}
