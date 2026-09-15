resource "commander_new_database" "example" {
  title           = "Test NSF Database Record"
  folder_location = "_REPLACE_WITH_FOLDER_PATH_OR_UID_"

  login    = "root"
  hostname = "127.0.0.1"
  port     = "8000"
  type     = "SQL"
  password = "ExamplePassword123!"

  notes = "test notes"

  # Optional NSF share permissions. Map key = user email or UID;
  # value = role (viewer | share-manager | content-manager |
  # content-share-manager | full-manager). Do not include the record owner.
  share = {
    "alice@example.com" = "viewer"
    "bob@example.com"   = "full-manager"
  }
}
