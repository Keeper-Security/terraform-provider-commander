// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package classic_share_test

import (
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func TestResourceShareAttribute_Shape(t *testing.T) {
	attrs := classic_share.ResourceShareAttribute()
	got, ok := attrs[classic_share.AttrShare]
	if !ok {
		t.Fatalf("expected key %q in resource attribute map", classic_share.AttrShare)
	}
	mapAttr, ok := got.(schema.MapNestedAttribute)
	if !ok {
		t.Fatalf("expected schema.MapNestedAttribute, got %T", got)
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
	if len(mapAttr.Validators) != 2 {
		t.Errorf("expected 2 validators on share (non-empty, key email), got %d", len(mapAttr.Validators))
	}
	nested := mapAttr.NestedObject.Attributes
	if _, ok := nested[classic_share.AttrCanShare]; !ok {
		t.Errorf("expected nested attribute %q on share", classic_share.AttrCanShare)
	}
	if _, ok := nested[classic_share.AttrCanEdit]; !ok {
		t.Errorf("expected nested attribute %q on share", classic_share.AttrCanEdit)
	}
}

func TestDataSourceShareAttribute_Shape(t *testing.T) {
	attrs := classic_share.DataSourceShareAttribute()
	got, ok := attrs[classic_share.AttrShare]
	if !ok {
		t.Fatalf("expected key %q in data source attribute map", classic_share.AttrShare)
	}
	mapAttr, ok := got.(dschema.MapNestedAttribute)
	if !ok {
		t.Fatalf("expected dschema.MapNestedAttribute, got %T", got)
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
