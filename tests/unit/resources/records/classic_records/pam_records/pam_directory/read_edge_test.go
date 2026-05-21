// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamdirectory_test

import (
	"context"
	"strings"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestRead_InvalidResponseData(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "get") {
			return "ok", "this-is-not-a-valid-record-map"
		}
		return "ok", nil
	}
	server := startMockServer(mock, responseForCommand)
	defer server.Close()

	r := newConfiguredResource(t, server)
	sch, objType := getSchema(t)

	hostVals := newHostnameValues("host.com", nil)
	rawState := tftypes.NewValue(objType, newPlanStateValues(
		"uid-abc", "Title", hostVals, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	))

	req := resource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: rawState}}
	resp := resource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: rawState}}
	r.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when response data is invalid")
	}
}
