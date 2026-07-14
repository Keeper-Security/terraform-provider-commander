# commander_new_folder (data source) — read an existing Keeper Drive folder from Commander.

# Look up by folder name.
data "commander_new_folder" "example" {
  new_folder = "Engineering"
}

# Alternative: look up by UID (e.g. the one returned from commander_new_folder.<name>.id).
# data "commander_new_folder" "by_uid" {
#   new_folder = "E6laPVJ1T3-sWchJCRaWOg"
# }

# data "commander_new_folder" "from_resource" {
#   new_folder = commander_new_folder.example.id
# }

output "new_folder_uid" {
  description = "Keeper Drive folder UID."
  value       = data.commander_new_folder.example.id
}

output "new_folder_name" {
  value = data.commander_new_folder.example.name
}

output "new_folder_share" {
  description = "Map of User Email/UID or Team Name/UID => role."
  value       = data.commander_new_folder.example.share
}
