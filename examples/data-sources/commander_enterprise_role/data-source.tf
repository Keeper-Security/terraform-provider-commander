# Look up an enterprise role by name or ID. Returns id, name, users, teams.
# Optional: managed_company (MSP only) to scope the lookup.

data "commander_enterprise_role" "example" {
  role = "Developer"

  # Optional, MSP only
  # managed_company = "Acme Corp"
}

output "role_id" {
  value = data.commander_enterprise_role.example.id
}

output "role_name" {
  value = data.commander_enterprise_role.example.name
}

output "role_users" {
  value = data.commander_enterprise_role.example.users
}

output "role_teams" {
  value = data.commander_enterprise_role.example.teams
}
