// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package wifi

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func (d *WifiDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dschema.Schema{
		Description:         SchemaDescription,
		MarkdownDescription: SchemaMarkdownDescription,
		Attributes: map[string]dschema.Attribute{
			"record_uid": dschema.StringAttribute{
				Required:            true,
				Description:         "Vault record UID of the wifiCredentials record to read.",
				MarkdownDescription: "Vault **record UID** of the `wifiCredentials` record to read.",
			},
			"id": dschema.StringAttribute{
				Computed:            true,
				Description:         "Same as record_uid from the vault.",
				MarkdownDescription: "Same as **record_uid** from the vault.",
			},
			"title": dschema.StringAttribute{
				Computed:            true,
				Description:         "Title of the WiFi credentials record.",
				MarkdownDescription: "**Title** of the WiFi credentials record.",
			},
			"folder": dschema.StringAttribute{
				Computed:            true,
				Description:         "Folder UID or path for the record.",
				MarkdownDescription: "**Folder** UID or path for the record.",
			},
			"notes": dschema.StringAttribute{
				Computed:            true,
				Description:         "Notes on the record, if any.",
				MarkdownDescription: "**Notes** on the record, if any.",
			},
			"ssid": dschema.StringAttribute{
				Computed:            true,
				Description:         "WiFi network SSID (network name).",
				MarkdownDescription: "WiFi network **SSID** (network name).",
			},
			"password": dschema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				Description:         "Password for the WiFi network.",
				MarkdownDescription: "**Password** for the WiFi network.",
			},
			"encryption": dschema.StringAttribute{
				Computed:            true,
				Description:         "Encryption type. One of: wep, wpa, noEncryption.",
				MarkdownDescription: "Encryption type. One of: `wep`, `wpa`, `noEncryption`.",
			},
			"is_ssid_hidden": dschema.BoolAttribute{
				Computed:            true,
				Description:         "Whether the SSID is hidden (not broadcast).",
				MarkdownDescription: "Whether the SSID is **hidden** (not broadcast).",
			},
			"custom": dschema.ListNestedAttribute{
				Computed:            true,
				Description:         "Custom fields stored in the record's `custom` array.",
				MarkdownDescription: "Custom fields stored in the record's `custom` array.",
				NestedObject: dschema.NestedAttributeObject{
					Attributes: map[string]dschema.Attribute{
						"type": dschema.StringAttribute{
							Computed:            true,
							Description:         "Keeper field type (e.g. text, email, secret).",
							MarkdownDescription: "Keeper field **type** (e.g. `text`, `email`, `secret`).",
						},
						"label": dschema.StringAttribute{
							Computed:            true,
							Description:         "Field label.",
							MarkdownDescription: "Field **label**.",
						},
						"value": dschema.StringAttribute{
							Computed:            true,
							Sensitive:           true,
							Description:         "Field value (JSON-encoded for complex types).",
							MarkdownDescription: "Field **value** (JSON-encoded for complex types).",
						},
						"sensitive": dschema.BoolAttribute{
							Computed:            true,
							Description:         "Whether the value should be treated as sensitive.",
							MarkdownDescription: "Whether the value should be treated as **sensitive**.",
						},
					},
				},
			},
			"share": dschema.MapNestedAttribute{
				Computed:            true,
				Description:         "Users this record is shared with. Keys are email addresses; values report each user's permissions.",
				MarkdownDescription: "Users this record is shared with. Map keys are **email addresses**; values report each user's `can_share` and `can_edit` permissions.",
				NestedObject: dschema.NestedAttributeObject{
					Attributes: map[string]dschema.Attribute{
						"can_share": dschema.BoolAttribute{
							Computed:            true,
							Description:         "Whether the user can re-share the record with others.",
							MarkdownDescription: "Whether the user can **re-share** the record with others.",
						},
						"can_edit": dschema.BoolAttribute{
							Computed:            true,
							Description:         "Whether the user can edit the record.",
							MarkdownDescription: "Whether the user can **edit** the record.",
						},
					},
				},
			},
		},
	}
}
