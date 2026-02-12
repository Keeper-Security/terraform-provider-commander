# Look up an enterprise team by name or ID. Returns id, name, users, roles.
# Optional: managed_company (MSP only) to scope the lookup.

data "commander_enterprise_team" "example" {
  team = "Backend Developers"

  # Optional, MSP only
  # managed_company = "Acme Corp"
}

output "team_id" {
  value = data.commander_enterprise_team.example.id
}

output "team_name" {
  value = data.commander_enterprise_team.example.name
}

output "team_users" {
  value = data.commander_enterprise_team.example.users
}

output "team_roles" {
  value = data.commander_enterprise_team.example.roles
}
