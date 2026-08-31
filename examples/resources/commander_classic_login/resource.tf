###############################################################################
# Example 1 - Minimal: title and login (username) only.
###############################################################################

resource "commander_classic_login" "minimal" {
  title = "Example Login"
  login = "user@example.com"
}

###############################################################################
# Example 2 - Full: password, website, folder, notes, custom fields, and
# per-user sharing.
#
# `share` is a map keyed by email address. Each value sets:
#   can_share — let the user re-share the record with others
#   can_edit  — let the user edit the record
# Removing an email from `share` on a subsequent apply revokes that user's
# access automatically.
###############################################################################

resource "commander_classic_login" "full" {
  title           = "Acme Portal"
  login           = "jane.smith@acme.com"
  password        = "ExamplePassword123!"
  website_address = "https://portal.acme.com"
  folder_location = "_REPLACE_WITH_FOLDER_PATH_OR_UID_"
  notes           = "Corporate SSO login for the Acme customer portal."

  share = {
    "alice@example.com" = {
      can_share = true
      can_edit  = true
    }
    "bob@example.com" = {
      can_share = false
      can_edit  = true
    }
    "viewer@example.com" = {
      can_share = false
      can_edit  = false
    }
  }

  /* Example of custom fields
    - For complex types, use jsonencode(JSON) matching the Keeper field schema.
    - For more information, see the Keeper field schema documentation:
      https://docs.keeper.io/en/keeperpam/secrets-manager/about/field-record-types
  */
  custom = [
    {
      type  = "text"
      label = "environment"
      value = "production"
    },
    {
      type  = "text"
      label = "team"
      value = "Platform Engineering"
    },
    {
      type  = "url"
      label = "admin_console"
      value = "https://portal.acme.com/admin"
    },
  ]
}