# commander_new_folder — create and manage a Nested Shared Folder (a.k.a. "new folder") via Commander.
#
# Required:
#   name — Folder leaf name (without parent path), e.g. "Platform".
#
# Computed:
#   id   — Folder UID (stable across renames; use for import and data source lookups).
#
# Optional:
#   folder_location — Parent vault path where the folder will be created (e.g.
#                     "Engineering/Platform"). Leave empty or omit for vault root.
#                     Whitespace around "/" separators is ignored (e.g. "A / B" == "A/B").
#                     Cannot be changed after creation; moving a Nested Shared
#                     Folder to a new parent is not supported by this resource.
#   share           — Map of User Email/UID or Team Name/UID => role granted on the folder.
#                     Allowed roles: viewer, share-manager, content-manager,
#                     content-share-manager, full-manager. Omit the attribute or use
#                     `lifecycle.ignore_changes = [share]` to leave existing
#                     out-of-band shares untouched.

# ----------------------------------------------------------------------------
# Usage 1 — Minimal folder at vault root (no shares)
# ----------------------------------------------------------------------------
resource "commander_new_folder" "minimal" {
  name = "Engineering"
}

# ----------------------------------------------------------------------------
# Usage 2 — Nested folder under an existing parent path
# ----------------------------------------------------------------------------
resource "commander_new_folder" "platform" {
  name            = "Platform"
  folder_location = "Engineering"
}

# ----------------------------------------------------------------------------
# Usage 3 — Nested folder with a full share map covering every supported role
# ----------------------------------------------------------------------------
resource "commander_new_folder" "engineering" {
  name            = "Platform Secrets"
  folder_location = "Engineering / Platform"

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

# ----------------------------------------------------------------------------
# Usage 4 — Rename in place
# Update the `name` and apply; the folder UID (id) and folder_location stay the
# same. Changing folder_location is rejected (move not supported).
# ----------------------------------------------------------------------------
# resource "commander_new_folder" "engineering" {
#   name            = "Platform Secrets (renamed)"
#   folder_location = "Engineering / Platform"
#   # share = { ... } // unchanged
# }

# ----------------------------------------------------------------------------
# Usage 5 — Leave existing shares untouched
# Omit the `share` attribute and tell Terraform to ignore drift on it; Read
# still refreshes the value into state, but plans will not propose any changes.
# ----------------------------------------------------------------------------
# resource "commander_new_folder" "untouched_shares" {
#   name            = "Templates"
#   folder_location = "Marketing"
#
#   lifecycle {
#     ignore_changes = [share]
#   }
# }

# ----------------------------------------------------------------------------
# Outputs — useful for wiring this folder into other resources / modules
# ----------------------------------------------------------------------------
output "engineering_folder_uid" {
  description = "Nested Shared Folder UID. Stable across renames; pass to data sources or downstream record resources."
  value       = commander_new_folder.engineering.id
}

output "engineering_folder_name" {
  value = commander_new_folder.engineering.name
}

output "engineering_folder_location" {
  description = "Parent vault path the folder was created under (null/empty for vault root)."
  value       = commander_new_folder.engineering.folder_location
}

output "engineering_folder_share_map" {
  description = "Currently-managed share grants (owner filtered out)."
  value       = commander_new_folder.engineering.share
}

# ----------------------------------------------------------------------------
# Importing an existing folder into Terraform state.
# See ./import.sh for the CLI form. Or use the import block (Terraform >= 1.5):
# ----------------------------------------------------------------------------
# import {
#   to = commander_new_folder.engineering
#   id = "E6laPVJ1T3-sWchJCRaWOg" # folder UID (preferred) or folder name
# }
