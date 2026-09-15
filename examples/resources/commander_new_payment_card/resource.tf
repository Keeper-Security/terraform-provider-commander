resource "commander_new_payment_card" "full" {
  title = "Corporate Amex - Acme Corp"

  payment_card = {
    card_number          = "371449635398431"
    card_expiration_date = "04/2030"
    card_security_code   = "1234"
  }

  cardholder_name = "Jane Smith"
  pin_code        = "0000"
  address_ref     = "_REPLACE_WITH_ADDRESS_RECORD_UID_"

  notes           = "Corporate card for procurement."
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
