# Look up an existing classic PAM database record by UID or name.
data "commander_classic_pam_database" "existing" {
  pam_database = "_REPLACE_WITH_RECORD_UID_OR_NAME_"
}

output "pam_database_id" {
  description = "Record UID from the data source read."
  value       = data.commander_classic_pam_database.existing.id
}

output "pam_database_title" {
  value = data.commander_classic_pam_database.existing.title
}

output "pam_database_hostname" {
  value = data.commander_classic_pam_database.existing.hostname_or_ip
}

output "pam_database_use_ssl" {
  value = data.commander_classic_pam_database.existing.use_ssl
}

output "pam_database_database_id" {
  value = data.commander_classic_pam_database.existing.database_id
}

output "pam_database_database_type" {
  value = data.commander_classic_pam_database.existing.database_type
}

output "pam_database_provider_group" {
  value = data.commander_classic_pam_database.existing.provider_group
}

output "pam_database_provider_region" {
  value = data.commander_classic_pam_database.existing.provider_region
}

output "pam_database_notes" {
  value = data.commander_classic_pam_database.existing.notes
}

output "pam_database_folder" {
  value = data.commander_classic_pam_database.existing.folder_location
}

output "pam_database_settings" {
  description = "PAM settings object (connection, tunnel, etc.)."
  value       = data.commander_classic_pam_database.existing.pam_settings
}

# Per-user share permissions. Map key = user email; value = { can_share, can_edit }.
output "pam_database_share" {
  description = "Per-user share permissions for this record."
  value       = data.commander_classic_pam_database.existing.share
}
