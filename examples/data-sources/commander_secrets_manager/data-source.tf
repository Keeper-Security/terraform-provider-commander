# Look up a Secrets Manager application by name or UID. Returns id, name, shares, app_users.

data "commander_secrets_manager" "example" {
  application = "Production Secrets App"
}

output "app_id" {
  value = data.commander_secrets_manager.example.id
}

output "app_name" {
  value = data.commander_secrets_manager.example.name
}

output "app_shares" {
  value = data.commander_secrets_manager.example.shares
}

output "app_users" {
  value = data.commander_secrets_manager.example.app_users
}
