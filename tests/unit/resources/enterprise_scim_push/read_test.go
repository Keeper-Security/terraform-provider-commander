// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisescimpush_test

import (
	"context"
	"testing"

	enterprisescimpush "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/enterprise_scim_push"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Read is a no-op: state is copied to response unchanged.

func TestEnterpriseScimPushResource_Read_Success(t *testing.T) {
	r := enterprisescimpush.NewEnterpriseScimPushResource().(*enterprisescimpush.EnterpriseScimPushResource)
	sch, objType := getSchema(t)
	rawState := tftypes.NewValue(objType, newPlanStateValues("id-1", "scim-1", "google", "record-uid-1", true))

	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	resp := resource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: rawState}}
	r.Read(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", resp.Diagnostics)
	}
}
