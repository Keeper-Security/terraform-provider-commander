// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisenode_test

import (
	"strings"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/tests/acceptance/provider"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEnterpriseNodeDataSource_basic(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "enterprise-info") && strings.Contains(cmd, "-n") && strings.Contains(cmd, "--format json") && strings.Contains(cmd, "--node") {
			return "ok", []map[string]interface{}{
				{
					"node_id":     float64(123),
					"name":        "Root",
					"parent_node": "",
					"parent_id":   float64(0),
				},
			}
		}
		return "ok", nil
	}
	server := helpers.StartCommandServer(mock, responseForCommand)
	defer server.Close()

	config := provider.AccProviderConfig(server.URL, "test-key") + `
data "commander_enterprise_node" "test" {
  node = "Root"
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.commander_enterprise_node.test", "id", "123"),
					resource.TestCheckResourceAttr("data.commander_enterprise_node.test", "name", "Root"),
					resource.TestCheckResourceAttr("data.commander_enterprise_node.test", "parent", ""),
					resource.TestCheckResourceAttr("data.commander_enterprise_node.test", "parent_id", "0"),
					resource.TestCheckResourceAttr("data.commander_enterprise_node.test", "node", "Root"),
				),
			},
		},
	})
}
