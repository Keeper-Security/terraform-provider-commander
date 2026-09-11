resource "commander_new_server" "example" {
  title = "Test Server Record Type"

  login    = "root"
  password = "ExamplePassword123!"
  hostname = "localhost"
  port     = "22"

  notes           = "test notes"
  folder_location = "_REPLACE_WITH_FOLDER_PATH_OR_UID_"

  share = {
    "alice@example.com" = "full-manager"
  }
}