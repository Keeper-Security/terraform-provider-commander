# Look up an SSN Card (Identity Card) record by title or UID.
data "commander_classic_ssn_card" "example" {
  ssn_card = "_REPLACE_WITH_RECORD_TITLE_OR_UID_"
}

output "ssn_card_id" {
  value = data.commander_classic_ssn_card.example.id
}

output "ssn_card_title" {
  value = data.commander_classic_ssn_card.example.title
}

output "ssn_card_notes" {
  value = data.commander_classic_ssn_card.example.notes
}

output "ssn_card_folder_location" {
  value = data.commander_classic_ssn_card.example.folder_location
}

output "ssn_card_account_number" {
  value     = data.commander_classic_ssn_card.example.account_number
  sensitive = true
}

output "ssn_card_name" {
  value = data.commander_classic_ssn_card.example.name
}

output "ssn_card_custom" {
  value = data.commander_classic_ssn_card.example.custom
}

output "ssn_card_share" {
  value = data.commander_classic_ssn_card.example.share
}
