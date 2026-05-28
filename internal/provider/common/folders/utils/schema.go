// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package utils

import (
	providerutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// ResourceCommonFolderAttributes returns computed id and required name attributes for folder resources.
func ResourceCommonFolderAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		AttrId: schema.StringAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			Description:         DescId,
			MarkdownDescription: DescIdMD,
		},
		AttrName: schema.StringAttribute{
			Required:            true,
			Description:         DescName,
			MarkdownDescription: DescNameMD,
			Validators: []validator.String{
				providerutils.StringMinLengthValidator(NameValidatorLabel, 1, false),
			},
		},
	}
}

// DataSourceCommonFolderAttributes returns computed id and name attributes for folder data sources.
func DataSourceCommonFolderAttributes() map[string]dschema.Attribute {
	return map[string]dschema.Attribute{
		AttrId: dschema.StringAttribute{
			Computed:            true,
			Description:         DescId,
			MarkdownDescription: DescIdMD,
		},
		AttrName: dschema.StringAttribute{
			Computed:            true,
			Description:         DescName,
			MarkdownDescription: DescNameMD,
		},
	}
}

// MergeResourceAttributes combines resource attribute maps; later maps override earlier keys.
func MergeResourceAttributes(maps ...map[string]schema.Attribute) map[string]schema.Attribute {
	result := map[string]schema.Attribute{}
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

// MergeDataSourceAttributes combines data source attribute maps; later maps override earlier keys.
func MergeDataSourceAttributes(maps ...map[string]dschema.Attribute) map[string]dschema.Attribute {
	result := map[string]dschema.Attribute{}
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}
