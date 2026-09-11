resource "commander_new_driver_license" "full" {
  title = "California Driver's License - John Doe"

  account_number = "DL123456789"

  name = {
    first = "John"
    last  = "Doe"
  }

  birth_date      = "1990-01-15"
  expiration_date = "2030-01-15"
  address_ref     = "_REPLACE_WITH_ADDRESS_RECORD_UID_"

  notes           = "Personal driver's license."
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
      label = "license_class"
      value = "Class C"
    },
  ]
}
