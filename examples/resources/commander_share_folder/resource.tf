# Share folder: name (required). Optional: folder_location, user_permissions, record_permissions, records, users.
# Defaults when omitted: user_permissions { manage_users = false, manage_records = false }, record_permissions { can_share = false, can_edit = false }.

resource "commander_share_folder" "example" {
  name            = "My Shared Folder"
  folder_location = "My Node/Subfolder"

  user_permissions = {
    manage_users   = true
    manage_records = false
  }

  record_permissions = {
    can_share = true
    can_edit  = true
  }

  records = {
    "record_uid_abc" = { can_share = true, can_edit = false }
    "record_uid_def" = { can_share = false, can_edit = true }
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
      expiration     = "90d"
    }
    "guest@example.com" = {
      manage_users   = false
      manage_records = false
      expiration     = "2026-12-31 23:59:59"
    }
  }
}
