// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package classic_share

import (
	providerutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// ResourceShareAttribute returns the `share` map-of-object attribute for
// resource schemas: Optional, MapNonEmptyValidator (rejects {} in config so
// MapResponseToModel can safely store null when the API filters to zero
// entries), MapKeysEmailValidator applied to keys. The attribute is
// intentionally NOT Computed; users express "no managed shares" by omitting
// the block entirely so the round-trip null == null stays clean.
func ResourceShareAttribute() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		AttrShare: schema.MapNestedAttribute{
			Optional:            true,
			Description:         DescShare,
			MarkdownDescription: DescShareMD,
			Validators: []validator.Map{
				providerutils.MapNonEmptyValidator(AttrShareValidatorLabel),
				providerutils.MapKeysEmailValidator(AttrShareValidatorLabel),
			},
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					AttrCanShare: schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
						Description:         DescCanShare,
						MarkdownDescription: DescCanShare,
					},
					AttrCanEdit: schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
						Description:         DescCanEdit,
						MarkdownDescription: DescCanEdit,
					},
				},
			},
		},
	}
}

// DataSourceShareAttribute returns the `share` map-of-object attribute for
// data source schemas: Computed only (read-only output populated from the
// API response).
func DataSourceShareAttribute() map[string]dschema.Attribute {
	return map[string]dschema.Attribute{
		AttrShare: dschema.MapNestedAttribute{
			Computed:            true,
			Description:         DescShare,
			MarkdownDescription: DescShareMD,
			NestedObject: dschema.NestedAttributeObject{
				Attributes: map[string]dschema.Attribute{
					AttrCanShare: dschema.BoolAttribute{
						Computed:            true,
						Description:         DescCanShare,
						MarkdownDescription: DescCanShare,
					},
					AttrCanEdit: dschema.BoolAttribute{
						Computed:            true,
						Description:         DescCanEdit,
						MarkdownDescription: DescCanEdit,
					},
				},
			},
		},
	}
}
