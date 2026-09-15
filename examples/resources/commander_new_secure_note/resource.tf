resource "commander_new_secure_note" "example" {
  title = "secure note record"

  secured_note = "hey this is test"
  date         = "2026-05-25"

  notes           = "optional management notes"
  folder_location = "_REPLACE_WITH_FOLDER_PATH_OR_UID_"

  share = {
    "alice@example.com" = "full-manager"
  }
}
