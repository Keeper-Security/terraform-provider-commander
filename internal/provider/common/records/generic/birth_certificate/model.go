// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

// Package birthcertificate holds shared model and helpers for the birthCertificate record type.
package birthcertificate

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// BirthCertificateModel maps a Keeper `birthCertificate` vault record.
// Shared between the resource and data source.
type BirthCertificateModel struct {
	utils.BaseVaultRecordModel

	Name      *utils.NameValue `tfsdk:"name"`
	BirthDate types.String     `tfsdk:"birth_date"`

	Custom []utils.CustomFieldModel `tfsdk:"custom"`
}
