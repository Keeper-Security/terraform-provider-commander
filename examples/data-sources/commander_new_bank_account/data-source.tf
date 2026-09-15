# Look up a New (NSF) bankAccount record by UID.
data "commander_new_bank_account" "example" {
  account = "_REPLACE_WITH_RECORD_UID_"
}

output "bank_account_id" {
  value = data.commander_new_bank_account.example.id
}

output "bank_account_title" {
  value = data.commander_new_bank_account.example.title
}

output "bank_account_notes" {
  value = data.commander_new_bank_account.example.notes
}

output "bank_account_folder_location" {
  value = data.commander_new_bank_account.example.folder_location
}

output "bank_account_details" {
  value     = data.commander_new_bank_account.example.bank_account
  sensitive = true
}

output "bank_account_name" {
  value = data.commander_new_bank_account.example.name
}

output "bank_account_login" {
  value = data.commander_new_bank_account.example.login
}

output "bank_account_password" {
  value     = data.commander_new_bank_account.example.password
  sensitive = true
}

output "bank_account_website_address" {
  value = data.commander_new_bank_account.example.website_address
}

output "bank_account_card_ref" {
  value = data.commander_new_bank_account.example.card_ref
}

output "bank_account_custom" {
  value     = data.commander_new_bank_account.example.custom
  sensitive = true
}

output "bank_account_share" {
  value = data.commander_new_bank_account.example.share
}
