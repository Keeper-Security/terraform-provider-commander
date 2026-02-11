// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managecompany_test

import (
	"strings"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/tests/acceptance/provider"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccManageCompanyResource_basic(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		switch {
		case cmd == "msp-down":
			return "ok", nil
		case strings.Contains(cmd, "msp-add"):
			return "4001", nil
		case strings.Contains(cmd, "msp-info") && strings.Contains(cmd, "-m") && strings.Contains(cmd, "--format json"):
			return "ok", []map[string]interface{}{
				{
					"company_id":   float64(4001),
					"company_name": "acc-test-company",
					"node":         "Root",
					"node_name":    "Root",
					"plan":         "business",
					"storage":      "100gb",
					"allocated":    float64(10),
					"addons":       []interface{}{},
				},
			}
		case strings.Contains(cmd, "msp-remove"):
			return "ok", nil
		}
		return "ok", nil
	}
	server := helpers.StartCommandServer(mock, responseForCommand)
	defer server.Close()

	config := provider.AccProviderConfig(server.URL, "test-key") + `
resource "commander_manage_company" "test" {
  name      = "acc-test-company"
  node      = "Root"
  plan      = "business"
  seats     = 10
  file_plan = "100gb"
  add_ons   = []
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("commander_manage_company.test", "id", "4001"),
					resource.TestCheckResourceAttr("commander_manage_company.test", "name", "acc-test-company"),
					resource.TestCheckResourceAttr("commander_manage_company.test", "node", "Root"),
					resource.TestCheckResourceAttr("commander_manage_company.test", "plan", "business"),
					resource.TestCheckResourceAttr("commander_manage_company.test", "seats", "10"),
					resource.TestCheckResourceAttr("commander_manage_company.test", "file_plan", "100gb"),
					resource.TestCheckResourceAttr("commander_manage_company.test", "add_ons.#", "0"),
				),
			},
		},
	})
}

func TestAccManageCompanyResource_updateAndDelete(t *testing.T) {
	mock := &helpers.CommandServer{}
	readCallCount := 0
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		switch {
		case cmd == "msp-down":
			return "ok", nil
		case strings.Contains(cmd, "msp-add"):
			return "4002", nil
		case strings.Contains(cmd, "msp-info") && strings.Contains(cmd, "-m") && strings.Contains(cmd, "--format json"):
			readCallCount++
			name, plan, storage := "update-company", "business", "1tb"
			if readCallCount > 1 {
				name, plan, storage = "updated-company-name", "enterprise", "1tb"
			}
			return "ok", []map[string]interface{}{
				{
					"company_id":   float64(4002),
					"company_name": name,
					"node":         "Root",
					"node_name":    "Root",
					"plan":         plan,
					"storage":      storage,
					"allocated":    float64(25),
					"addons":       []interface{}{},
				},
			}
		case strings.Contains(cmd, "msp-update"):
			return "ok", nil
		case strings.Contains(cmd, "msp-remove"):
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
resource "commander_manage_company" "test" {
  name      = "update-company"
  node      = "Root"
  plan      = "business"
  seats     = 25
  file_plan = "1tb"
  add_ons   = []
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("commander_manage_company.test", "id", "4002"),
					resource.TestCheckResourceAttr("commander_manage_company.test", "name", "update-company"),
					resource.TestCheckResourceAttr("commander_manage_company.test", "node", "Root"),
					resource.TestCheckResourceAttr("commander_manage_company.test", "plan", "business"),
				),
			},
			{
				Config: provider.AccProviderConfig(server.URL, "test-key") + `
resource "commander_manage_company" "test" {
  name      = "updated-company-name"
  node      = "Root"
  plan      = "enterprise"
  seats     = 25
  file_plan = "1tb"
  add_ons   = []
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("commander_manage_company.test", "id", "4002"),
					resource.TestCheckResourceAttr("commander_manage_company.test", "name", "updated-company-name"),
					resource.TestCheckResourceAttr("commander_manage_company.test", "node", "Root"),
					resource.TestCheckResourceAttr("commander_manage_company.test", "plan", "enterprise"),
				),
			},
			{
				Config: provider.AccProviderConfig(server.URL, "test-key"),
			},
		},
	})
}
