# Look up an NSF Address record by title or UID.
data "commander_new_address" "example" {
  location = "_REPLACE_WITH_RECORD_TITLE_OR_UID_"
}

output "address_id" {
  value = data.commander_new_address.example.id
}

output "address_title" {
  value = data.commander_new_address.example.title
}

output "address_notes" {
  value = data.commander_new_address.example.notes
}

output "address_folder_location" {
  value = data.commander_new_address.example.folder_location
}

output "address_details" {
  value = data.commander_new_address.example.address
}

output "address_custom" {
  value = data.commander_new_address.example.custom
}

output "address_share" {
  value = data.commander_new_address.example.share
}
