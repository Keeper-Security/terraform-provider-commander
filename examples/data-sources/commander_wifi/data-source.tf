# commander_wifi (data source)
#
# Reads an existing WiFi credentials record (`wifiCredentials`) from the vault
# by its record UID. Returns title, ssid, password (sensitive), encryption,
# is_ssid_hidden, folder, notes, and any custom fields stored on the record.

terraform {
  required_providers {
    commander = {
      source = "keeper-security/commander"
    }
  }
}

provider "commander" {
  service_mode_url     = "http://localhost:8080/api/v2/"
  service_mode_api_key = "XXXXXXXXXXXXXX"
}

###############################################################################
# Usage 1 - Look up a WiFi record by record UID
###############################################################################

data "commander_wifi" "office_guest" {
  record_uid = "_REPLACE_WITH_RECORD_UID_"
}

###############################################################################
# Usage 2 - Chain from a managed resource (no hard-coded UID)
###############################################################################

# data "commander_wifi" "from_managed_resource" {
#   record_uid = commander_wifi.home.id
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
