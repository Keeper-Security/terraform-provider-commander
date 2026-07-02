// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamremotebrowser_test

import (
	"context"
	"testing"

	commonpamremotebrowser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam/pam_remote_browser"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestAudioBitDepthValidator_Valid(t *testing.T) {
	ctx := context.Background()
	v := commonpamremotebrowser.AudioBitDepthValidator{}
	p := path.Root("audio_bit_depth")

	for _, val := range []int64{8, 16} {
		var resp validator.Int64Response
		v.ValidateInt64(ctx, validator.Int64Request{
			Path:        p,
			ConfigValue: types.Int64Value(val),
		}, &resp)
		if resp.Diagnostics.HasError() {
			t.Errorf("expected no error for bit depth %d, got: %v", val, resp.Diagnostics)
		}
	}
}

func TestAudioBitDepthValidator_Invalid(t *testing.T) {
	ctx := context.Background()
	v := commonpamremotebrowser.AudioBitDepthValidator{}
	p := path.Root("audio_bit_depth")

	for _, val := range []int64{0, 4, 12, 24, 32} {
		var resp validator.Int64Response
		v.ValidateInt64(ctx, validator.Int64Request{
			Path:        p,
			ConfigValue: types.Int64Value(val),
		}, &resp)
		if !resp.Diagnostics.HasError() {
			t.Errorf("expected error for bit depth %d", val)
		}
	}
}

func TestAudioBitDepthValidator_NullAndUnknown(t *testing.T) {
	ctx := context.Background()
	v := commonpamremotebrowser.AudioBitDepthValidator{}
	p := path.Root("audio_bit_depth")

	var resp validator.Int64Response
	v.ValidateInt64(ctx, validator.Int64Request{Path: p, ConfigValue: types.Int64Null()}, &resp)
	if resp.Diagnostics.HasError() {
		t.Errorf("null should be skipped: %v", resp.Diagnostics)
	}

	resp = validator.Int64Response{}
	v.ValidateInt64(ctx, validator.Int64Request{Path: p, ConfigValue: types.Int64Unknown()}, &resp)
	if resp.Diagnostics.HasError() {
		t.Errorf("unknown should be skipped: %v", resp.Diagnostics)
	}
}

func TestAudioBitDepthValidator_Descriptions(t *testing.T) {
	ctx := context.Background()
	v := commonpamremotebrowser.AudioBitDepthValidator{}
	if v.Description(ctx) == "" {
		t.Error("Description should not be empty")
	}
	if v.MarkdownDescription(ctx) == "" {
		t.Error("MarkdownDescription should not be empty")
	}
}

func TestAudioChannelsValidator_Valid(t *testing.T) {
	ctx := context.Background()
	v := commonpamremotebrowser.AudioChannelsValidator{}
	p := path.Root("audio_channels")

	for _, val := range []int32{1, 2} {
		var resp validator.Int32Response
		v.ValidateInt32(ctx, validator.Int32Request{
			Path:        p,
			ConfigValue: types.Int32Value(val),
		}, &resp)
		if resp.Diagnostics.HasError() {
			t.Errorf("expected no error for channels %d, got: %v", val, resp.Diagnostics)
		}
	}
}

func TestAudioChannelsValidator_Invalid(t *testing.T) {
	ctx := context.Background()
	v := commonpamremotebrowser.AudioChannelsValidator{}
	p := path.Root("audio_channels")

	for _, val := range []int32{0, 3, 5, 8} {
		var resp validator.Int32Response
		v.ValidateInt32(ctx, validator.Int32Request{
			Path:        p,
			ConfigValue: types.Int32Value(val),
		}, &resp)
		if !resp.Diagnostics.HasError() {
			t.Errorf("expected error for channels %d", val)
		}
	}
}

func TestAudioChannelsValidator_NullAndUnknown(t *testing.T) {
	ctx := context.Background()
	v := commonpamremotebrowser.AudioChannelsValidator{}
	p := path.Root("audio_channels")

	var resp validator.Int32Response
	v.ValidateInt32(ctx, validator.Int32Request{Path: p, ConfigValue: types.Int32Null()}, &resp)
	if resp.Diagnostics.HasError() {
		t.Errorf("null should be skipped: %v", resp.Diagnostics)
	}

	resp = validator.Int32Response{}
	v.ValidateInt32(ctx, validator.Int32Request{Path: p, ConfigValue: types.Int32Unknown()}, &resp)
	if resp.Diagnostics.HasError() {
		t.Errorf("unknown should be skipped: %v", resp.Diagnostics)
	}
}

func TestAudioChannelsValidator_Descriptions(t *testing.T) {
	ctx := context.Background()
	v := commonpamremotebrowser.AudioChannelsValidator{}
	if v.Description(ctx) == "" {
		t.Error("Description should not be empty")
	}
	if v.MarkdownDescription(ctx) == "" {
		t.Error("MarkdownDescription should not be empty")
	}
}
