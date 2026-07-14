# Secrets Manager application: name (required). Optional: shares, app_users.
# id is read-only (assigned by Keeper after creation).

resource "commander_secrets_manager" "example" {
  name = "Production Secrets App"

  app_users = [
    "alice@example.com",
    "bob@example.com",
  ]

  shares = [
    {
      secret   = "RECORD_UID_abc123"
      editable = false
    },
    {
      secret   = "SHARED_FOLDER_UID_def456"
      editable = true
    },
  ]
}
