resource "commander_new_login" "example" {
  title           = "Test NSF Login Record"
  folder_location = "_REPLACE_WITH_FOLDER_PATH_OR_UID_"

  login           = "user@example.com"
  password        = "ExamplePassword123!"
  website_address = "https://example.com"

  notes = "test notes"

  # Optional NSF share permissions. Map key = user email or UID;
  # value = role (viewer | share-manager | content-manager |
  # content-share-manager | full-manager). Do not include the record owner.
  share = {
    "alice@example.com" = "viewer"
    "bob@example.com"   = "full-manager"
  }
}
