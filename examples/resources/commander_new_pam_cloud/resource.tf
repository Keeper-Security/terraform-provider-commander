# ------------------------------------------------------------------
# Example 1 - Minimal: only required fields (title)
# ------------------------------------------------------------------
resource "commander_new_pam_cloud" "minimal" {
  title = "PAM Cloud - aws-prod"
}

# ------------------------------------------------------------------
# Example 2 - All cloud fields + share permissions
# ------------------------------------------------------------------
resource "commander_new_pam_cloud" "full_fields" {
  title              = "PAM Cloud - aws-prod"
  url                = "https://console.aws.amazon.com"
  account_identifier = "123456789012"
  notes              = "Production AWS account."
  folder_location    = "_REPLACE_WITH_NESTED_SHARED_FOLDER_UID_"

  # ----------------------------------------------------------------
  # Per-user share permissions (optional).
  # Map key = user email. Each value is one of:
  #   viewer, share-manager, content-manager,
  #   content-share-manager, full-manager.
  # Omit the block entirely to skip share reconciliation.
  # ----------------------------------------------------------------
  share = {
    "alice@example.com" = "full-manager"
    "bob@example.com"   = "viewer"
  }
}
