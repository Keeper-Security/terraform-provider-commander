// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamuser

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonpamuser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam/pam_user"
)

// PamUserResourceModel is the classic PAM User resource state model: the
// shared PAM User fields plus the `share` attribute reconciled via the
// classic_share package and the `share-record` Commander CLI.
type PamUserResourceModel struct {
	commonpamuser.PamUserSharedModel
	classic_share.ShareModel
}

// PamUserRotationSettings is aliased to the shared rotation settings model.
type PamUserRotationSettings = commonpamuser.PamUserRotationSettings
