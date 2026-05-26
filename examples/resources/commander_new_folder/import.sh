# Import an existing Nested Shared Folder into Terraform state.
#
# The import ID is either the folder UID (preferred) or the folder name.
# Both forms are accepted; the canonical UID is written to state after the
# subsequent Read.
#
# Preferred — by UID:
#
# terraform import commander_new_folder.engineering "Cuuc9aK6VuATH49ewBf0zg"
#
# Alternative — by folder name:
#
# terraform import commander_new_folder.engineering "Engineering / Platform"
# Or use the import block (Terraform >= 1.5) directly in configuration:
#
# import {
#   to = commander_new_folder.engineering
#   id = "Cuuc9aK6VuATH49ewBf0zg"
# }
