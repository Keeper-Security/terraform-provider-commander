# Look up a passport record by title or UID.
data "commander_classic_passport" "example" {
  passport = "_REPLACE_WITH_RECORD_TITLE_OR_UID_"
}

output "passport_id" {
  value = data.commander_classic_passport.example.id
}

output "passport_title" {
  value = data.commander_classic_passport.example.title
}

output "passport_notes" {
  value = data.commander_classic_passport.example.notes
}

output "passport_folder_location" {
  value = data.commander_classic_passport.example.folder_location
}

output "passport_account_number" {
  value     = data.commander_classic_passport.example.account_number
  sensitive = true
}

output "passport_name" {
  value = data.commander_classic_passport.example.name
}

output "passport_birth_date" {
  value = data.commander_classic_passport.example.birth_date
}

output "passport_address_ref" {
  value = data.commander_classic_passport.example.address_ref
}

output "passport_expiration_date" {
  value = data.commander_classic_passport.example.expiration_date
}

output "passport_date_issued" {
  value = data.commander_classic_passport.example.date_issued
}

output "passport_password" {
  value     = data.commander_classic_passport.example.password
  sensitive = true
}

output "passport_custom" {
  value = data.commander_classic_passport.example.custom
}

output "passport_share" {
  value = data.commander_classic_passport.example.share
}
