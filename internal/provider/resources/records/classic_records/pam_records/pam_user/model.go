// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamuser

import (
	commonpamuser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records/pam_user"
)

// PamUserResourceModel is the classic PAM User resource model, aliased to the
// shared model used by both classic and new PAM User resources.
type PamUserResourceModel = commonpamuser.PamUserSharedModel

// PamUserRotationSettings is aliased to the shared rotation settings model.
type PamUserRotationSettings = commonpamuser.PamUserRotationSettings
