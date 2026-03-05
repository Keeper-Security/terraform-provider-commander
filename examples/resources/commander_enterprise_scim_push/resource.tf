# commander_enterprise_scim_push: one-time push of SCIM data to a Keeper SCIM endpoint.
#
# Required: scim_id, source (google | ad | record), record (record UID), auto_approve.
# Optional: managed_company (MSP only).
#
# When you apply, Terraform runs the push once. Changing any value runs the push again.
# No read or delete on the server; removing the resource only removes it from state.

resource "commander_enterprise_scim_push" "example" {
  scim_id      = "1234567890"      # ID of your SCIM configuration
  source       = "google"          # google | ad | record
  record       = "your_record_uid" # Record UID with SCIM configuration
  auto_approve = true              # true = auto-approve SCIM teams; false = require approval

  # Optional, MSP only: run in a specific managed company
  # managed_company = "Acme Corp"
}
