# Look up a nested (NSF) login record by title or UID.
data "commander_new_login" "example" {
  login_record = "_REPLACE_WITH_RECORD_TITLE_OR_UID_"
}

output "login_id" {
  value = data.commander_new_login.example.id
}

output "login_title" {
  value = data.commander_new_login.example.title
}

output "login_username" {
  value = data.commander_new_login.example.login
}

output "login_website_address" {
  value = data.commander_new_login.example.website_address
}

output "login_password" {
  value     = data.commander_new_login.example.password
  sensitive = true
}

output "login_custom" {
  value     = data.commander_new_login.example.custom
  sensitive = true
}

output "login_share" {
  value = data.commander_new_login.example.share
}
