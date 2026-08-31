resource "commander_classic_address" "full" {
  title = "Home Address - John Doe"

  address = {
    street1 = "123 Main Street"
    street2 = "Apt 4B"
    city    = "San Francisco"
    state   = "CA"
    zip     = "94105"
    country = "US"
  }

  notes           = "Primary residential address."
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
      label = "address_type"
      value = "Residential"
    },
  ]
}
