// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamremotebrowser

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (d *PamRemoteBrowserDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dschema.Schema{
		Description:         "Reads an existing Keeper PAM remote browser vault record by record UID.",
		MarkdownDescription: "Reads an existing **Keeper PAM** remote browser vault record by **record UID** (`get <uid> --format json`).",
		Attributes: map[string]dschema.Attribute{
			"record_uid": dschema.StringAttribute{
				Required:            true,
				Description:         "Vault record UID of the pamRemoteBrowser record to read.",
				MarkdownDescription: "Vault **record UID** of the `pamRemoteBrowser` record to read.",
			},
			"id": dschema.StringAttribute{
				Computed:            true,
				Description:         "Same as record_uid from the vault.",
				MarkdownDescription: "Same as **record_uid** from the vault.",
			},
			"title": dschema.StringAttribute{
				Computed:            true,
				Description:         "Title of the remote browser record.",
				MarkdownDescription: "**Title** of the remote browser record.",
			},
			"url": dschema.StringAttribute{
				Computed:            true,
				Description:         "Target URL for the remote browser session (rbiUrl).",
				MarkdownDescription: "**Target URL** for the remote browser session (`rbiUrl`).",
			},
			"notes": dschema.StringAttribute{
				Computed:            true,
				Description:         "Notes on the record, if any.",
				MarkdownDescription: "**Notes** on the record, if any.",
			},
			"folder": dschema.StringAttribute{
				Computed:            true,
				Description:         "Folder UID or path for the record, if returned by the API.",
				MarkdownDescription: "**Folder** UID or path for the record, if returned by the API.",
			},
			"pam_remote_browser_settings": dschema.SingleNestedAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "Session and isolation settings parsed from the vault record when present.",
				MarkdownDescription: "Session and **isolation settings** parsed from the vault record when present.",
				Attributes:          pamRemoteBrowserRBISettingsDataSourceAttributes(),
			},
		},
	}
}

func pamRemoteBrowserRBISettingsDataSourceAttributes() map[string]dschema.Attribute {
	return map[string]dschema.Attribute{
		"configuration": dschema.StringAttribute{
			Computed:            true,
			Description:         "PAM Configuration UID for remote browser settings (may be unknown from get JSON).",
			MarkdownDescription: "**PAM Configuration UID** for remote browser settings (may be unknown from `get` JSON).",
		},
		"remote_browser_isolation": dschema.BoolAttribute{
			Computed:            true,
			Description:         "Remote browser isolation (may be unknown from get JSON).",
			MarkdownDescription: "**Remote browser isolation** (may be unknown from `get` JSON).",
		},
		"connections_recording": dschema.BoolAttribute{
			Computed:            true,
			Description:         "Graphical session recording (may be unknown from get JSON).",
			MarkdownDescription: "**Graphical session recording** (may be unknown from `get` JSON).",
		},
		"key_events": dschema.BoolAttribute{
			Computed:            true,
			Description:         "Key events for session recording.",
			MarkdownDescription: "**Key events** for session recording.",
		},
		"allow_url_navigation": dschema.BoolAttribute{
			Computed:            true,
			Description:         "Allow navigation via direct URL manipulation.",
			MarkdownDescription: "Allow **navigation** via direct URL manipulation.",
		},
		"ignore_server_cert": dschema.BoolAttribute{
			Computed:            true,
			Description:         "Ignore server certificate.",
			MarkdownDescription: "**Ignore server certificate**.",
		},
		"allowed_urls": dschema.SetAttribute{
			Computed:            true,
			ElementType:         types.StringType,
			Description:         "Allowed URL patterns.",
			MarkdownDescription: "**Allowed URL patterns**.",
		},
		"allowed_resource_urls": dschema.SetAttribute{
			Computed:            true,
			ElementType:         types.StringType,
			Description:         "Allowed resource URL patterns.",
			MarkdownDescription: "**Allowed resource URL patterns**.",
		},
		"auto_fill_targets": dschema.SetAttribute{
			Computed:            true,
			ElementType:         types.StringType,
			Description:         "Browser autofill targets.",
			MarkdownDescription: "**Browser autofill targets**.",
		},
		"auto_fill_credentials": dschema.StringAttribute{
			Computed:            true,
			Description:         "Credentials record UID for autofill.",
			MarkdownDescription: "**Credentials** record UID for autofill.",
		},
		"allow_copy": dschema.BoolAttribute{
			Computed:            true,
			Description:         "Whether copy to clipboard is allowed.",
			MarkdownDescription: "Whether **copy** to clipboard is allowed.",
		},
		"allow_paste": dschema.BoolAttribute{
			Computed:            true,
			Description:         "Whether paste from clipboard is allowed.",
			MarkdownDescription: "Whether **paste** from clipboard is allowed.",
		},
		"disable_audio": dschema.BoolAttribute{
			Computed:            true,
			Description:         "Whether audio is disabled.",
			MarkdownDescription: "Whether **audio** is disabled.",
		},
		"audio_channels": dschema.Int32Attribute{
			Computed:            true,
			Description:         "Number of audio channels.",
			MarkdownDescription: "Number of **audio channels**.",
		},
		"audio_bit_depth": dschema.Int64Attribute{
			Computed:            true,
			Description:         "Audio bit depth.",
			MarkdownDescription: "Audio **bit depth**.",
		},
		"audio_sample_rate": dschema.Int64Attribute{
			Computed:            true,
			Description:         "Audio sample rate in Hz.",
			MarkdownDescription: "Audio **sample rate** in Hz.",
		},
	}
}
