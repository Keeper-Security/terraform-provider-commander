# Look up a classic saasConfiguration record by UID.
data "commander_classic_saas_configuration" "example" {
  saas_configuration = "_REPLACE_WITH_RECORD_UID_"
}

output "saas_configuration_id" {
  value = data.commander_classic_saas_configuration.example.id
}

output "saas_configuration_title" {
  value = data.commander_classic_saas_configuration.example.title
}

output "saas_configuration_custom" {
  value     = data.commander_classic_saas_configuration.example.custom
  sensitive = true
}

output "saas_configuration_share" {
  value = data.commander_classic_saas_configuration.example.share
}
