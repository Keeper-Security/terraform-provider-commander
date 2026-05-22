# commander_shared_folder — create and manage a Keeper classic shared folder via Commander.
#
# Required:
#   name — Full vault path to the classic shared folder (e.g. "My Folder" at vault root, or
#          "Templates/Team/Project Vault" where parent folders must exist). Updates:
#          same parent + rename leaf uses rndir; different parent uses mv.
#
# Computed:
#   id — Shared folder UID (stable; use for import and for data source lookups).
#
# Optional (defaults applied when blocks/attributes omitted):
#   user_permissions    — default manage_users / manage_records (both default false)
#   record_permissions    — default can_share / can_edit (both default false)
#   records               — map of record_uid => per-record can_share / can_edit
#   users                 — map of email or user UID => manage_*, expiration
#
# Expiration per user: "never" or absolute time as yyyy-MM-ddTHH:mm:ss.
# manage_users must be false when expiration is a datetime (not "never").

resource "commander_shared_folder" "example" {
  # Vault path (include parent path when not at root).
  name = "Shared Folders/Engineering/My Project Vault"

  user_permissions = {
    manage_users   = true
    manage_records = false
  }

  record_permissions = {
    can_share = true
    can_edit  = true
  }

  records = {
    "_replace_with_record_uid_1_" = {
      can_share = true
      can_edit  = false
    }
    "_replace_with_record_uid_2_" = {
      can_share = false
      can_edit  = true
    }
  }

  users = {
    "alice@example.com" = {
      manage_users   = true
      manage_records = false
      expiration     = "never"
    }
    "bob@example.com" = {
      manage_users   = false
      manage_records = true
      expiration     = "never"
    }
    "contractor@example.com" = {
      manage_users   = false
      manage_records = false
      expiration     = "2026-12-31T23:59:59"
    }
  }
}

# Minimal folder (defaults only; empty records/users maps).
# resource "commander_shared_folder" "minimal" {
#   name = "Shared Folders/Team Read-only"
# }

output "shared_folder_id" {
  description = "UID for lookups, import id, and Commander APIs."
  value       = commander_shared_folder.example.id
}

output "shared_folder_path" {
  value = commander_shared_folder.example.name
}

output "shared_folder_default_record_permissions" {
  value = commander_shared_folder.example.record_permissions
}
