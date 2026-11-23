// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamcloud

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// SharedAttributes returns the PAM Cloud resource attribute map shared
// between classic and new resources. Callers add any share-extension
// attribute separately.
func SharedAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:            true,
			Description:         IDDescription,
			MarkdownDescription: IDMarkdownDescription,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"title": schema.StringAttribute{
			Required:            true,
			Description:         TitleDescription,
			MarkdownDescription: TitleMarkdownDescription,
			Validators: []validator.String{
				utils.StringMinLengthValidator("Title", 1, false),
			},
		},
		"url": schema.StringAttribute{
			Optional:            true,
			Description:         URLDescription,
			MarkdownDescription: URLMarkdownDescription,
			Validators: []validator.String{
				utils.StringMinLengthValidator("URL", 1, true),
			},
		},
		"account_identifier": schema.StringAttribute{
			Optional:            true,
			Description:         AccountIdentifierDescription,
			MarkdownDescription: AccountIdentifierMarkdownDescription,
			Validators: []validator.String{
				utils.StringMinLengthValidator("Account Identifier", 1, true),
			},
		},
		"notes": schema.StringAttribute{
			Optional:            true,
			Description:         NotesDescription,
			MarkdownDescription: NotesMarkdownDescription,
			Validators: []validator.String{
				utils.StringMinLengthValidator("Notes", 0, true),
			},
		},
		"folder_location": schema.StringAttribute{
			Optional:            true,
			Description:         FolderDescription,
			MarkdownDescription: FolderMarkdownDescription,
			Validators: []validator.String{
				utils.StringMinLengthValidator("Folder Location", 1, true),
			},
		},
	}
}

// SharedDataSourceAttributes returns computed PAM Cloud data source attributes
// shared between classic and new data sources. Callers add the lookup key
// (e.g. pam_cloud) and share-extension attributes separately.
func SharedDataSourceAttributes() map[string]dschema.Attribute {
	return map[string]dschema.Attribute{
		"id": dschema.StringAttribute{
			Computed:            true,
			Description:         IDDescription,
			MarkdownDescription: IDMarkdownDescription,
		},
		"title": dschema.StringAttribute{
			Computed:            true,
			Description:         TitleDescription,
			MarkdownDescription: TitleMarkdownDescription,
		},
		"url": dschema.StringAttribute{
			Computed:            true,
			Description:         URLDescription,
			MarkdownDescription: URLMarkdownDescription,
		},
		"account_identifier": dschema.StringAttribute{
			Computed:            true,
			Description:         AccountIdentifierDescription,
			MarkdownDescription: AccountIdentifierMarkdownDescription,
		},
		"notes": dschema.StringAttribute{
			Computed:            true,
			Description:         NotesDescription,
			MarkdownDescription: NotesMarkdownDescription,
		},
		"folder_location": dschema.StringAttribute{
			Computed:            true,
			Description:         FolderDescription,
			MarkdownDescription: FolderMarkdownDescription,
		},
	}
}
