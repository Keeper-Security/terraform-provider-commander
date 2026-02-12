// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriseteam_test

import (
	"strings"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/tests/acceptance/provider"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEnterpriseTeamResource_basic(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		switch {
		case strings.Contains(cmd, "enterprise-team") && strings.Contains(cmd, "--add") && !strings.Contains(cmd, "--delete"):
			return "Team ID: team-3001", nil
		case strings.Contains(cmd, "enterprise-info") && strings.Contains(cmd, "-t") && strings.Contains(cmd, "--format json"):
			return "ok", []map[string]interface{}{
				{
					"team_uid":  "team-3001",
					"name":      "acc-test-team",
					"restricts": "",
					"node":      "Root",
					"users":     []interface{}{},
					"roles":     []interface{}{},
				},
			}
		case strings.Contains(cmd, "enterprise-info") && strings.Contains(cmd, "-n") && strings.Contains(cmd, "--format json"):
			return "ok", []map[string]interface{}{
				{"node_id": float64(0), "name": "Root", "parent_node": "", "parent_id": float64(0)},
				{"node_id": float64(1), "name": "acc-test-team-node", "parent_node": "Root", "parent_id": float64(0)},
			}
		case strings.Contains(cmd, "enterprise-team") && strings.Contains(cmd, "--delete"):
			return "ok", nil
		}
		return "ok", nil
	}
	server := helpers.StartCommandServer(mock, responseForCommand)
	defer server.Close()

	config := provider.AccProviderConfig(server.URL, "test-key") + `
resource "commander_enterprise_team" "test" {
  name = "acc-test-team"
  node = "Root"
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("commander_enterprise_team.test", "id", "team-3001"),
					resource.TestCheckResourceAttr("commander_enterprise_team.test", "name", "acc-test-team"),
					resource.TestCheckResourceAttr("commander_enterprise_team.test", "node", "Root"),
				),
			},
		},
	})
}

func TestAccEnterpriseTeamResource_updateAndDelete(t *testing.T) {
	mock := &helpers.CommandServer{}
	readCallCount := 0
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		switch {
		case strings.Contains(cmd, "enterprise-team") && strings.Contains(cmd, "--add") && !strings.Contains(cmd, "--delete"):
			return "Team ID: team-3002", nil
		case strings.Contains(cmd, "enterprise-info") && strings.Contains(cmd, "-t") && strings.Contains(cmd, "--format json"):
			readCallCount++
			name := "update-team"
			if readCallCount > 1 {
				name = "updated-team-name"
			}
			return "ok", []map[string]interface{}{
				{
					"team_uid":  "team-3002",
					"name":      name,
					"restricts": "",
					"node":      "Root",
					"users":     []interface{}{},
					"roles":     []interface{}{},
				},
			}
		case strings.Contains(cmd, "enterprise-info") && strings.Contains(cmd, "-n") && strings.Contains(cmd, "--format json"):
			return "ok", []map[string]interface{}{
				{"node_id": float64(0), "name": "Root", "parent_node": "", "parent_id": float64(0)},
				{"node_id": float64(1), "name": "updated-team-node", "parent_node": "Root", "parent_id": float64(0)},
			}
		case strings.Contains(cmd, "enterprise-team") && strings.Contains(cmd, "team-3002") && !strings.Contains(cmd, "--delete"):
			return "ok", nil
		case strings.Contains(cmd, "enterprise-team") && strings.Contains(cmd, "--delete"):
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
resource "commander_enterprise_team" "test" {
  name = "update-team"
  node = "Root"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("commander_enterprise_team.test", "id", "team-3002"),
					resource.TestCheckResourceAttr("commander_enterprise_team.test", "name", "update-team"),
					resource.TestCheckResourceAttr("commander_enterprise_team.test", "node", "Root"),
				),
			},
			{
				Config: provider.AccProviderConfig(server.URL, "test-key") + `
resource "commander_enterprise_team" "test" {
  name = "updated-team-name"
  node = "Root"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("commander_enterprise_team.test", "id", "team-3002"),
					resource.TestCheckResourceAttr("commander_enterprise_team.test", "name", "updated-team-name"),
					resource.TestCheckResourceAttr("commander_enterprise_team.test", "node", "Root"),
				),
			},
			{
				Config: provider.AccProviderConfig(server.URL, "test-key"),
			},
		},
	})
}
