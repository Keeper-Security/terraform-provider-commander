// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package wifi

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func (r *WifiResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{
		Description:         SchemaDescription,
		MarkdownDescription: SchemaMarkdownDescription,
		Attributes: utils.MergeResourceAttributes(
			commonrecordsutils.BaseRecordAttributes(),
			map[string]schema.Attribute{
				"ssid": schema.StringAttribute{
					Required:            true,
					Description:         "WiFi network SSID (network name).",
					MarkdownDescription: "WiFi network SSID (network name). Maps to the record's `text` field with label `SSID`.",
					Validators: []validator.String{
						utils.StringMinLengthValidator("SSID", 1, false),
					},
				},
				"password": schema.StringAttribute{
					Optional:            true,
					Sensitive:           true,
					Description:         "Password for the WiFi network.",
					MarkdownDescription: "Password for the WiFi network. Maps to the record's `password` field.",
					Validators: []validator.String{
						utils.StringMinLengthValidator("Password", 1, true),
					},
				},
				"encryption": schema.StringAttribute{
					Optional:            true,
					Description:         "Encryption type. One of: wep, wpa, noEncryption.",
					MarkdownDescription: "Encryption type. One of: `wep`, `wpa`, `noEncryption`. Maps to the record's `wifiEncryption` field.",
					Validators: []validator.String{
						utils.StringOneOfValidator("Encryption", AllowedEncryptions, true),
					},
				},
				"is_ssid_hidden": schema.BoolAttribute{
					Optional:            true,
					Description:         "Whether the SSID is hidden (not broadcast).",
					MarkdownDescription: "Whether the SSID is hidden (not broadcast). Maps to the record's `isSSIDHidden` field.",
				},
			},
			classic_share.ResourceShareAttribute(),
		),
	}
}
