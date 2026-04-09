// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package enterpriseuser_test

import (
	"strings"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/tests/acceptance/provider"
	"github.com/Keeper-Security/terraform-provider-commander/tests/helpers"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEnterpriseUserResource_basic(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		switch {
		case strings.Contains(cmd, "enterprise-user") && strings.Contains(cmd, "--add") && !strings.Contains(cmd, "--delete"):
			return "ok", nil
		case strings.Contains(cmd, "enterprise-info") && strings.Contains(cmd, "-u") && strings.Contains(cmd, "--format json"):
			return "ok", []map[string]interface{}{
				{
					"user_id":   float64(10001),
					"email":     "testuser@example.com",
					"name":      "",
					"job_title": "",
					"node":      "Root",
					"status":    "Invited",
					"roles":     []interface{}{},
					"teams":     []interface{}{},
				},
			}
		case strings.Contains(cmd, "enterprise-info") && strings.Contains(cmd, "-n") && strings.Contains(cmd, "--format json"):
			return "ok", []map[string]interface{}{
				{
					"node_id":     float64(1),
					"name":        "Root",
					"parent_node": "",
					"parent_id":   float64(0),
				},
			}
		case strings.Contains(cmd, "enterprise-user") && strings.Contains(cmd, "--delete"):
			return "ok", nil
		}
		return "ok", nil
	}
	server := helpers.StartCommandServer(mock, responseForCommand)
	defer server.Close()

	config := provider.AccProviderConfig(server.URL, "test-key") + `
resource "commander_enterprise_user" "test" {
  email = "testuser@example.com"
  node  = "Root"
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("commander_enterprise_user.test", "email", "testuser@example.com"),
					resource.TestCheckResourceAttr("commander_enterprise_user.test", "node", "Root"),
					resource.TestCheckResourceAttr("commander_enterprise_user.test", "status", "Invited"),
					resource.TestCheckResourceAttrSet("commander_enterprise_user.test", "id"),
				),
			},
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("commander_enterprise_user.test", "id", "10001"),
					resource.TestCheckResourceAttr("commander_enterprise_user.test", "email", "testuser@example.com"),
					resource.TestCheckResourceAttr("commander_enterprise_user.test", "node", "Root"),
					resource.TestCheckResourceAttr("commander_enterprise_user.test", "status", "Invited"),
				),
			},
		},
	})
}

func TestAccEnterpriseUserResource_withNameAndJobTitle(t *testing.T) {
	mock := &helpers.CommandServer{}
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		switch {
		case strings.Contains(cmd, "enterprise-user") && strings.Contains(cmd, "--add") && !strings.Contains(cmd, "--delete"):
			return "ok", nil
		case strings.Contains(cmd, "enterprise-info") && strings.Contains(cmd, "-u") && strings.Contains(cmd, "--format json"):
			return "ok", []map[string]interface{}{
				{
					"user_id":   float64(10002),
					"email":     "dev@example.com",
					"name":      "Dev User",
					"job_title": "Developer",
					"node":      "Root",
					"status":    "Invited",
					"roles":     []interface{}{},
					"teams":     []interface{}{},
				},
			}
		case strings.Contains(cmd, "enterprise-info") && strings.Contains(cmd, "-n") && strings.Contains(cmd, "--format json"):
			return "ok", []map[string]interface{}{
				{
					"node_id":     float64(1),
					"name":        "Root",
					"parent_node": "",
					"parent_id":   float64(0),
				},
			}
		case strings.Contains(cmd, "enterprise-user") && strings.Contains(cmd, "--delete"):
			return "ok", nil
		}
		return "ok", nil
	}
	server := helpers.StartCommandServer(mock, responseForCommand)
	defer server.Close()

	config := provider.AccProviderConfig(server.URL, "test-key") + `
resource "commander_enterprise_user" "test" {
  email     = "dev@example.com"
  name      = "Dev User"
  job_title = "Developer"
  node      = "Root"
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("commander_enterprise_user.test", "email", "dev@example.com"),
					resource.TestCheckResourceAttr("commander_enterprise_user.test", "name", "Dev User"),
					resource.TestCheckResourceAttr("commander_enterprise_user.test", "job_title", "Developer"),
					resource.TestCheckResourceAttr("commander_enterprise_user.test", "node", "Root"),
					resource.TestCheckResourceAttr("commander_enterprise_user.test", "status", "Invited"),
				),
			},
		},
	})
}

func TestAccEnterpriseUserResource_updateAndDelete(t *testing.T) {
	mock := &helpers.CommandServer{}
	readCallCount := 0
	responseForCommand := func(cmd string, idx int) (string, interface{}) {
		switch {
		case strings.Contains(cmd, "enterprise-user") && strings.Contains(cmd, "--add") && !strings.Contains(cmd, "--delete"):
			return "ok", nil
		case strings.Contains(cmd, "enterprise-info") && strings.Contains(cmd, "-u") && strings.Contains(cmd, "--format json"):
			readCallCount++
			name, jobTitle := "", ""
			if readCallCount > 1 {
				name, jobTitle = "Updated Name", "Updated Title"
			}
			return "ok", []map[string]interface{}{
				{
					"user_id":   float64(10003),
					"email":     "update@example.com",
					"name":      name,
					"job_title": jobTitle,
					"node":      "Root",
					"status":    "Invited",
					"roles":     []interface{}{},
					"teams":     []interface{}{},
				},
			}
		case strings.Contains(cmd, "enterprise-info") && strings.Contains(cmd, "-n") && strings.Contains(cmd, "--format json"):
			return "ok", []map[string]interface{}{
				{
					"node_id":     float64(1),
					"name":        "Root",
					"parent_node": "",
					"parent_id":   float64(0),
				},
			}
		case strings.Contains(cmd, "enterprise-user") && strings.Contains(cmd, "'10003'") && !strings.Contains(cmd, "--delete"):
			return "ok", nil
		case strings.Contains(cmd, "enterprise-user") && strings.Contains(cmd, "--delete"):
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
resource "commander_enterprise_user" "test" {
  email = "update@example.com"
  node  = "Root"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("commander_enterprise_user.test", "email", "update@example.com"),
					resource.TestCheckResourceAttr("commander_enterprise_user.test", "node", "Root"),
				),
			},
			{
				Config: provider.AccProviderConfig(server.URL, "test-key") + `
resource "commander_enterprise_user" "test" {
  email     = "update@example.com"
  name      = "Updated Name"
  job_title = "Updated Title"
  node      = "Root"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("commander_enterprise_user.test", "email", "update@example.com"),
					resource.TestCheckResourceAttr("commander_enterprise_user.test", "name", "Updated Name"),
					resource.TestCheckResourceAttr("commander_enterprise_user.test", "job_title", "Updated Title"),
					resource.TestCheckResourceAttr("commander_enterprise_user.test", "node", "Root"),
				),
			},
			{
				Config: provider.AccProviderConfig(server.URL, "test-key"),
			},
		},
	})
}
