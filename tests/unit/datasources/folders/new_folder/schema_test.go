// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newfolderds_test

import (
	"testing"

	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func TestNewFolderDataSource_Schema_HasExpectedAttributes(t *testing.T) {
	sch, _ := getDSSchema(t)

	for _, attr := range []string{"new_folder", "id", "name", "share"} {
		if sch.Attributes[attr] == nil {
			t.Errorf("expected %q attribute on data source schema", attr)
		}
	}
}

func TestNewFolderDataSource_Schema_NewFolderIsRequiredString(t *testing.T) {
	sch, _ := getDSSchema(t)

	a, ok := sch.Attributes["new_folder"].(dschema.StringAttribute)
	if !ok {
		t.Fatalf("new_folder is not a StringAttribute, got %T", sch.Attributes["new_folder"])
	}
	if !a.Required {
		t.Error("new_folder should be Required")
	}
	if a.Optional || a.Computed {
		t.Errorf("new_folder should not be Optional/Computed; got Optional=%v Computed=%v", a.Optional, a.Computed)
	}
}

func TestNewFolderDataSource_Schema_IdNameShareAreComputed(t *testing.T) {
	sch, _ := getDSSchema(t)

	cases := []struct {
		name string
		attr dschema.Attribute
	}{
		{"id", sch.Attributes["id"]},
		{"name", sch.Attributes["name"]},
		{"share", sch.Attributes["share"]},
	}
	for _, tc := range cases {
		if !tc.attr.IsComputed() {
			t.Errorf("%s should be Computed", tc.name)
		}
		if tc.attr.IsRequired() {
			t.Errorf("%s should not be Required on a data source", tc.name)
		}
		if tc.attr.IsOptional() {
			t.Errorf("%s should not be Optional on a data source", tc.name)
		}
	}
}
