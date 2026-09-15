resource "commander_new_saas_configuration" "example" {
  title = "SaaS Rotation Config"

  notes           = "Configuration for SaaS account rotation."
  folder_location = "_REPLACE_WITH_FOLDER_PATH_OR_UID_"

  custom = [
    {
      type  = "text"
      label = "SaaS Type"
      value = "Okta"
    },
    {
      type  = "text"
      label = "AppName"
      value = "Example SaaS App"
    },
  ]

  share = {
    "alice@example.com" = "full-manager"
  }
}
