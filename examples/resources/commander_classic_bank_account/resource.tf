resource "commander_classic_bank_account" "full" {
  title = "Corporate Checking - Acme Corp"

  bank_account = {
    account_type   = "Checking"
    routing_number = "021000021"
    account_number = "000123456789"
  }

  name = {
    first = "Jane"
    last  = "Smith"
  }

  login           = "jane.smith"
  password        = "ExamplePassword123!"
  website_address = "https://examplebank.com"
  card_ref        = "_REPLACE_WITH_PAYMENT_CARD_RECORD_UID_"

  notes           = "Corporate operating account."
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
      label = "billing_department"
      value = "Finance"
    },
    {
      type  = "text"
      label = "cost_center"
      value = "CC-1002"
    },
  ]
}
