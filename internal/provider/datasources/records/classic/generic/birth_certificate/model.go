// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package birthcertificate

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonrecordbirthcertificate "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/birth_certificate"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// BirthCertificateDataSourceModel adds a lookup key (`birth_certificate`) to the shared model.
type BirthCertificateDataSourceModel struct {
	BirthCertificate types.String `tfsdk:"birth_certificate"`
	commonrecordbirthcertificate.BirthCertificateModel
	classic_share.ShareModel
}
