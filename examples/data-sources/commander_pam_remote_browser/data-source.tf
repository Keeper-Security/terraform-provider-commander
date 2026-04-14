# Look up an existing RBI record by UID.
data "commander_pam_remote_browser" "existing" {
  record_uid = "_REPLACE_WITH_RECORD_UID_"
}

# Typical pattern: read a record you manage elsewhere, or chain from a resource.
# data "commander_pam_remote_browser" "from_resource" {
#   record_uid = commander_pam_remote_browser.intranet_app.id
# }

output "rbi_record_id" {
  description = "Record UID from the data source read."
  value       = data.commander_pam_remote_browser.existing.id
}

output "rbi_title" {
  value = data.commander_pam_remote_browser.existing.title
}

output "rbi_target_url" {
  value = data.commander_pam_remote_browser.existing.url
}

output "rbi_notes" {
  value = data.commander_pam_remote_browser.existing.notes
}

output "rbi_folder" {
  value = data.commander_pam_remote_browser.existing.folder
}

output "rbi_settings" {
  description = "Nested settings object when returned by Commander (may be partially null)."
  value       = data.commander_pam_remote_browser.existing.pam_remote_browser_settings
}
