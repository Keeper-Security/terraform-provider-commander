// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

// Package pamrecordresources holds shared Terraform schema fragments for PAM vault record resources.
package pamrecordresources

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// RecordIDAttribute is the standard computed record UID with UseStateForUnknown.
func RecordIDAttribute(description, markdownDescription string) schema.Attribute {
	return schema.StringAttribute{
		Computed:            true,
		Description:         description,
		MarkdownDescription: markdownDescription,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
}

// RecordTitleAttribute is the standard required PAM record title.
func RecordTitleAttribute(description, markdownDescription string) schema.Attribute {
	return schema.StringAttribute{
		Required:            true,
		Description:         description,
		MarkdownDescription: markdownDescription,
		Validators: []validator.String{
			utils.StringMinLengthValidator("Title", 1, false),
		},
	}
}

// RecordNotesAttribute is optional notes shared by PAM record types.
func RecordNotesAttribute() schema.Attribute {
	return schema.StringAttribute{
		Optional:            true,
		Description:         "Optional notes for this configuration.",
		MarkdownDescription: "Optional **notes** for this configuration.",
		Validators: []validator.String{
			utils.StringMinLengthValidator("Notes", 0, true),
		},
	}
}

// RecordFolderAttribute is optional vault folder for storing the PAM record.
func RecordFolderAttribute(description, markdownDescription string) schema.Attribute {
	return schema.StringAttribute{
		Optional:            true,
		Description:         description,
		MarkdownDescription: markdownDescription,
		Validators: []validator.String{
			utils.StringMinLengthValidator("Folder", 1, true),
		},
	}
}

// RbiRecordURLAttribute is the required RBI target URL field (rbiUrl) for remote browser records.
func RbiRecordURLAttribute(description, markdownDescription string) schema.Attribute {
	return schema.StringAttribute{
		Required:            true,
		Description:         description,
		MarkdownDescription: markdownDescription,
		Validators: []validator.String{
			utils.StringMinLengthValidator("URL", 1, false),
		},
	}
}

// PamRemoteBrowserRBISettingsAttributes returns inner attributes for pam_remote_browser_settings / pam rbi edit.
func PamRemoteBrowserRBISettingsAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"configuration": schema.StringAttribute{
			Required:            true,
			Description:         "PAM Configuration UID for remote browser settings.",
			MarkdownDescription: "**PAM Configuration UID** for remote browser settings.",
			Validators: []validator.String{
				utils.StringMinLengthValidator("PAM Configuration UID", 1, false),
			},
		},
		"remote_browser_isolation": schema.BoolAttribute{
			Optional:            true,
			Description:         "Enable remote browser isolation.",
			MarkdownDescription: "Enable **remote browser isolation**.",
		},
		"connections_recording": schema.BoolAttribute{
			Optional:            true,
			Description:         "Manage graphical session recording.",
			MarkdownDescription: "**Manage graphical session recording**.",
		},
		"key_events": schema.BoolAttribute{
			Optional:            true,
			Description:         "Manage key events for session recording.",
			MarkdownDescription: "**Manage key events for session recording**.",
		},
		"allow_url_navigation": schema.BoolAttribute{
			Optional:            true,
			Description:         "Allow navigation via direct URL manipulation.",
			MarkdownDescription: "Allow **navigation** via direct URL manipulation.",
		},
		"ignore_server_cert": schema.BoolAttribute{
			Optional:            true,
			Description:         "Ignore Server Certificate.",
			MarkdownDescription: "**Ignore Server Certificate**.",
		},
		"allowed_urls": schema.SetAttribute{
			Optional:            true,
			ElementType:         types.StringType,
			Description:         "Allowed URL patterns. When set, must contain at least one non-empty value. Omit the attribute entirely to clear via the CLI (empty argument); an empty list is not allowed.",
			MarkdownDescription: "**Allowed URL patterns.** When set, must contain at least one non-empty value. **Omit** the attribute to clear via the CLI; an empty list is not allowed.",
			Validators: []validator.Set{
				utils.SetNotEmptyValidator("Allowed URLs"),
				utils.SetNoEmptyStringsValidator("Allowed URLs"),
			},
		},
		"allowed_resource_urls": schema.SetAttribute{
			Optional:            true,
			ElementType:         types.StringType,
			Description:         "Allowed resource URL patterns. When set, must contain at least one non-empty value. Omit the attribute entirely to clear via the CLI (empty argument); an empty list is not allowed.",
			MarkdownDescription: "**Allowed resource URL patterns.** When set, must contain at least one non-empty value. **Omit** the attribute to clear via the CLI; an empty list is not allowed.",
			Validators: []validator.Set{
				utils.SetNotEmptyValidator("Allowed Resource URLs"),
				utils.SetNoEmptyStringsValidator("Allowed Resource URLs"),
			},
		},
		"auto_fill_targets": schema.SetAttribute{
			Optional:            true,
			ElementType:         types.StringType,
			Description:         "Browser autofill targets. When set, must contain at least one non-empty value. Omit the attribute entirely to clear via the CLI (empty argument); an empty list is not allowed.",
			MarkdownDescription: "**Browser autofill targets.** When set, must contain at least one non-empty value. **Omit** the attribute to clear via the CLI; an empty list is not allowed.",
			Validators: []validator.Set{
				utils.SetNotEmptyValidator("Auto-fill Targets"),
				utils.SetNoEmptyStringsValidator("Auto-fill Targets"),
			},
		},
		"auto_fill_credentials": schema.StringAttribute{
			Optional:            true,
			Description:         "Record UID of Credentials stored in given pam configuration.",
			MarkdownDescription: "Record UID of **Credentials** stored in given pam configuration.",
			Validators: []validator.String{
				utils.StringMinLengthValidator("Auto-fill Credentials", 1, true),
			},
		},
		"allow_copy": schema.BoolAttribute{
			Optional:            true,
			Description:         "Can copy to clipboard.",
			MarkdownDescription: "Can **copy** to clipboard.",
		},
		"allow_paste": schema.BoolAttribute{
			Optional:            true,
			Description:         "Can paste from clipboard.",
			MarkdownDescription: "Can **paste** from clipboard.",
		},
		"disable_audio": schema.BoolAttribute{
			Optional:            true,
			Description:         "Disable audio.",
			MarkdownDescription: "**Disable audio**.",
		},
		"audio_channels": schema.Int32Attribute{
			Computed:            true,
			Optional:            true,
			Description:         "Number of audio channels; must be 1 for mono or 2 for stereo.",
			MarkdownDescription: "Number of **audio channels**; must be `1` for **mono** or `2` for **stereo**.",
			Validators: []validator.Int32{
				AudioChannelsValidator{},
			},
			Default: int32default.StaticInt32(2),
		},
		"audio_bit_depth": schema.Int64Attribute{
			Computed:            true,
			Optional:            true,
			Description:         "Audio bit depth; must be 8 for 8-bit or 16 for 16-bit.",
			MarkdownDescription: "Audio **bit depth**; must be `8` for **8-bit** or `16` for **16-bit**.",
			Validators: []validator.Int64{
				AudioBitDepthValidator{},
			},
			Default: int64default.StaticInt64(16),
		},
		"audio_sample_rate": schema.Int64Attribute{
			Computed:            true,
			Optional:            true,
			Description:         "Audio sample rate in Hz (for example 48000).",
			MarkdownDescription: "Audio **sample rate** in Hz (for example `48000`).",
			Default:             int64default.StaticInt64(44100),
		},
	}
}
