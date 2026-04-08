// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package enterprisenode_test

import (
	"regexp"
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
					"isolated":    false,
				},
			}
		case strings.Contains(cmd, "enterprise-info") && strings.Contains(cmd, "-n") && strings.Contains(cmd, "--format json"):
			return "ok", []map[string]interface{}{
				{"node_id": float64(0), "name": "Root", "parent_node": "", "parent_id": float64(0), "isolated": false},
				{"node_id": float64(1001), "name": "acc-test-node", "parent_node": "Root", "parent_id": float64(0), "isolated": false},
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

// TestAccEnterpriseNodeResource_withToggleIsolated creates a node then updates toggle_isolated to true
// (toggle_isolated is not supported on create).
func TestAccEnterpriseNodeResource_withToggleIsolated(t *testing.T) {
	mock := &helpers.CommandServer{}
	updateSeen := false
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "enterprise-node") && strings.Contains(cmd, "'1002'") && strings.Contains(cmd, "--toggle-isolated") && !strings.Contains(cmd, "--delete") {
			updateSeen = true
		}
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
					"isolated":    updateSeen,
				},
			}
		case strings.Contains(cmd, "enterprise-info") && strings.Contains(cmd, "-n") && strings.Contains(cmd, "--format json"):
			return "ok", []map[string]interface{}{
				{"node_id": float64(0), "name": "Root", "parent_node": "", "parent_id": float64(0), "isolated": false},
				{"node_id": float64(1002), "name": "isolated-node", "parent_node": "Root", "parent_id": float64(0), "isolated": updateSeen},
			}
		case strings.Contains(cmd, "enterprise-node") && strings.Contains(cmd, "'1002'") && !strings.Contains(cmd, "--delete"):
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
  name   = "isolated-node"
  parent = "Root"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("commander_enterprise_node.test", "id", "1002"),
					resource.TestCheckResourceAttr("commander_enterprise_node.test", "name", "isolated-node"),
					resource.TestCheckResourceAttr("commander_enterprise_node.test", "parent", "Root"),
					resource.TestCheckResourceAttr("commander_enterprise_node.test", "toggle_isolated", "false"),
				),
			},
			{
				Config: provider.AccProviderConfig(server.URL, "test-key") + `
resource "commander_enterprise_node" "test" {
  name            = "isolated-node"
  parent          = "Root"
  toggle_isolated = true
}
`,
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

// TestAccEnterpriseNodeResource_CreateToggleIsolatedTrueFails verifies that create fails when toggle_isolated = true.
func TestAccEnterpriseNodeResource_CreateToggleIsolatedTrueFails(t *testing.T) {
	server := helpers.StartCommandServer(&helpers.CommandServer{}, nil)
	defer server.Close()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: provider.AccProviderConfig(server.URL, "test-key") + `
resource "commander_enterprise_node" "test" {
  name            = "isolated-node"
  parent          = "Root"
  toggle_isolated = true
}
`,
				ExpectError: regexp.MustCompile(`toggle_isolated is not supported in create`),
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
					"isolated":    false,
				},
			}
		case strings.Contains(cmd, "enterprise-info") && strings.Contains(cmd, "-n") && strings.Contains(cmd, "--format json"):
			readCallCount++
			name := "update-test-node"
			if readCallCount > 1 {
				name = "updated-node-name"
			}
			return "ok", []map[string]interface{}{
				{"node_id": float64(0), "name": "Root", "parent_node": "", "parent_id": float64(0), "isolated": false},
				{"node_id": float64(1003), "name": name, "parent_node": "Root", "parent_id": float64(0), "isolated": false},
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
