# commander_non_shared_folder (data source) — read an existing vault folder from Commander.
#
# Required:
#   folder — Folder UID or vault path (passed to Commander: get '<value>' --format json).
#
# Computed:
#   id              — Folder UID from the API.
#   name            — Folder name (leaf name).
#   folder_location — Parent folder path. Empty if at vault root.
#   records         — Set of record UIDs linked to this folder.

# Look up by vault path
data "commander_non_shared_folder" "example" {
  folder = "My personal Record"
}

# Alternative: look up by UID
# data "commander_non_shared_folder" "by_uid" {
#   folder = "NJiANrRnbuvVEOgnqjiYaw"
# }

# Reference from a resource
# data "commander_non_shared_folder" "from_resource" {
#   folder = commander_non_shared_folder.example.id
# }

output "folder_uid" {
  description = "Folder UID."
  value       = data.commander_non_shared_folder.example.id
}

output "folder_name" {
  value = data.commander_non_shared_folder.example.name
}

output "folder_location" {
  value = data.commander_non_shared_folder.example.folder_location
}

output "folder_records" {
  value = data.commander_non_shared_folder.example.records
}