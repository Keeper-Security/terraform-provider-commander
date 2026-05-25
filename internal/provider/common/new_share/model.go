// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package new_share

import "github.com/hashicorp/terraform-plugin-framework/types"

// ShareModel is an embeddable Terraform model fragment that exposes the
// `share` map attribute. Compose it into any resource/data source model that
// needs share permissions.
//
// Example:
//
//	type Model struct {
//	    folderutils.IdentityModel
//	    new_share.ShareModel
//	}
type ShareModel struct {
	Share types.Map `tfsdk:"share"`
}
