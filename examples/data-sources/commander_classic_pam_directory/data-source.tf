# ------------------------------------------------------------------
# Example 1 - Look up a classic PAM directory record by UID or name
# ------------------------------------------------------------------
data "commander_classic_pam_directory" "existing" {
  pam_directory = "_REPLACE_WITH_RECORD_UID_OR_NAME_"
}

output "pam_directory_id" {
  description = "Record UID from the data source read."
  value       = data.commander_classic_pam_directory.existing.id
}

output "pam_directory_title" {
  value = data.commander_classic_pam_directory.existing.title
}

output "pam_directory_hostname" {
  value = data.commander_classic_pam_directory.existing.hostname_or_ip
}

output "pam_directory_use_ssl" {
  value = data.commander_classic_pam_directory.existing.use_ssl
}

output "pam_directory_domain_name" {
  value = data.commander_classic_pam_directory.existing.domain_name
}

output "pam_directory_alternative_ips" {
  value = data.commander_classic_pam_directory.existing.alternative_ips
}

output "pam_directory_directory_id" {
  value = data.commander_classic_pam_directory.existing.directory_id
}

output "pam_directory_directory_type" {
  value = data.commander_classic_pam_directory.existing.directory_type
}

output "pam_directory_user_match" {
  value = data.commander_classic_pam_directory.existing.user_match
}

output "pam_directory_provider_group" {
  value = data.commander_classic_pam_directory.existing.provider_group
}

output "pam_directory_provider_region" {
  value = data.commander_classic_pam_directory.existing.provider_region
}

output "pam_directory_notes" {
  value = data.commander_classic_pam_directory.existing.notes
}

output "pam_directory_folder" {
  value = data.commander_classic_pam_directory.existing.folder
}

output "pam_directory_settings" {
  description = "PAM settings object (connection, tunnel, etc.)."
  value       = data.commander_classic_pam_directory.existing.pam_settings
}

# Per-user share permissions. Map key = user email; value = { can_share, can_edit }.
output "pam_directory_share" {
  description = "Per-user share permissions for this record."
  value       = data.commander_classic_pam_directory.existing.share
}

# ------------------------------------------------------------------
# Example 2 - Use data source to reference a directory in another resource
# ------------------------------------------------------------------
data "commander_classic_pam_directory" "corp_dc" {
  pam_directory = "_REPLACE_WITH_RECORD_UID_OR_NAME_"
}

# Reference the directory's configuration UID in another resource.
output "corp_dc_configuration" {
  description = "PAM configuration UID linked to this directory."
  value       = data.commander_classic_pam_directory.corp_dc.pam_settings.configuration
}
