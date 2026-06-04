# commander_new_pam_user (data source)
#
# Reads an existing new (nested-shared) PAM User record (`pamUser`) from the
# vault by record UID or name. Returns title, login, password (sensitive),
# folder, notes, distinguished_name, private_pem_key (sensitive),
# connect_database, managed, rotation_settings, and per-user share permissions.

###############################################################################
# Usage 1 - Look up a new PAM User by record UID or name
###############################################################################

data "commander_new_pam_user" "mysql_app_account" {
  pam_user = "_REPLACE_WITH_RECORD_UID_OR_NAME_"
}

###############################################################################
# Usage 2 - Chain from a managed resource (no hard-coded UID)
###############################################################################

# data "commander_new_pam_user" "from_managed_resource" {
#   pam_user = commander_new_pam_user.mysql_app_account.id
# }

###############################################################################
# Outputs - top-level fields
###############################################################################

output "pam_user_id" {
  description = "Record UID of the PAM User."
  value       = data.commander_new_pam_user.mysql_app_account.id
}

output "pam_user_title" {
  value = data.commander_new_pam_user.mysql_app_account.title
}

output "pam_user_login" {
  value = data.commander_new_pam_user.mysql_app_account.login
}

output "pam_user_folder" {
  value = data.commander_new_pam_user.mysql_app_account.folder
}

output "pam_user_distinguished_name" {
  value = data.commander_new_pam_user.mysql_app_account.distinguished_name
}

output "pam_user_connect_database" {
  value = data.commander_new_pam_user.mysql_app_account.connect_database
}

output "pam_user_managed" {
  value = data.commander_new_pam_user.mysql_app_account.managed
}

###############################################################################
# Outputs - sensitive fields
###############################################################################

output "pam_user_password" {
  value     = data.commander_new_pam_user.mysql_app_account.password
  sensitive = true
}

output "pam_user_private_pem_key" {
  value     = data.commander_new_pam_user.mysql_app_account.private_pem_key
  sensitive = true
}

###############################################################################
# Outputs - rotation_settings (may be null if rotation isn't configured)
###############################################################################

output "pam_user_rotation_enabled" {
  description = "Whether automated rotation is enabled (null if not configured)."
  value = (
    data.commander_new_pam_user.mysql_app_account.rotation_settings != null
    ? data.commander_new_pam_user.mysql_app_account.rotation_settings.enabled
    : null
  )
}

output "pam_user_rotation_profile" {
  description = "Rotation profile type: general | iam_user | scripts_only (null if not configured)."
  value = (
    data.commander_new_pam_user.mysql_app_account.rotation_settings != null
    ? data.commander_new_pam_user.mysql_app_account.rotation_settings.rotation_profile
    : null
  )
}

output "pam_user_rotation_settings" {
  description = "Full rotation settings block (null if not configured)."
  value       = data.commander_new_pam_user.mysql_app_account.rotation_settings
}

###############################################################################
# Outputs - share (per-user permissions, populated from the API)
###############################################################################

output "pam_user_share" {
  description = "Per-user share permissions. Map key = user email; value = role string (viewer | share-manager | content-manager | content-share-manager | full-manager)."
  value       = data.commander_new_pam_user.mysql_app_account.share
}
