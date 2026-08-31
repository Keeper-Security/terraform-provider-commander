resource "commander_classic_database" "example" {
  title = "test DB"

  login    = "root"
  hostname = "127.0.0.1"
  port     = "8000"
  type     = "SQL"
  password = "ExamplePassword123!"

  notes           = "test notes"
  folder_location = "_REPLACE_WITH_FOLDER_PATH_OR_UID_"

  share = {
    "alice@example.com" = {
      can_share = true
      can_edit  = true
    }
  }
}
