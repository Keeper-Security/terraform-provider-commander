resource "commander_new_health_insurance" "example" {
  title = "Blue Cross Blue Shield - Jane Smith"

  account_number = "12345678901"

  name = {
    first = "Jane"
    last  = "Smith"
  }

  login           = "jane.smith"
  password        = "ExamplePassword123!"
  website_address = "https://bcbs.com"

  notes           = "Family health insurance policy."
  folder_location = "_REPLACE_WITH_FOLDER_PATH_OR_UID_"

  share = {
    "alice@example.com" = "full-manager"
  }

  /* Example of custom fields
    - For Complex types, use jsonencode(JSON) matching the Keeper field schema.
    - For more information, see the Keeper field schema documentation: https://docs.keeper.io/en/keeperpam/secrets-manager/about/field-record-types
  */
  custom = [
    {
      type  = "text"
      label = "plan_type"
      value = "PPO"
    },
    {
      type  = "text"
      label = "group_number"
      value = "12345"
    },
  ]
}
