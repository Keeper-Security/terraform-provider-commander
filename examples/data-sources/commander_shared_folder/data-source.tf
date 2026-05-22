# commander_shared_folder (data source) — read an existing classic shared folder from Commander.
#
# Required:
#   shared_folder — Shared folder UID or vault path/name (passed to Commander: get '<value>' --format json).
#
# Computed:
#   id                    — Canonical classic shared folder UID from the API (same as resource id after create).
#   name                  — Vault path of the folder
#   user_permissions      — Default manage_users / manage_records
#   record_permissions    — Default can_share / can_edit
#   records               — Map record_uid => { can_share, can_edit }
#   users                 — Map email/uid => { manage_users, manage_records, expiration }

# Look up by vault path (use the same path you use in commander_shared_folder.name, or any path Commander accepts).
data "commander_shared_folder" "example" {
  shared_folder = "My Project Vault"
}

# Alternative: look up by UID from Commander or from commander_shared_folder.<name>.id
# data "commander_shared_folder" "by_uid" {
#   shared_folder = "BTbjhOmqw9iYal3OQJ9UAQ"
# }

# data "commander_shared_folder" "from_resource" {
#   shared_folder = commander_shared_folder.example.id
# }

output "folder_uid" {
  description = "Computed UID after read."
  value       = data.commander_shared_folder.example.id
}

output "folder_path" {
  value = data.commander_shared_folder.example.name
}

output "folder_default_user_permissions" {
  value = data.commander_shared_folder.example.user_permissions
}

output "folder_default_record_permissions" {
  value = data.commander_shared_folder.example.record_permissions
}

output "folder_record_overrides" {
  value = data.commander_shared_folder.example.records
}

output "folder_user_overrides" {
  value = data.commander_shared_folder.example.users
}
