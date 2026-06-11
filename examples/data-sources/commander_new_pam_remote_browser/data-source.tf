# Look up an existing new (nested-shared) PAM remote browser record by UID or name.
data "commander_new_pam_remote_browser" "existing" {
  remote_browser = "_REPLACE_WITH_RECORD_UID_OR_NAME_"
}

output "rbi_record_id" {
  description = "Record UID from the data source read."
  value       = data.commander_new_pam_remote_browser.existing.id
}

output "rbi_title" {
  value = data.commander_new_pam_remote_browser.existing.title
}

output "rbi_target_url" {
  value = data.commander_new_pam_remote_browser.existing.url
}

output "rbi_notes" {
  value = data.commander_new_pam_remote_browser.existing.notes
}

output "rbi_folder" {
  value = data.commander_new_pam_remote_browser.existing.folder_location
}

output "rbi_settings" {
  description = "Nested settings object when returned by Commander (may be partially null)."
  value       = data.commander_new_pam_remote_browser.existing.pam_remote_browser_settings
}

# Per-user share permissions. Map key = user email; value = role string
# (viewer | share-manager | content-manager | content-share-manager | full-manager).
output "rbi_share" {
  description = "Per-user share permissions for this record."
  value       = data.commander_new_pam_remote_browser.existing.share
}
