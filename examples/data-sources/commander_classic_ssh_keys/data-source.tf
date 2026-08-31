# Look up a classic sshKeys record by title or UID.
data "commander_classic_ssh_keys" "example" {
  ssh_keys = "_REPLACE_WITH_RECORD_TITLE_OR_UID_"
}

output "ssh_keys_id" {
  value = data.commander_classic_ssh_keys.example.id
}

output "ssh_keys_title" {
  value = data.commander_classic_ssh_keys.example.title
}

output "ssh_keys_login" {
  value = data.commander_classic_ssh_keys.example.login
}

output "ssh_keys_hostname" {
  value = data.commander_classic_ssh_keys.example.hostname
}

output "ssh_keys_port" {
  value = data.commander_classic_ssh_keys.example.port
}

output "ssh_keys_passphrase" {
  value     = data.commander_classic_ssh_keys.example.passphrase
  sensitive = true
}

output "ssh_keys_public_key" {
  value     = data.commander_classic_ssh_keys.example.public_key
  sensitive = true
}

output "ssh_keys_private_key" {
  value     = data.commander_classic_ssh_keys.example.private_key
  sensitive = true
}

output "ssh_keys_custom" {
  value     = data.commander_classic_ssh_keys.example.custom
  sensitive = true
}

output "ssh_keys_share" {
  value = data.commander_classic_ssh_keys.example.share
}
