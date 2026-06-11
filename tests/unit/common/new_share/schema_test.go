// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package new_share_test

import (
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func TestResourceShareAttribute_Shape(t *testing.T) {
	attrs := new_share.ResourceShareAttribute()
	got, ok := attrs[new_share.AttrShare]
	if !ok {
		t.Fatalf("expected key %q in resource attribute map", new_share.AttrShare)
	}
	mapAttr, ok := got.(schema.MapAttribute)
	if !ok {
		t.Fatalf("expected schema.MapAttribute, got %T", got)
	}
	if !mapAttr.Optional {
		t.Error("expected share to be Optional")
	}
	if mapAttr.Computed {
		t.Error("expected share to NOT be Computed (Optional only)")
	}
	if mapAttr.Required {
		t.Error("expected share to NOT be Required")
	}
	if len(mapAttr.Validators) != 3 {
		t.Errorf("expected 3 validators on share (non-empty, key min-length, value enum), got %d", len(mapAttr.Validators))
	}
}

func TestDataSourceShareAttribute_Shape(t *testing.T) {
	attrs := new_share.DataSourceShareAttribute()
	got, ok := attrs[new_share.AttrShare]
	if !ok {
		t.Fatalf("expected key %q in data source attribute map", new_share.AttrShare)
	}
	mapAttr, ok := got.(dschema.MapAttribute)
	if !ok {
		t.Fatalf("expected dschema.MapAttribute, got %T", got)
	}
	if !mapAttr.Computed {
		t.Error("expected data source share to be Computed")
	}
	if mapAttr.Optional {
		t.Error("expected data source share to NOT be Optional")
	}
	if mapAttr.Required {
		t.Error("expected data source share to NOT be Required")
	}
}
