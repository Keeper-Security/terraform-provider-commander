# Look up an existing new (NSF) PAM cloud record by UID.
data "commander_new_pam_cloud" "existing" {
  pam_cloud = "_REPLACE_WITH_RECORD_UID_"
}

output "pam_cloud_id" {
  description = "Record UID from the data source read."
  value       = data.commander_new_pam_cloud.existing.id
}

output "pam_cloud_title" {
  value = data.commander_new_pam_cloud.existing.title
}

output "pam_cloud_url" {
  value = data.commander_new_pam_cloud.existing.url
}

output "pam_cloud_account_identifier" {
  value = data.commander_new_pam_cloud.existing.account_identifier
}

output "pam_cloud_notes" {
  value = data.commander_new_pam_cloud.existing.notes
}

output "pam_cloud_folder" {
  value = data.commander_new_pam_cloud.existing.folder_location
}

# Per-user share permissions. Map key = user email; value = role (viewer, share-manager, content-manager, content-share-manager, full-manager).
output "pam_cloud_share" {
  description = "Per-user share permissions for this record."
  value       = data.commander_new_pam_cloud.existing.share
}
