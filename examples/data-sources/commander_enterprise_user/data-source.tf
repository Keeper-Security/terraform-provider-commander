# Look up an enterprise user by email or ID. Returns id, name, email, job_title, roles, teams, status.
# Optional: managed_company (MSP only) to scope the lookup.

data "commander_enterprise_user" "example" {
  user = "alice@example.com"

  # Optional, MSP only
  # managed_company = "Acme Corp"
}

output "user_id" {
  value = data.commander_enterprise_user.example.id
}

output "user_name" {
  value = data.commander_enterprise_user.example.name
}

output "user_email" {
  value = data.commander_enterprise_user.example.email
}

output "user_job_title" {
  value = data.commander_enterprise_user.example.job_title
}

output "user_roles" {
  value = data.commander_enterprise_user.example.roles
}

output "user_teams" {
  value = data.commander_enterprise_user.example.teams
}
