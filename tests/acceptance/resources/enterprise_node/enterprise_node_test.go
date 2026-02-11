// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisenode_test

import (
	"strings"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/tests/acceptance/provider"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEnterpriseNodeResource_basic(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		switch {
		case strings.Contains(cmd, "enterprise-node") && strings.Contains(cmd, "--add"):
			return "Node ID: 1001", nil
		case strings.Contains(cmd, "enterprise-info") && strings.Contains(cmd, "-n") && strings.Contains(cmd, "--format json") && strings.Contains(cmd, "--node"):
			return "ok", []map[string]interface{}{
				{
					"node_id":     float64(1001),
					"name":        "acc-test-node",
					"parent_node": "Root",
					"parent_id":   float64(0),
				},
			}
		case strings.Contains(cmd, "enterprise-info") && strings.Contains(cmd, "-n") && strings.Contains(cmd, "--format json"):
			return "ok", []map[string]interface{}{
				{"node_id": float64(0), "name": "Root", "parent_node": "", "parent_id": float64(0)},
				{"node_id": float64(1001), "name": "acc-test-node", "parent_node": "Root", "parent_id": float64(0)},
			}
		case strings.Contains(cmd, "enterprise-node") && strings.Contains(cmd, "--delete"):
			return "ok", nil
		}
		return "ok", nil
	}
	server := helpers.StartCommandServer(mock, responseForCommand)
	defer server.Close()

	config := provider.AccProviderConfig(server.URL, "test-key") + `
resource "commander_enterprise_node" "test" {
  name   = "acc-test-node"
  parent = "Root"
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("commander_enterprise_node.test", "id", "1001"),
					resource.TestCheckResourceAttr("commander_enterprise_node.test", "name", "acc-test-node"),
					resource.TestCheckResourceAttr("commander_enterprise_node.test", "parent", "Root"),
				),
			},
		},
	})
}

func TestAccEnterpriseNodeResource_withToggleIsolated(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		switch {
		case strings.Contains(cmd, "enterprise-node") && strings.Contains(cmd, "--add"):
			return "Node ID: 1002", nil
		case strings.Contains(cmd, "enterprise-info") && strings.Contains(cmd, "-n") && strings.Contains(cmd, "--format json") && strings.Contains(cmd, "--node"):
			return "ok", []map[string]interface{}{
				{
					"node_id":     float64(1002),
					"name":        "isolated-node",
					"parent_node": "Root",
					"parent_id":   float64(0),
				},
			}
		case strings.Contains(cmd, "enterprise-info") && strings.Contains(cmd, "-n") && strings.Contains(cmd, "--format json"):
			return "ok", []map[string]interface{}{
				{"node_id": float64(0), "name": "Root", "parent_node": "", "parent_id": float64(0)},
				{"node_id": float64(1002), "name": "isolated-node", "parent_node": "Root", "parent_id": float64(0)},
			}
		case strings.Contains(cmd, "enterprise-node") && strings.Contains(cmd, "--delete"):
			return "ok", nil
		}
		return "ok", nil
	}
	server := helpers.StartCommandServer(mock, responseForCommand)
	defer server.Close()

	config := provider.AccProviderConfig(server.URL, "test-key") + `
resource "commander_enterprise_node" "test" {
  name            = "isolated-node"
  parent          = "Root"
  toggle_isolated = true
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("commander_enterprise_node.test", "id", "1002"),
					resource.TestCheckResourceAttr("commander_enterprise_node.test", "name", "isolated-node"),
					resource.TestCheckResourceAttr("commander_enterprise_node.test", "parent", "Root"),
					resource.TestCheckResourceAttr("commander_enterprise_node.test", "toggle_isolated", "true"),
				),
			},
		},
	})
}

func TestAccEnterpriseNodeResource_updateAndDelete(t *testing.T) {
	mock := &helpers.CommandServer{}
	readCallCount := 0
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		switch {
		case strings.Contains(cmd, "enterprise-node") && strings.Contains(cmd, "--add") && !strings.Contains(cmd, "--delete"):
			return "Node ID: 1003", nil
		case strings.Contains(cmd, "enterprise-info") && strings.Contains(cmd, "-n") && strings.Contains(cmd, "--format json") && strings.Contains(cmd, "--node"):
			readCallCount++
			name := "update-test-node"
			if readCallCount > 1 {
				name = "updated-node-name"
			}
			return "ok", []map[string]interface{}{
				{
					"node_id":     float64(1003),
					"name":        name,
					"parent_node": "Root",
					"parent_id":   float64(0),
				},
			}
		case strings.Contains(cmd, "enterprise-info") && strings.Contains(cmd, "-n") && strings.Contains(cmd, "--format json"):
			readCallCount++
			name := "update-test-node"
			if readCallCount > 1 {
				name = "updated-node-name"
			}
			return "ok", []map[string]interface{}{
				{"node_id": float64(0), "name": "Root", "parent_node": "", "parent_id": float64(0)},
				{"node_id": float64(1003), "name": name, "parent_node": "Root", "parent_id": float64(0)},
			}
		case strings.Contains(cmd, "enterprise-node") && strings.Contains(cmd, "'1003'") && !strings.Contains(cmd, "--delete"):
			return "ok", nil
		case strings.Contains(cmd, "enterprise-node") && strings.Contains(cmd, "--delete"):
			return "ok", nil
		}
		return "ok", nil
	}
	server := helpers.StartCommandServer(mock, responseForCommand)
	defer server.Close()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: provider.AccProviderConfig(server.URL, "test-key") + `
resource "commander_enterprise_node" "test" {
  name   = "update-test-node"
  parent = "Root"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("commander_enterprise_node.test", "id", "1003"),
					resource.TestCheckResourceAttr("commander_enterprise_node.test", "name", "update-test-node"),
					resource.TestCheckResourceAttr("commander_enterprise_node.test", "parent", "Root"),
				),
			},
			{
				Config: provider.AccProviderConfig(server.URL, "test-key") + `
resource "commander_enterprise_node" "test" {
  name   = "updated-node-name"
  parent = "Root"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("commander_enterprise_node.test", "id", "1003"),
					resource.TestCheckResourceAttr("commander_enterprise_node.test", "name", "updated-node-name"),
					resource.TestCheckResourceAttr("commander_enterprise_node.test", "parent", "Root"),
				),
			},
			{
				Config: provider.AccProviderConfig(server.URL, "test-key"),
			},
		},
	})
}
