// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package sharedfolder_test

import (
	"context"
	"testing"

	sharedfolder "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/folders/classic_folders/shared_folder"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func TestSharedFolderResource_Schema(t *testing.T) {
	r := sharedfolder.NewSharedFolderResource().(*sharedfolder.SharedFolderResource)
	req := resource.SchemaRequest{}
	var resp resource.SchemaResponse
	r.Schema(context.Background(), req, &resp)
	if resp.Schema.Attributes["id"] == nil {
		t.Error("expected id attribute")
	}
	if resp.Schema.Attributes["name"] == nil {
		t.Error("expected name attribute")
	}
	if resp.Schema.Attributes["folder_location"] != nil {
		t.Error("did not expect folder_location attribute")
	}

	up, ok := resp.Schema.Attributes["user_permissions"].(schema.SingleNestedAttribute)
	if !ok || up.ObjectDefaultValue() == nil {
		t.Error("expected user_permissions SingleNestedAttribute with object default")
	}
	if up.Attributes == nil {
		t.Fatal("expected user_permissions nested attributes")
	}
	if err := assertBoolDefault(t, up.Attributes["manage_users"], "user_permissions.manage_users"); err != nil {
		t.Error(err)
	}
	if err := assertBoolDefault(t, up.Attributes["manage_records"], "user_permissions.manage_records"); err != nil {
		t.Error(err)
	}

	rp, ok := resp.Schema.Attributes["record_permissions"].(schema.SingleNestedAttribute)
	if !ok || rp.ObjectDefaultValue() == nil {
		t.Error("expected record_permissions SingleNestedAttribute with object default")
	}
	if err := assertBoolDefault(t, rp.Attributes["can_share"], "record_permissions.can_share"); err != nil {
		t.Error(err)
	}
	if err := assertBoolDefault(t, rp.Attributes["can_edit"], "record_permissions.can_edit"); err != nil {
		t.Error(err)
	}

	if resp.Schema.Attributes["records"] == nil {
		t.Error("expected records attribute")
	}
	rec, ok := resp.Schema.Attributes["records"].(schema.MapNestedAttribute)
	if !ok {
		t.Fatal("expected records to be MapNestedAttribute")
	}
	if err := assertBoolDefault(t, rec.NestedObject.Attributes["can_share"], "records.can_share"); err != nil {
		t.Error(err)
	}
	if err := assertBoolDefault(t, rec.NestedObject.Attributes["can_edit"], "records.can_edit"); err != nil {
		t.Error(err)
	}

	users, ok := resp.Schema.Attributes["users"].(schema.MapNestedAttribute)
	if !ok {
		t.Fatal("expected users to be MapNestedAttribute")
	}
	if len(users.NestedObject.Validators) == 0 {
		t.Error("expected users nested object validators (manage_users vs expiration)")
	}
	userAttrs := users.NestedObject.Attributes
	if err := assertBoolDefault(t, userAttrs["manage_users"], "users.manage_users"); err != nil {
		t.Error(err)
	}
	if err := assertBoolDefault(t, userAttrs["manage_records"], "users.manage_records"); err != nil {
		t.Error(err)
	}
	expAttr, ok := userAttrs["expiration"].(schema.StringAttribute)
	if !ok {
		t.Fatal("expected users.expiration StringAttribute")
	}
	if expAttr.StringDefaultValue() == nil {
		t.Error("expected users.expiration to have a string default")
	}
	if len(expAttr.Validators) == 0 {
		t.Error("expected users.expiration validators")
	}
}

func assertBoolDefault(t *testing.T, attr schema.Attribute, label string) error {
	t.Helper()
	b, ok := attr.(schema.BoolAttribute)
	if !ok {
		return &schemaAssertError{label: label, msg: "expected BoolAttribute"}
	}
	if b.BoolDefaultValue() == nil {
		return &schemaAssertError{label: label, msg: "expected bool default"}
	}
	return nil
}

type schemaAssertError struct {
	label, msg string
}

func (e *schemaAssertError) Error() string {
	return e.label + ": " + e.msg
}
