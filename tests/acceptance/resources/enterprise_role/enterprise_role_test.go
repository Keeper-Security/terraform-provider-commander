// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriserole_test

import (
	"strings"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/tests/acceptance/provider"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEnterpriseRoleResource_basic(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		switch {
		case strings.Contains(cmd, "enterprise-role") && strings.Contains(cmd, "--add") && !strings.Contains(cmd, "--delete"):
			return "Role ID : 2001", nil
		case strings.Contains(cmd, "enterprise-info") && strings.Contains(cmd, "-r") && strings.Contains(cmd, "--format json"):
			return "ok", []map[string]interface{}{
				{
					"role_id":                   float64(2001),
					"name":                      "acc-test-role",
					"node":                      "Root",
					"users":                     []interface{}{},
					"teams":                     []interface{}{},
					"managed_nodes_permissions": []interface{}{},
					"enforcements":              map[string]interface{}{},
				},
			}
		case strings.Contains(cmd, "enterprise-info") && strings.Contains(cmd, "-n") && strings.Contains(cmd, "--format json"):
			return "ok", []map[string]interface{}{
				{"node_id": float64(0), "name": "Root", "parent_node": "", "parent_id": float64(0)},
				{"node_id": float64(1), "name": "acc-test-role-node", "parent_node": "Root", "parent_id": float64(0)},
			}
		case strings.Contains(cmd, "enterprise-role") && strings.Contains(cmd, "--delete"):
			return "ok", nil
		}
		return "ok", nil
	}
	server := helpers.StartCommandServer(mock, responseForCommand)
	defer server.Close()

	config := provider.AccProviderConfig(server.URL, "test-key") + `
resource "commander_enterprise_role" "test" {
  name = "acc-test-role"
  node = "Root"
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("commander_enterprise_role.test", "id", "2001"),
					resource.TestCheckResourceAttr("commander_enterprise_role.test", "name", "acc-test-role"),
					resource.TestCheckResourceAttr("commander_enterprise_role.test", "node", "Root"),
				),
			},
		},
	})
}

func TestAccEnterpriseRoleResource_updateAndDelete(t *testing.T) {
	mock := &helpers.CommandServer{}
	readCallCount := 0
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		switch {
		case strings.Contains(cmd, "enterprise-role") && strings.Contains(cmd, "--add") && !strings.Contains(cmd, "--delete"):
			return "Role ID : 2002", nil
		case strings.Contains(cmd, "enterprise-info") && strings.Contains(cmd, "-r") && strings.Contains(cmd, "--format json"):
			readCallCount++
			name, nodeName := "update-role", "Root"
			if readCallCount > 1 {
				name, nodeName = "updated-role-name", "Root"
			}
			return "ok", []map[string]interface{}{
				{
					"role_id":                   float64(2002),
					"name":                      name,
					"node":                      nodeName,
					"users":                     []interface{}{},
					"teams":                     []interface{}{},
					"managed_nodes_permissions": []interface{}{},
					"enforcements":              map[string]interface{}{},
				},
			}
		case strings.Contains(cmd, "enterprise-info") && strings.Contains(cmd, "-n") && strings.Contains(cmd, "--format json"):
			return "ok", []map[string]interface{}{
				{"node_id": float64(0), "name": "Root", "parent_node": "", "parent_id": float64(0)},
				{"node_id": float64(2002), "name": "updated-role-node", "parent_node": "Root", "parent_id": float64(0)},
			}
		case strings.Contains(cmd, "enterprise-role") && strings.Contains(cmd, "'2002'") && !strings.Contains(cmd, "--delete"):
			return "ok", nil
		case strings.Contains(cmd, "enterprise-role") && strings.Contains(cmd, "--delete"):
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
resource "commander_enterprise_role" "test" {
  name = "update-role"
  node = "Root"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("commander_enterprise_role.test", "id", "2002"),
					resource.TestCheckResourceAttr("commander_enterprise_role.test", "name", "update-role"),
					resource.TestCheckResourceAttr("commander_enterprise_role.test", "node", "Root"),
				),
			},
			{
				Config: provider.AccProviderConfig(server.URL, "test-key") + `
resource "commander_enterprise_role" "test" {
  name = "updated-role-name"
  node = "Root"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("commander_enterprise_role.test", "id", "2002"),
					resource.TestCheckResourceAttr("commander_enterprise_role.test", "name", "updated-role-name"),
					resource.TestCheckResourceAttr("commander_enterprise_role.test", "node", "Root"),
				),
			},
			{
				Config: provider.AccProviderConfig(server.URL, "test-key"),
			},
		},
	})
}
