// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package membership

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonrecordmembership "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/membership"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// MembershipDataSourceModel adds a lookup key (`membership`) to the shared
// membership model.
type MembershipDataSourceModel struct {
	Membership types.String `tfsdk:"membership"`
	commonrecordmembership.MembershipModel
	classic_share.ShareModel
}
