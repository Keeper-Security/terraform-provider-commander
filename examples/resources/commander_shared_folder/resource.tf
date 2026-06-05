# commander_shared_folder — create and manage a Keeper classic shared folder via Commander.
#
# Required:
#   name — Classic shared folder leaf name (without parent path), e.g. "My Project Vault".
#          Updates: same parent + rename leaf uses rndir; different parent uses mv.
#
# Computed:
#   id — Shared folder UID (stable; use for import and for data source lookups).
#
# Optional (defaults applied when blocks/attributes omitted):
#   folder_location       — Parent folder path where the classic shared folder will be created
#                           (e.g. "Shared Folders/Engineering"). Leave empty or omit for vault root.
#   user_permissions      — default manage_users / manage_records (both default false)
#   record_permissions    — default can_share / can_edit (both default false)
#   records               — map of record_uid => per-record can_share / can_edit
#   users                 — map of email or user UID => manage_*, expiration
#
# Expiration per user: "never" or absolute time as yyyy-MM-ddTHH:mm:ss.
# manage_users must be false when expiration is a datetime (not "never").

resource "commander_shared_folder" "example" {
  # Leaf name only.
  name            = "My Project Vault"
  folder_location = "Shared Folders/Engineering"

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