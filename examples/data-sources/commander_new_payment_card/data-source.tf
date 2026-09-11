# Look up a payment card record by title or UID.
data "commander_new_payment_card" "example" {
  bank_card = "_REPLACE_WITH_RECORD_TITLE_OR_UID_"
}

output "payment_card_id" {
  value = data.commander_new_payment_card.example.id
}

output "payment_card_title" {
  value = data.commander_new_payment_card.example.title
}

output "payment_card_notes" {
  value = data.commander_new_payment_card.example.notes
}

output "payment_card_folder_location" {
  value = data.commander_new_payment_card.example.folder_location
}

output "payment_card_details" {
  value     = data.commander_new_payment_card.example.payment_card
  sensitive = true
}

output "payment_card_cardholder_name" {
  value = data.commander_new_payment_card.example.cardholder_name
}

output "payment_card_pin_code" {
  value     = data.commander_new_payment_card.example.pin_code
  sensitive = true
}

output "payment_card_address_ref" {
  value = data.commander_new_payment_card.example.address_ref
}

output "payment_card_custom" {
  value = data.commander_new_payment_card.example.custom
}

output "payment_card_share" {
  value = data.commander_new_payment_card.example.share
}
