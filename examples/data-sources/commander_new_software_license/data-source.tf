# Look up a New (NSF) softwareLicense record by UID.
data "commander_new_software_license" "example" {
  software_license = "_REPLACE_WITH_RECORD_UID_"
}

output "software_license_id" {
  value = data.commander_new_software_license.example.id
}

output "software_license_title" {
  value = data.commander_new_software_license.example.title
}

output "software_license_key" {
  value     = data.commander_new_software_license.example.software_license_key
  sensitive = true
}

output "software_license_expiration_date" {
  value = data.commander_new_software_license.example.expiration_date
}

output "software_license_date_active" {
  value = data.commander_new_software_license.example.date_active
}

output "software_license_custom" {
  value     = data.commander_new_software_license.example.custom
  sensitive = true
}

output "software_license_share" {
  value = data.commander_new_software_license.example.share
}
