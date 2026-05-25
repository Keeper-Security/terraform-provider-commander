# commander_new_folder — create and manage a nested shared folder (a.k.a. new folder) via Commander.
#
# Required:
#   name — Folder name. Must be non-empty. Renames are supported in place via
#          the Update path (folder UID stays the same).
#
# Computed:
#   id — Keeper Drive folder UID. Stable across renames; use it for import,
#        for data source lookups, and when referencing the folder from other
#        resources.
#
# Optional:
#   share — Map of share grants. Each key is a user email; each value is one
#           of the allowed roles below. Omit the block entirely if you do not
#           want Terraform to manage shares for this folder.
#
# Allowed share roles (one of):
#   - "viewer"                  Read-only access to the folder contents.
#   - "share-manager"           May re-share but not edit content.
#   - "content-manager"         May edit / add / remove records and subfolders.
#   - "content-share-manager"   content-manager + share-manager combined.
#   - "full-manager"            All of the above; effectively co-owner
#                               (the original owner remains the Keeper-side
#                               owner of the folder).
#
# Notes on `share`:
#   - The folder *owner* is managed by Keeper and is NOT represented in this
#     block. If you read the API response back, the owner entry is filtered
#     out automatically.
#   - On every plan, Terraform reconciles the declared map against the API:
#     adds new emails (grant), removes emails no longer present (revoke), and
#     re-applies grant for any role change.
#   - The block is fully optional. If you do not declare `share`, Terraform
#     will NOT touch existing shares on the folder (no drift surfaced).

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
# Omit the share block entirely; Terraform will not grant or revoke anything,
# even if shares exist on the folder server-side.
# -----------------------------------------------------------------------------
# resource "commander_new_folder" "untouched_shares" {
#   name = "Marketing / Templates"
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
