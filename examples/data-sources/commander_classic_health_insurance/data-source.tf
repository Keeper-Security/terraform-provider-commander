# Look up a health insurance record by title or UID.
data "commander_classic_health_insurance" "example" {
  health_insurance = "_REPLACE_WITH_RECORD_TITLE_OR_UID_"
}

output "health_insurance_id" {
  value = data.commander_classic_health_insurance.example.id
}

output "health_insurance_title" {
  value = data.commander_classic_health_insurance.example.title
}

output "health_insurance_notes" {
  value = data.commander_classic_health_insurance.example.notes
}

output "health_insurance_folder_location" {
  value = data.commander_classic_health_insurance.example.folder_location
}

output "health_insurance_account_number" {
  value     = data.commander_classic_health_insurance.example.account_number
  sensitive = true
}

output "health_insurance_name" {
  value = data.commander_classic_health_insurance.example.name
}

output "health_insurance_login" {
  value = data.commander_classic_health_insurance.example.login
}

output "health_insurance_password" {
  value     = data.commander_classic_health_insurance.example.password
  sensitive = true
}

output "health_insurance_website_address" {
  value = data.commander_classic_health_insurance.example.website_address
}

output "health_insurance_custom" {
  value = data.commander_classic_health_insurance.example.custom
}

output "health_insurance_share" {
  value = data.commander_classic_health_insurance.example.share
}
