// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package managedcompany_test

import (
	"strings"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/tests/acceptance/provider"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccManagedCompanyDataSource_basic(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if cmd == "msp-down" {
			return "ok", nil
		}
		if strings.Contains(cmd, "msp-info") && strings.Contains(cmd, "-m") && strings.Contains(cmd, "--format json") {
			return "ok", []map[string]interface{}{
				{
					"company_id":   float64(100),
					"company_name": "Acme Corp",
					"node":         "node-1",
					"plan":         "business",
					"storage":      "1TB",
				},
			}
		}
		return "ok", nil
	}
	server := helpers.StartCommandServer(mock, responseForCommand)
	defer server.Close()

	config := provider.AccProviderConfig(server.URL, "test-key") + `
data "commander_managed_company" "test" {
  managed_company = "Acme Corp"
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.commander_managed_company.test", "id", "100"),
					resource.TestCheckResourceAttr("data.commander_managed_company.test", "name", "Acme Corp"),
					resource.TestCheckResourceAttr("data.commander_managed_company.test", "managed_company", "Acme Corp"),
					resource.TestCheckResourceAttr("data.commander_managed_company.test", "node", "node-1"),
					resource.TestCheckResourceAttr("data.commander_managed_company.test", "plan", "business"),
					resource.TestCheckResourceAttr("data.commander_managed_company.test", "file_plan", "1tb"),
				),
			},
		},
	})
}

func TestAccManagedCompanyDataSource_byId(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		if cmd == "msp-down" {
			return "ok", nil
		}
		if strings.Contains(cmd, "msp-info") && strings.Contains(cmd, "-m") && strings.Contains(cmd, "--format json") {
			return "ok", []map[string]interface{}{
				{
					"company_id":   float64(200),
					"company_name": "Other Company",
					"node":         "node-2",
					"plan":         "enterprise",
					"storage":      "5TB",
				},
			}
		}
		return "ok", nil
	}
	server := helpers.StartCommandServer(mock, responseForCommand)
	defer server.Close()

	config := provider.AccProviderConfig(server.URL, "test-key") + `
data "commander_managed_company" "test" {
  managed_company = "200"
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.commander_managed_company.test", "id", "200"),
					resource.TestCheckResourceAttr("data.commander_managed_company.test", "name", "Other Company"),
					resource.TestCheckResourceAttr("data.commander_managed_company.test", "managed_company", "200"),
					resource.TestCheckResourceAttr("data.commander_managed_company.test", "node", "node-2"),
					resource.TestCheckResourceAttr("data.commander_managed_company.test", "plan", "enterprise"),
					resource.TestCheckResourceAttr("data.commander_managed_company.test", "file_plan", "5tb"),
				),
			},
		},
	})
}
