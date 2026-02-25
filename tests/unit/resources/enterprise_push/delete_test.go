// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisepush_test

import (
	"context"
	"testing"

	enterprisepush "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/enterprise_push"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Delete is a no-op: no API call; resource is removed from state by Terraform.

func TestEnterprisePushResource_Delete_Success(t *testing.T) {
	r := enterprisepush.NewEnterprisePushResource().(*enterprisepush.EnterprisePushResource)
	sch, objType := getSchema(t)
	rawState := tftypes.NewValue(objType, newPlanStateValues("id123", "/path/to/file.json", "sha256hash", []interface{}{"user@example.com"}, nil))

	req := resource.DeleteRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	var resp resource.DeleteResponse
	r.Delete(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Delete failed: %v", resp.Diagnostics)
	}
}
