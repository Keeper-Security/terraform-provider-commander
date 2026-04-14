// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamrecordresources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// AudioBitDepthValidator validates RBI audio bit depth (8 or 16).
type AudioBitDepthValidator struct{}

func (v AudioBitDepthValidator) Description(_ context.Context) string {
	return "audio bit depth must be 8 for 8-bit or 16 for 16-bit"
}

func (v AudioBitDepthValidator) MarkdownDescription(_ context.Context) string {
	return "audio bit depth must be `8` for **8-bit** or `16` for **16-bit**"
}

func (v AudioBitDepthValidator) ValidateInt64(_ context.Context, req validator.Int64Request, resp *validator.Int64Response) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	val := req.ConfigValue.ValueInt64()
	if val != 8 && val != 16 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid audio_bit_depth value",
			fmt.Sprintf("expected 8 for 8-bit or 16 for 16-bit, got %d", val),
		)
	}
}

// AudioChannelsValidator validates RBI audio channel count (1 or 2).
type AudioChannelsValidator struct{}

func (v AudioChannelsValidator) Description(_ context.Context) string {
	return "audio channels must be `1` for **mono** or `2` for **stereo**"
}

func (v AudioChannelsValidator) MarkdownDescription(_ context.Context) string {
	return "audio channels must be `1` for **mono** or `2` for **stereo**"
}

func (v AudioChannelsValidator) ValidateInt32(_ context.Context, req validator.Int32Request, resp *validator.Int32Response) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	val := req.ConfigValue.ValueInt32()
	if val != 1 && val != 2 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid audio_channels value",
			fmt.Sprintf("expected 1 for mono or 2 for stereo, got %d", val),
		)
	}
}
