// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package enterpriserole_test

import (
	"strings"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/tests/acceptance/provider"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEnterpriseRoleDataSource_basic(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "enterprise-info") && strings.Contains(cmd, "-r") && strings.Contains(cmd, "--format json") {
			return "ok", []map[string]interface{}{
				{
					"role_id": float64(123),
					"name":    "Admin",
					"users":   []interface{}{},
					"teams":   []interface{}{},
				},
			}
		}
		return "ok", nil
	}
	server := helpers.StartCommandServer(mock, responseForCommand)
	defer server.Close()

	config := provider.AccProviderConfig(server.URL, "test-key") + `
data "commander_enterprise_role" "test" {
  role = "Admin"
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.commander_enterprise_role.test", "id", "123"),
					resource.TestCheckResourceAttr("data.commander_enterprise_role.test", "name", "Admin"),
					resource.TestCheckResourceAttr("data.commander_enterprise_role.test", "role", "Admin"),
					resource.TestCheckResourceAttr("data.commander_enterprise_role.test", "users.#", "0"),
					resource.TestCheckResourceAttr("data.commander_enterprise_role.test", "teams.#", "0"),
				),
			},
		},
	})
}

func TestAccEnterpriseRoleDataSource_withUsersAndTeams(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "enterprise-info") && strings.Contains(cmd, "-r") && strings.Contains(cmd, "--format json") {
			return "ok", []map[string]interface{}{
				{
					"role_id": float64(456),
					"name":    "Developer",
					"users":   []interface{}{"user@example.com"},
					"teams":   []interface{}{"Engineering"},
				},
			}
		}
		return "ok", nil
	}
	server := helpers.StartCommandServer(mock, responseForCommand)
	defer server.Close()

	config := provider.AccProviderConfig(server.URL, "test-key") + `
data "commander_enterprise_role" "test" {
  role = "Developer"
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.commander_enterprise_role.test", "id", "456"),
					resource.TestCheckResourceAttr("data.commander_enterprise_role.test", "name", "Developer"),
					resource.TestCheckResourceAttr("data.commander_enterprise_role.test", "role", "Developer"),
					resource.TestCheckTypeSetElemAttr("data.commander_enterprise_role.test", "users.*", "user@example.com"),
					resource.TestCheckTypeSetElemAttr("data.commander_enterprise_role.test", "teams.*", "Engineering"),
				),
			},
		},
	})
}
