# Look up a classic database record by UID.
data "commander_classic_database" "example" {
  database = "_REPLACE_WITH_RECORD_UID_"
}

output "database_id" {
  value = data.commander_classic_database.example.id
}

output "database_title" {
  value = data.commander_classic_database.example.title
}

output "database_login" {
  value = data.commander_classic_database.example.login
}

output "database_hostname" {
  value = data.commander_classic_database.example.hostname
}

output "database_port" {
  value = data.commander_classic_database.example.port
}

output "database_type" {
  value = data.commander_classic_database.example.type
}

output "database_password" {
  value     = data.commander_classic_database.example.password
  sensitive = true
}

output "database_custom" {
  value     = data.commander_classic_database.example.custom
  sensitive = true
}

output "database_share" {
  value = data.commander_classic_database.example.share
}
