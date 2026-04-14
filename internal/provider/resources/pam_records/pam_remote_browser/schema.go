// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamremotebrowser

import (
	"context"

	pamrecordresources "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/pam_records/resources"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func (r *PamRemoteBrowserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages Keeper PAM remote browser configuration.",
		MarkdownDescription: "Manages **Keeper PAM** remote browser configuration.",
		Attributes: map[string]schema.Attribute{
			"id": pamrecordresources.RecordIDAttribute(
				"Remote browser configuration identifier assigned by Keeper after creation.",
				"**Remote browser configuration identifier** assigned by Keeper after creation.",
			),
			"title": pamrecordresources.RecordTitleAttribute(
				"Title of the remote browser configuration.",
				"**Title** of the remote browser configuration.",
			),
			"url": pamrecordresources.RbiRecordURLAttribute(
				"Target URL for the remote browser session.",
				"**Target URL** for the remote browser session.",
			),
			"notes": pamrecordresources.RecordNotesAttribute(),
			"folder": pamrecordresources.RecordFolderAttribute(
				"Folder UID or name to store PAM remote browser record. If not provided, the record will be stored in the root path of vault.",
				"Folder UID or name to store PAM remote browser record. If not provided, the record will be stored in the root path of vault.",
			),
			"pam_remote_browser_settings": schema.SingleNestedAttribute{
				Optional:            true,
				Description:         "Session and isolation settings for the remote browser.",
				MarkdownDescription: "Session and **isolation settings** for the remote browser.",
				Attributes:          pamrecordresources.PamRemoteBrowserRBISettingsAttributes(),
			},
		},
	}
}
