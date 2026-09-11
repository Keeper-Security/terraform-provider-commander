# Look up a driver's license record by title or UID.
data "commander_new_driver_license" "example" {
  driver_license = "_REPLACE_WITH_RECORD_TITLE_OR_UID_"
}

output "driver_license_id" {
  value = data.commander_new_driver_license.example.id
}

output "driver_license_title" {
  value = data.commander_new_driver_license.example.title
}

output "driver_license_notes" {
  value = data.commander_new_driver_license.example.notes
}

output "driver_license_folder_location" {
  value = data.commander_new_driver_license.example.folder_location
}

output "driver_license_account_number" {
  value     = data.commander_new_driver_license.example.account_number
  sensitive = true
}

output "driver_license_name" {
  value = data.commander_new_driver_license.example.name
}

output "driver_license_birth_date" {
  value = data.commander_new_driver_license.example.birth_date
}

output "driver_license_address_ref" {
  value = data.commander_new_driver_license.example.address_ref
}

output "driver_license_expiration_date" {
  value = data.commander_new_driver_license.example.expiration_date
}

output "driver_license_custom" {
  value = data.commander_new_driver_license.example.custom
}

output "driver_license_share" {
  value = data.commander_new_driver_license.example.share
}
