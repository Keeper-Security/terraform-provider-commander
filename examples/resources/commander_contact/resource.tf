resource "commander_contact" "full" {
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

  address_ref = "_REPLACE_WITH_ADDRESS_RECORD_UID_"
  notes       = "Primary contact for Acme Corp procurement."
  folder      = "_REPLACE_WITH_FOLDER_UID_"

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
}