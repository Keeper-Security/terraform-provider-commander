# commander_new_folder — create and manage a nested shared folder (a.k.a. new folder) via Commander.
# -----------------------------------------------------------------------------
# Example 1: minimal folder — just a name.
# -----------------------------------------------------------------------------
resource "commander_new_folder" "minimal" {
  name = "Engineering"
}

# -----------------------------------------------------------------------------
# Example 2: folder with a full share map covering every supported role.
# -----------------------------------------------------------------------------
resource "commander_new_folder" "engineering" {
  name = "Engineering / Platform"

  share = {
    # Read-only access for a contractor.
    "contractor@example.com" = "viewer"

    # Senior engineer who can re-share but cannot edit content.
    "lead-share@example.com" = "share-manager"

    # Team member who can add/remove records and subfolders.
    "platform-member@example.com" = "content-manager"

    # Team lead with both content and share authority.
    "platform-lead@example.com" = "content-share-manager"

    # Co-owner (full authority other than ownership transfer).
    "director@example.com" = "full-manager"
  }
}

# -----------------------------------------------------------------------------
# Example 3: rename in place.
# Update the `name` and apply — the folder UID (id) stays the same.
# -----------------------------------------------------------------------------
# resource "commander_new_folder" "engineering" {
#   name = "Engineering / Platform (renamed)"
#   # share = { ... } // unchanged
# }

# -----------------------------------------------------------------------------
# Example 4: leave existing shares untouched.
# Omit the `share` block and tell Terraform to ignore drift on it; Read will
# still refresh the value into state, but plans will not propose any changes.
# -----------------------------------------------------------------------------
# resource "commander_new_folder" "untouched_shares" {
#   name = "Marketing / Templates"
#
#   lifecycle {
#     ignore_changes = [share]
#   }
# }

# -----------------------------------------------------------------------------
# Outputs — useful for wiring this folder into other resources / modules.
# -----------------------------------------------------------------------------
output "engineering_folder_uid" {
  description = "Keeper Drive folder UID for the Engineering folder. Stable across renames."
  value       = commander_new_folder.engineering.id
}

output "engineering_folder_name" {
  value = commander_new_folder.engineering.name
}

output "engineering_folder_share_map" {
  description = "Currently-managed share grants (owner filtered out)."
  value       = commander_new_folder.engineering.share
}

# -----------------------------------------------------------------------------
# Importing an existing folder into Terraform state.
# See ./import.sh for the CLI form. Or use the import block (Terraform >= 1.5):
# -----------------------------------------------------------------------------
# import {
#   to = commander_new_folder.engineering
#   id = "E6laPVJ1T3-sWchJCRaWOg" # folder UID
# }
