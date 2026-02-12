// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriseuser_test

import (
	"strings"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/tests/acceptance/provider"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEnterpriseUserDataSource_basic(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "enterprise-info") && strings.Contains(cmd, "-u") && strings.Contains(cmd, "--format json") {
			return "ok", []map[string]interface{}{
				{
					"user_id":   float64(1001),
					"name":      "Jane Doe",
					"email":     "jane@example.com",
					"job_title": "Engineer",
					"status":    "active",
					"roles":     []interface{}{},
					"teams":     []interface{}{},
				},
			}
		}
		return "ok", nil
	}
	server := helpers.StartCommandServer(mock, responseForCommand)
	defer server.Close()

	config := provider.AccProviderConfig(server.URL, "test-key") + `
data "commander_enterprise_user" "test" {
  user = "jane@example.com"
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.commander_enterprise_user.test", "id", "1001"),
					resource.TestCheckResourceAttr("data.commander_enterprise_user.test", "name", "Jane Doe"),
					resource.TestCheckResourceAttr("data.commander_enterprise_user.test", "email", "jane@example.com"),
					resource.TestCheckResourceAttr("data.commander_enterprise_user.test", "job_title", "Engineer"),
					resource.TestCheckResourceAttr("data.commander_enterprise_user.test", "status", "active"),
					resource.TestCheckResourceAttr("data.commander_enterprise_user.test", "user", "jane@example.com"),
					resource.TestCheckResourceAttr("data.commander_enterprise_user.test", "roles.#", "0"),
					resource.TestCheckResourceAttr("data.commander_enterprise_user.test", "teams.#", "0"),
				),
			},
		},
	})
}

func TestAccEnterpriseUserDataSource_withRolesAndTeams(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if strings.Contains(cmd, "enterprise-info") && strings.Contains(cmd, "-u") && strings.Contains(cmd, "--format json") {
			return "ok", []map[string]interface{}{
				{
					"user_id":   float64(1002),
					"name":      "John Smith",
					"email":     "john@example.com",
					"job_title": "Admin",
					"status":    "active",
					"roles":     []interface{}{"Admin"},
					"teams":     []interface{}{"Engineering"},
				},
			}
		}
		return "ok", nil
	}
	server := helpers.StartCommandServer(mock, responseForCommand)
	defer server.Close()

	config := provider.AccProviderConfig(server.URL, "test-key") + `
data "commander_enterprise_user" "test" {
  user = "john@example.com"
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.commander_enterprise_user.test", "id", "1002"),
					resource.TestCheckResourceAttr("data.commander_enterprise_user.test", "name", "John Smith"),
					resource.TestCheckResourceAttr("data.commander_enterprise_user.test", "email", "john@example.com"),
					resource.TestCheckResourceAttr("data.commander_enterprise_user.test", "job_title", "Admin"),
					resource.TestCheckResourceAttr("data.commander_enterprise_user.test", "status", "active"),
					resource.TestCheckResourceAttr("data.commander_enterprise_user.test", "user", "john@example.com"),
					resource.TestCheckTypeSetElemAttr("data.commander_enterprise_user.test", "roles.*", "Admin"),
					resource.TestCheckTypeSetElemAttr("data.commander_enterprise_user.test", "teams.*", "Engineering"),
				),
			},
		},
	})
}
