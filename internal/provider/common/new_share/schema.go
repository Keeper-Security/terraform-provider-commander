// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package new_share

import (
	providerutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ResourceShareAttribute returns the `share` map attribute for resource
// schemas: Optional, validators applied to keys (non-empty user UID or
// email) and values (one of AllowedRoles). The attribute is intentionally
// NOT Computed; if the user omits the block, callers should also skip share
// reconciliation in Read so the state stays null.
func ResourceShareAttribute() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		AttrShare: schema.MapAttribute{
			Optional:            true,
			ElementType:         types.StringType,
			Description:         DescShare,
			MarkdownDescription: DescShareMD,
			Validators: []validator.Map{
				providerutils.MapKeysMinLengthValidator(AttrShareValidatorLabel, 1),
				providerutils.MapValuesStringOneOfValidator(AttrShareValueLabel, AllowedRoles),
			},
		},
	}
}

// DataSourceShareAttribute returns the `share` map attribute for data source
// schemas: Computed only (read-only output populated from the API response).
func DataSourceShareAttribute() map[string]dschema.Attribute {
	return map[string]dschema.Attribute{
		AttrShare: dschema.MapAttribute{
			Computed:            true,
			ElementType:         types.StringType,
			Description:         DescShare,
			MarkdownDescription: DescShareMD,
		},
	}
}
