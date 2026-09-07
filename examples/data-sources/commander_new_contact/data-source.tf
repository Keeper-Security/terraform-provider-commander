# Look up a New (NSF) contact record by title or UID.
data "commander_new_contact" "example" {
  contact = "_REPLACE_WITH_RECORD_TITLE_OR_UID_"
}

output "contact_id" {
  value = data.commander_new_contact.example.id
}

output "contact_title" {
  value = data.commander_new_contact.example.title
}

output "contact_notes" {
  value = data.commander_new_contact.example.notes
}

output "contact_folder_location" {
  value = data.commander_new_contact.example.folder_location
}

output "contact_name" {
  value = data.commander_new_contact.example.name
}

output "contact_name_first" {
  value = data.commander_new_contact.example.name.first
}

output "contact_name_middle" {
  value = data.commander_new_contact.example.name.middle
}

output "contact_name_last" {
  value = data.commander_new_contact.example.name.last
}

output "contact_company" {
  value = data.commander_new_contact.example.company
}

output "contact_email" {
  value = data.commander_new_contact.example.email
}

output "contact_phone" {
  value = data.commander_new_contact.example.phone
}

output "contact_address_ref" {
  value = data.commander_new_contact.example.address_ref
}

output "contact_custom" {
  value = data.commander_new_contact.example.custom
}

output "contact_share" {
  value = data.commander_new_contact.example.share
}
