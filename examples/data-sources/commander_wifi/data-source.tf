###############################################################################
# Usage 1 - Look up a WiFi record by record UID
###############################################################################

data "commander_wifi" "office_guest" {
  wifi = "_REPLACE_WITH_RECORD_UID_OR_TITLE_"
}

###############################################################################
# Usage 2 - Chain from a managed resource (no hard-coded UID)
###############################################################################

# data "commander_wifi" "from_managed_resource" {
#   wifi = commander_wifi.home.id
# }

###############################################################################
# Outputs - top-level fields
###############################################################################

output "wifi_id" {
  description = "Record UID of the WiFi record."
  value       = data.commander_wifi.office_guest.id
}

output "wifi_title" {
  value = data.commander_wifi.office_guest.title
}

output "wifi_ssid" {
  description = "WiFi network name (SSID)."
  value       = data.commander_wifi.office_guest.ssid
}

output "wifi_encryption" {
  description = "Encryption type: wep | wpa | noEncryption."
  value       = data.commander_wifi.office_guest.encryption
}

output "wifi_is_ssid_hidden" {
  value = data.commander_wifi.office_guest.is_ssid_hidden
}

output "wifi_folder" {
  value = data.commander_wifi.office_guest.folder
}

output "wifi_notes" {
  value = data.commander_wifi.office_guest.notes
}

###############################################################################
# Outputs - sensitive fields
###############################################################################

output "wifi_password" {
  description = "WiFi password (returns null for open networks)."
  value       = data.commander_wifi.office_guest.password
  sensitive   = true
}

###############################################################################
# Outputs - custom fields
#
# `custom` is a list of objects: { type, label, value, sensitive }.
# Always check length / null before indexing; the list is empty when the
# record has no custom fields.
###############################################################################

output "wifi_custom" {
  description = "All custom fields stored on the record."
  value       = data.commander_wifi.office_guest.custom
  sensitive   = true
}

###############################################################################
# Outputs - sharing
#
# `share` is a map keyed by email address with { can_share, can_edit } values.
# The map is null when the record is not shared with any user.
###############################################################################

output "wifi_share" {
  description = "Users this record is shared with and their permissions."
  value       = data.commander_wifi.office_guest.share
}
