// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package birthcertificate

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonrecordbirthcertificate "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/birth_certificate"
)

// BirthCertificateResourceModel is the new Birth Certificate resource state model:
// shared Birth Certificate fields plus the `share` attribute reconciled via new_share.
type BirthCertificateResourceModel struct {
	commonrecordbirthcertificate.BirthCertificateModel
	new_share.ShareModel
}
