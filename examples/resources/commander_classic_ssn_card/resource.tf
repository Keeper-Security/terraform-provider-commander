resource "commander_classic_ssn_card" "full" {
  title = "Social Security Card - John Doe"

  account_number = "123-45-6789"

  name = {
    first = "John"
    last  = "Doe"
  }

  notes           = "Personal social security card."
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
      label = "issue_date"
      value = "2008-01-15"
    },
    {
      type  = "text"
      label = "issue_state"
      value = "New York"
    },
  ]
}
