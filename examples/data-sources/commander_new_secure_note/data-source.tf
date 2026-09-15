# Look up a New (NSF) Secure Note record by UID.
data "commander_new_secure_note" "example" {
  secure_note = "_REPLACE_WITH_RECORD_UID_"
}

output "secure_note_id" {
  value = data.commander_new_secure_note.example.id
}

output "secure_note_title" {
  value = data.commander_new_secure_note.example.title
}

output "secure_note_secured_note" {
  value     = data.commander_new_secure_note.example.secured_note
  sensitive = true
}

output "secure_note_date" {
  value = data.commander_new_secure_note.example.date
}

output "secure_note_custom" {
  value     = data.commander_new_secure_note.example.custom
  sensitive = true
}

output "secure_note_share" {
  value = data.commander_new_secure_note.example.share
}
