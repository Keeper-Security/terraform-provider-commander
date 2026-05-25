# Import an existing Keeper Drive folder into Terraform state.
#
# The import ID is either the folder UID (preferred) or the folder name.
# Both forms are accepted; the canonical UID is written to state after the
# subsequent Read.
#
# Preferred — by UID:
#
# terraform import commander_new_folder.engineering "E6laPVJ1T3-sWchJCRaWOg"
#
# Alternative — by folder name:
#
# terraform import commander_new_folder.engineering "Engineering / Platform"
#
# Notes:
#   - After import, Terraform runs Read to refresh `name` (and to canonicalize
#     `id` to the UID even if you imported by name).
#   - `share` is NOT populated by import. If you want Terraform to manage
#     shares on this folder, add a `share = { ... }` block to the resource
#     config; the next apply will reconcile it against the API.
#
# Or use the import block (Terraform >= 1.5) directly in configuration:
#
# import {
#   to = commander_new_folder.engineering
#   id = "E6laPVJ1T3-sWchJCRaWOg"
# }
