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

func TestAccEnterpriseTeamDataSource_basic(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "enterprise-info") && strings.Contains(cmd, "-t") && strings.Contains(cmd, "--format json") {
			return "ok", []map[string]interface{}{
				{
					"team_uid": "team-uid-123",
					"name":     "Engineering",
					"users":    []interface{}{},
					"roles":    []interface{}{},
				},
			}
		}
		return "ok", nil
	}
	server := helpers.StartCommandServer(mock, responseForCommand)
	defer server.Close()

	config := provider.AccProviderConfig(server.URL, "test-key") + `
data "commander_enterprise_team" "test" {
  team = "Engineering"
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.commander_enterprise_team.test", "id", "team-uid-123"),
					resource.TestCheckResourceAttr("data.commander_enterprise_team.test", "name", "Engineering"),
					resource.TestCheckResourceAttr("data.commander_enterprise_team.test", "team", "Engineering"),
					resource.TestCheckResourceAttr("data.commander_enterprise_team.test", "users.#", "0"),
					resource.TestCheckResourceAttr("data.commander_enterprise_team.test", "roles.#", "0"),
				),
			},
		},
	})
}

func TestAccEnterpriseTeamDataSource_withUsersAndRoles(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "enterprise-info") && strings.Contains(cmd, "-t") && strings.Contains(cmd, "--format json") {
			return "ok", []map[string]interface{}{
				{
					"team_uid": "team-uid-456",
					"name":     "DevTeam",
					"users":    []interface{}{"user@example.com"},
					"roles":    []interface{}{"Admin"},
				},
			}
		}
		return "ok", nil
	}
	server := helpers.StartCommandServer(mock, responseForCommand)
	defer server.Close()

	config := provider.AccProviderConfig(server.URL, "test-key") + `
data "commander_enterprise_team" "test" {
  team = "DevTeam"
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.commander_enterprise_team.test", "id", "team-uid-456"),
					resource.TestCheckResourceAttr("data.commander_enterprise_team.test", "name", "DevTeam"),
					resource.TestCheckResourceAttr("data.commander_enterprise_team.test", "team", "DevTeam"),
					resource.TestCheckTypeSetElemAttr("data.commander_enterprise_team.test", "users.*", "user@example.com"),
					resource.TestCheckTypeSetElemAttr("data.commander_enterprise_team.test", "roles.*", "Admin"),
				),
			},
		},
	})
}
