// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamuser

import (
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// PamUserSharedModel is the shared resource model for PAM User records, used
// by both classic and new PAM User resources. It excludes share-extension
// fields, which classic and new resources add separately.
type PamUserSharedModel struct {
	commonrecordsutils.BaseVaultRecordModel

	Login                types.String             `tfsdk:"login"`
	Password             types.String             `tfsdk:"password"`
	DistinguishedName    types.String             `tfsdk:"distinguished_name"`
	PrivatePEMKey        types.String             `tfsdk:"private_pem_key"`
	PublicKey            types.String             `tfsdk:"public_key"`
	PrivateKeyPassphrase types.String             `tfsdk:"private_key_passphrase"`
	ConnectDatabase      types.String             `tfsdk:"connect_database"`
	Managed              types.Bool               `tfsdk:"managed"`
	RotationSettings     *PamUserRotationSettings `tfsdk:"rotation_settings"`
}

// PamUserRotationSettings models the nested `rotation_settings` block applied
// via `pam rotation edit` and read back via `pam rotation info`.
type PamUserRotationSettings struct {
	RotationProfile types.String `tfsdk:"rotation_profile"`
	Configuration   types.String `tfsdk:"configuration"`
	IamAadConfig    types.String `tfsdk:"iam_aad_config"`
	Resource        types.String `tfsdk:"resource"`
	Enabled         types.Bool   `tfsdk:"enabled"`
	ScheduleCron    types.String `tfsdk:"schedule_cron"`
	ScheduleJSON    types.String `tfsdk:"schedule_json"`
	OnDemand        types.Bool   `tfsdk:"on_demand"`
	ScheduleConfig  types.Bool   `tfsdk:"schedule_config"`
	Complexity      types.String `tfsdk:"complexity"`
}
