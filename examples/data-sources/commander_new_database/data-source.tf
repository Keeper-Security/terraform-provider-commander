# Look up a New (NSF) database credentials record by UID.
data "commander_new_database" "example" {
  database = "_REPLACE_WITH_RECORD_UID_"
}

output "database_id" {
  value = data.commander_new_database.example.id
}

output "database_title" {
  value = data.commander_new_database.example.title
}

output "database_login" {
  value = data.commander_new_database.example.login
}

output "database_hostname" {
  value = data.commander_new_database.example.hostname
}

output "database_port" {
  value = data.commander_new_database.example.port
}

output "database_type" {
  value = data.commander_new_database.example.type
}

output "database_password" {
  value     = data.commander_new_database.example.password
  sensitive = true
}

output "database_notes" {
  value = data.commander_new_database.example.notes
}

output "database_folder" {
  value = data.commander_new_database.example.folder_location
}

output "database_custom" {
  value     = data.commander_new_database.example.custom
  sensitive = true
}

# Per-user share permissions. Map key = user email or UID; value = role string
# (viewer | share-manager | content-manager | content-share-manager | full-manager).
output "database_share" {
  value = data.commander_new_database.example.share
}
