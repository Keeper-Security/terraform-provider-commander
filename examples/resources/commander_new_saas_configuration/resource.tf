resource "commander_new_saas_configuration" "example" {
  title = "SaaS Rotation Config"

  notes           = "Configuration for SaaS account rotation."
  folder_location = "_REPLACE_WITH_FOLDER_PATH_OR_UID_"

  configuration = "_REPLACE_WITH_CONFIGURATION_UID_"
  gateway       = "_REPLACE_WITH_GATEWAY_UID_"

  custom = [
    {
      type  = "text"
      label = "SaaS Type"
      value = "REST" // Valid values changes based upon provied gateway and configuration.
    },
    {
      type  = "text"
      label = "Active"
      value = "TRUE" // Valid values are "TRUE" or "FALSE"
    },
  ]

  share = {
    "alice@example.com" = "full-manager"
  }
}
