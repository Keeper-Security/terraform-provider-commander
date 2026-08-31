# Look up a classic serverCredentials record by UID.
data "commander_classic_server" "example" {
  server = "_REPLACE_WITH_RECORD_UID_"
}

output "server_id" {
  value = data.commander_classic_server.example.id
}

output "server_title" {
  value = data.commander_classic_server.example.title
}

output "server_login" {
  value = data.commander_classic_server.example.login
}

output "server_hostname" {
  value = data.commander_classic_server.example.hostname
}

output "server_port" {
  value = data.commander_classic_server.example.port
}

output "server_password" {
  value     = data.commander_classic_server.example.password
  sensitive = true
}

output "server_custom" {
  value     = data.commander_classic_server.example.custom
  sensitive = true
}

output "server_share" {
  value = data.commander_classic_server.example.share
}
