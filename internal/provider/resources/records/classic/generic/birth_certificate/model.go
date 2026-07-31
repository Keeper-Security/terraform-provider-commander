// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package birthcertificate

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonrecordbirthcertificate "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/birth_certificate"
)

// BirthCertificateResourceModel is the classic Birth Certificate resource state model:
// shared Birth Certificate fields plus the `share` attribute reconciled via classic_share.
type BirthCertificateResourceModel struct {
	commonrecordbirthcertificate.BirthCertificateModel
	classic_share.ShareModel
}
