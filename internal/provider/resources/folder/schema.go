// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package folder

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *FolderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         DescResource,
		MarkdownDescription: DescResource,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description:         DescId,
				MarkdownDescription: DescId,
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         DescName,
				MarkdownDescription: DescName,
				Validators: []validator.String{
					utils.StringMinLengthValidator("Folder Name", 1, false),
				},
			},
			"folder_location": schema.StringAttribute{
				Optional:            true,
				Description:         DescFolderLocation,
				MarkdownDescription: DescFolderLocation,
			},
			"color": schema.StringAttribute{
				Optional:            true,
				Description:         DescColor,
				MarkdownDescription: DescColor,
				Validators: []validator.String{
					ColorValidator(),
				},
			},
			"records": schema.SetAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				Description:         DescRecords,
				MarkdownDescription: DescRecords,
			},
		},
	}
}
