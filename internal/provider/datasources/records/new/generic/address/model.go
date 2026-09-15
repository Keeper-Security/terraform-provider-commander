// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package address

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonrecordaddress "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/address"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AddressDataSourceModel adds a lookup key (`location`) to the shared address
// model. A distinct lookup key name is required because `address` is already
// used by the nested address attribute in AddressModel.
type AddressDataSourceModel struct {
	Location types.String `tfsdk:"location"`
	commonrecordaddress.AddressModel
	new_share.ShareModel
}
