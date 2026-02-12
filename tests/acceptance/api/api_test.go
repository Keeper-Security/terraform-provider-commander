// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package api_test

import (
	"strings"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/tests/acceptance/provider"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccApi_ExecuteCommandFlow_DataSource verifies the API flow (POST executecommand-async
// -> poll GET /result/req-N) works when the provider runs a datasource read.
func TestAccApi_ExecuteCommandFlow_DataSource(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "enterprise-info") && strings.Contains(cmd, "-n") && strings.Contains(cmd, "--format json") && strings.Contains(cmd, "--node") {
			return "ok", []map[string]interface{}{
				{
					"node_id":     float64(5001),
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
					resource.TestCheckResourceAttr("data.commander_enterprise_node.test", "id", "5001"),
					resource.TestCheckResourceAttr("data.commander_enterprise_node.test", "name", "Root"),
					resource.TestCheckResourceAttr("data.commander_enterprise_node.test", "node", "Root"),
				),
			},
		},
	})

	// API acceptance: the mock server must have received at least one command (submit + poll)
	if n := mock.CommandCount(); n < 1 {
		t.Errorf("expected at least one API command (executecommand-async); got %d", n)
	}
}

// TestAccApi_ExecuteCommandFlow_ResourceCreate verifies the API flow works when the provider
// runs a resource create (submit + poll for create command, then submit + poll for read).
func TestAccApi_ExecuteCommandFlow_ResourceCreate(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		switch {
		case strings.Contains(cmd, "enterprise-node") && strings.Contains(cmd, "--add"):
			return "Node ID: 5002", nil
		case strings.Contains(cmd, "enterprise-info") && strings.Contains(cmd, "-n") && strings.Contains(cmd, "--format json") && strings.Contains(cmd, "--node"):
			return "ok", []map[string]interface{}{
				{
					"node_id":     float64(5002),
					"name":        "api-test-node",
					"parent_node": "Root",
					"parent_id":   float64(0),
				},
			}
		case strings.Contains(cmd, "enterprise-info") && strings.Contains(cmd, "-n") && strings.Contains(cmd, "--format json"):
			return "ok", []map[string]interface{}{
				{"node_id": float64(0), "name": "Root", "parent_node": "", "parent_id": float64(0)},
				{"node_id": float64(5002), "name": "api-test-node", "parent_node": "Root", "parent_id": float64(0)},
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
  name   = "api-test-node"
  parent = "Root"
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("commander_enterprise_node.test", "id", "5002"),
					resource.TestCheckResourceAttr("commander_enterprise_node.test", "name", "api-test-node"),
					resource.TestCheckResourceAttr("commander_enterprise_node.test", "parent", "Root"),
				),
			},
		},
	})

	// API acceptance: multiple commands (provider init, create submit+poll, read submit+poll)
	if n := mock.CommandCount(); n < 2 {
		t.Errorf("expected at least 2 API commands (e.g. create + read); got %d", n)
	}
}
