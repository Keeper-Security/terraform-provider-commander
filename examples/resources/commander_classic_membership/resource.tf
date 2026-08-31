resource "commander_classic_membership" "full" {
  title = "Gold's Gym - Jane Smith"

  account_number = "GYM123456"

  name = {
    first = "Jane"
    last  = "Smith"
  }

  password = "ExamplePassword123!"

  notes           = "Premium gym membership."
  folder_location = "_REPLACE_WITH_FOLDER_PATH_OR_UID_"

  share = {
    "alice@example.com" = {
      can_share = true
      can_edit  = true
    }
    "bob@example.com" = {
      can_share = false
      can_edit  = true
    }
    "viewer@example.com" = {
      can_share = false
      can_edit  = false
    }
  }

  /* Example of custom fields
    - For Complex types, use jsonencode(JSON) matching the Keeper field schema. 
    - For more information, see the Keeper field schema documentation: https://docs.keeper.io/en/keeperpam/secrets-manager/about/field-record-types
  */
  custom = [
    {
      type  = "text"
      label = "membership_type"
      value = "Premium"
    },
    {
      type  = "date"
      label = "expiration_date"
      value = "2024-12-31"
    },
  ]
}
