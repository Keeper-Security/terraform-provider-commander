// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamdirectory_test

import (
	"context"
	"strings"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

func TestRead_DS_InvalidResponseData(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "get") {
			return "ok", "this-is-not-a-valid-record-map"
		}
		return "ok", nil
	}
	server := startDSMockServer(mock, responseForCommand)
	defer server.Close()

	d := newConfiguredDataSource(t, server)
	sch := getDSSchema(t)

	configRaw := newDSConfigRaw(t, sch, "uid-abc")
	emptyState := newDSEmptyState(t, sch)

	req := datasource.ReadRequest{Config: tfsdk.Config{Schema: sch, Raw: configRaw}}
	resp := datasource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: emptyState}}
	d.Read(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics when response data is invalid")
	}
}
