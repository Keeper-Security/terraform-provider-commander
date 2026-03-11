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

// Update copies plan to state (no-op for this one-time resource; replace handles changes).

func TestEnterpriseScimPushResource_Update_Success(t *testing.T) {
	r := enterprisescimpush.NewEnterpriseScimPushResource().(*enterprisescimpush.EnterpriseScimPushResource)
	sch, objType := getSchema(t)
	rawPlan := tftypes.NewValue(objType, newPlanStateValues("id-1", "scim-1", "google", "record-uid-1", false))
	emptyState := tftypes.NewValue(objType, newPlanStateValues(nil, nil, nil, nil, nil))

	req := resource.UpdateRequest{
		Plan:  tfsdk.Plan{Schema: sch, Raw: rawPlan},
		State: tfsdk.State{Schema: sch, Raw: emptyState},
	}
	resp := resource.UpdateResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	r.Update(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Update failed: %v", resp.Diagnostics)
	}
}
