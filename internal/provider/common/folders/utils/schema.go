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

// ResourceCommonFolderAttributes returns computed id, required name and optional
// folder_location attributes for folder resources.
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
		AttrFolderLocation: schema.StringAttribute{
			Optional:            true,
			Description:         DescFolderLocation,
			MarkdownDescription: DescFolderLocation,
			PlanModifiers: []planmodifier.String{
				FolderLocationSemanticEquality(),
			},
		},
	}
}

// DataSourceCommonFolderAttributes returns computed id, name and folder_location attributes for folder data sources.
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
		AttrFolderLocation: dschema.StringAttribute{
			Computed:            true,
			Description:         DescFolderLocation,
			MarkdownDescription: DescFolderLocation,
		},
	}
}
