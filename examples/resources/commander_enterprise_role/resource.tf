# Enterprise role: name and node (required). Optional: users, teams, managing_nodes, enforcement_policies, managed_company (MSP only).

resource "commander_enterprise_role" "example" {
  name = "Developer"
  node = "Engineering"

  users = ["alice@example.com", "bob@example.com"]
  teams = ["Backend Developers"]

  # Grant admin privileges on specific nodes (map key = node name/ID)
  managing_nodes = {
    "Engineering" = {
      privileges = ["manage_user", "manage_roles", "manage_teams"]
      cascade    = true
    }
  }

  # Enforcement policies (key = policy name, value = string e.g. "true" or "false")
  enforcement_policies = {
    "REQUIRE_TWO_FACTOR" = "true"
  }

  # Optional, MSP only
  # managed_company = "Acme Corp"
}
