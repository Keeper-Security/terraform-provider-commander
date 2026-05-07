# commander_epm_policy — create and manage an EPM (Endpoint Policy Management) policy.
#
# Required for every policy:
#   policy_name, policy_type, status
#
# policy_type: elevation | file_access | command | least_privilege
# status:      enforce | monitor | monitor_and_notify | off
#               (least_privilege: only enforce or off.)
#
# Collections: use collection IDs from your Keeper EPM tenant, or ["*"] alone for “all” (do not mix * with other IDs).
#
# Quick rules (see commented blocks below for full examples):
#   elevation / file_access + enforce → control (≥1) + user_groups + machine_collections + applications
#   elevation + monitor(_and_notify) → all three collections; control optional
#   file_access + monitor(_and_notify) → at least one of the three collection sets non-empty
#   command + enforce → control (≥1) + user_groups + machine_collections; never set applications
#   command + monitor(_and_notify) → user_groups + machine_collections; never set applications
#   elevation/file_access/command + off → only required attributes (omit optional sets)
#   least_privilege → only policy_name, policy_type, status, and optionally machine_collections
#
# Computed: id — use with data "commander_epm_policy" or terraform import.

# ---------------------------------------------------------------------------
# Default example (active): elevation + enforce, full scope, weekday schedule
# ---------------------------------------------------------------------------
resource "commander_epm_policy" "example" {
  policy_name = "Terraform example — elevation enforce"
  policy_type = "elevation"
  status      = "enforce"

  control = ["audit", "justify"]

  user_groups         = ["*"]
  machine_collections = ["*"]
  applications        = ["*"]

  day_filter  = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday"]
  time_filter = ["9-17"]

  # date_filter = ["2025-01-01:2025-12-31"]

  # Optional, MSP only — scope create/read/update/delete to a managed company
  # managed_company = "Acme Corp"
}

output "epm_policy_id" {
  description = "Keeper-assigned ID after create (use for data source or import)."
  value       = commander_epm_policy.example.id
}

# =============================================================================
# DETAILED EXAMPLES
# =============================================================================

# --- All control actions (enforce policies only; not valid for least_privilege) ---
#
# control = ["audit", "notify", "mfa", "justify", "approval"]

# --- Multiple time windows (ranges must not overlap; not for least_privilege) ---
#
# time_filter = ["8-12", "13-17"]

# --- Multiple date ranges (ISO YYYY-MM-DD:YYYY-MM-DD; ranges must not overlap) ---
#
# date_filter = ["2025-01-01:2025-03-31", "2025-04-01:2025-06-30"]

# --- Elevation + monitor — requires user_groups, machine_collections, and applications ---
#
# resource "commander_epm_policy" "elevation_monitor" {
#   policy_name = "Example — elevation monitor"
#   policy_type = "elevation"
#   status      = "monitor"
#
#   user_groups         = ["*"]
#   machine_collections = ["*"]
#   applications        = ["*"]
#
#   # control is optional for monitor / monitor_and_notify
#   # control = ["audit"]
# }

# --- Elevation + monitor_and_notify (same collection rules as monitor) ---
#
# resource "commander_epm_policy" "elevation_monitor_notify" {
#   policy_name = "Example — elevation monitor and notify"
#   policy_type = "elevation"
#   status      = "monitor_and_notify"
#
#   user_groups         = ["user-collection-id-1"]
#   machine_collections = ["machine-collection-id-1"]
#   applications        = ["application-collection-id-1"]
# }

# --- Elevation + off — only required attributes (omit collections and control) ---
#
# resource "commander_epm_policy" "elevation_off" {
#   policy_name = "Example — elevation off"
#   policy_type = "elevation"
#   status      = "off"
# }

# --- File access + enforce — same requirements as elevation + enforce ---
#
# resource "commander_epm_policy" "file_access_enforce" {
#   policy_name = "Example — file access enforce"
#   policy_type = "file_access"
#   status      = "enforce"
#
#   control = ["mfa", "approval"]
#
#   user_groups         = ["*"]
#   machine_collections = ["*"]
#   applications        = ["*"]
# }

# --- File access + monitor — at least one of user_groups / machine_collections / applications ---
#
# resource "commander_epm_policy" "file_access_monitor" {
#   policy_name = "Example — file access monitor"
#   policy_type = "file_access"
#   status      = "monitor"
#
#   user_groups = ["*"]
#   # You may add machine_collections and/or applications when needed
#   # machine_collections = ["machine-collection-id-1"]
#   # applications        = ["application-collection-id-1"]
# }

# --- Command + enforce — control + user_groups + machine_collections; do NOT set applications ---
#
# resource "commander_epm_policy" "command_enforce" {
#   policy_name = "Example — command enforce"
#   policy_type = "command"
#   status      = "enforce"
#
#   control = ["audit", "notify"]
#
#   user_groups         = ["*"]
#   machine_collections = ["*"]
# }

# --- Command + monitor — user_groups + machine_collections only ---
#
# resource "commander_epm_policy" "command_monitor" {
#   policy_name = "Example — command monitor"
#   policy_type = "command"
#   status      = "monitor"
#
#   user_groups         = ["*"]
#   machine_collections = ["*"]
# }

# --- Command + off ---
#
# resource "commander_epm_policy" "command_off" {
#   policy_name = "Example — command off"
#   policy_type = "command"
#   status      = "off"
# }

# --- Least privilege + enforce — only name, type, status (optional machine_collections) ---
#
# resource "commander_epm_policy" "least_privilege_enforce" {
#   policy_name = "Example — least privilege enforce"
#   policy_type = "least_privilege"
#   status      = "enforce"
#
#   # machine_collections = ["*"]
# }

# --- Least privilege + off ---
#
# resource "commander_epm_policy" "least_privilege_off" {
#   policy_name = "Example — least privilege off"
#   policy_type = "least_privilege"
#   status      = "off"
# }

# --- Scoped collections (replace IDs; do not mix "*" with explicit IDs in the same set) ---
#
# user_groups         = ["11111111-1111-1111-1111-111111111111"]
# machine_collections = ["22222222-2222-2222-2222-222222222222"]
# applications        = ["33333333-3333-3333-3333-333333333333"]
