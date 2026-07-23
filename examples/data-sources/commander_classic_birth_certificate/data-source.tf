# Look up a Birth Certificate record by title or UID.
data "commander_classic_birth_certificate" "example" {
  birth_certificate = "_REPLACE_WITH_RECORD_TITLE_OR_UID_"
}

output "birth_certificate_id" {
  value = data.commander_classic_birth_certificate.example.id
}

output "birth_certificate_title" {
  value = data.commander_classic_birth_certificate.example.title
}

output "birth_certificate_notes" {
  value = data.commander_classic_birth_certificate.example.notes
}

output "birth_certificate_folder_location" {
  value = data.commander_classic_birth_certificate.example.folder_location
}

output "birth_certificate_name" {
  value = data.commander_classic_birth_certificate.example.name
}

output "birth_certificate_birth_date" {
  value = data.commander_classic_birth_certificate.example.birth_date
}

output "birth_certificate_custom" {
  value = data.commander_classic_birth_certificate.example.custom
}

output "birth_certificate_share" {
  value = data.commander_classic_birth_certificate.example.share
}
