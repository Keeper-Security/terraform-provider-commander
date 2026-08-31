resource "commander_classic_software_license" "example" {
  title = "my metron software record"

  software_license_key = "132456789807867564331423565"
  expiration_date      = "2026-05-20"
  date_active          = "2026-07-09"

  notes           = "this is test record"
  folder_location = "_REPLACE_WITH_FOLDER_PATH_OR_UID_"

  share = {
    "alice@example.com" = {
      can_share = true
      can_edit  = true
    }
  }
}
