# Look up a membership record by title or UID.
data "commander_classic_membership" "example" {
  membership = "_REPLACE_WITH_RECORD_TITLE_OR_UID_"
}

output "membership_id" {
  value = data.commander_classic_membership.example.id
}

output "membership_title" {
  value = data.commander_classic_membership.example.title
}

output "membership_notes" {
  value = data.commander_classic_membership.example.notes
}

output "membership_folder_location" {
  value = data.commander_classic_membership.example.folder_location
}

output "membership_account_number" {
  value     = data.commander_classic_membership.example.account_number
  sensitive = true
}

output "membership_name" {
  value = data.commander_classic_membership.example.name
}

output "membership_password" {
  value     = data.commander_classic_membership.example.password
  sensitive = true
}

output "membership_custom" {
  value = data.commander_classic_membership.example.custom
}

output "membership_share" {
  value = data.commander_classic_membership.example.share
}
