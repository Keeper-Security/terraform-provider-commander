resource "commander_new_birth_certificate" "full" {
  title = "Birth Certificate - John Doe"

  name = {
    first  = "John"
    middle = "Michael"
    last   = "Doe"
  }

  birth_date = "1990-05-15"

  notes           = "Official birth certificate record."
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
      label = "place_of_birth"
      value = "San Francisco, CA"
    },
    {
      type  = "text"
      label = "certificate_number"
      value = "BC-1990-001234"
    },
  ]
}
