# Look up an existing classic PAM machine record by UID or name.
data "commander_classic_pam_machine" "existing" {
  pam_machine = "_REPLACE_WITH_RECORD_UID_OR_NAME_"
}

output "pam_machine_id" {
  description = "Record UID from the data source read."
  value       = data.commander_classic_pam_machine.existing.id
}

output "pam_machine_title" {
  value = data.commander_classic_pam_machine.existing.title
}

output "pam_machine_hostname" {
  value = data.commander_classic_pam_machine.existing.hostname_or_ip
}

output "pam_machine_operating_system" {
  value = data.commander_classic_pam_machine.existing.operating_system
}

output "pam_machine_instance_name" {
  value = data.commander_classic_pam_machine.existing.instance_name
}

output "pam_machine_instance_id" {
  value = data.commander_classic_pam_machine.existing.instance_id
}

output "pam_machine_provider_group" {
  value = data.commander_classic_pam_machine.existing.provider_group
}

output "pam_machine_provider_region" {
  value = data.commander_classic_pam_machine.existing.provider_region
}

output "pam_machine_notes" {
  value = data.commander_classic_pam_machine.existing.notes
}

output "pam_machine_folder" {
  value = data.commander_classic_pam_machine.existing.folder
}

output "pam_machine_settings" {
  description = "PAM settings object (connection, tunnel, etc.)."
  value       = data.commander_classic_pam_machine.existing.pam_settings
}

# Per-user share permissions. Map key = user email; value = { can_share, can_edit }.
output "pam_machine_share" {
  description = "Per-user share permissions for this record."
  value       = data.commander_classic_pam_machine.existing.share
}
