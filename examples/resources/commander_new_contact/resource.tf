resource "commander_new_contact" "full" {
  title = "Jane Smith - Acme Corp"

  name = {
    first  = "Jane"
    middle = "A."
    last   = "Smith"
  }

  company = "Acme Corp"
  email   = "jane.smith@acme.com"

  phone = [
    {
      type   = "Mobile"
      region = "US"
      number = "+1-555-0100"
    },
    {
      type   = "Work"
      region = "US"
      number = "+1-555-0200"
      ext    = "4321"
    },
  ]

  address_ref     = "_REPLACE_WITH_ADDRESS_RECORD_UID_"
  notes           = "Primary contact for Acme Corp procurement."
  folder_location = "_REPLACE_WITH_FOLDER_UID_"

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
      label = "department"
      value = "Engineering"
    },
    {
      type  = "text"
      label = "slack_channel"
      value = "#eng-support"
    },
    {
      type  = "address"
      label = "office_address"
      value = jsonencode({
        street1 = "100 Main Street"
        street2 = "apt 2"
        city    = "San Francisco"
        state   = "CA"
        zip     = "94105"
        country = "US"
      })
    },
  ]
}
