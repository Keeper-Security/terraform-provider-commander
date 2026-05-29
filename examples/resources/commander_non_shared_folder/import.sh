# Import by folder UID (required for correct state; Read refreshes name from Commander).
# Get the UID from Commander or from a prior data.commander_non_shared_folder.*.id / resource id.
#
# terraform import commander_non_shared_folder.example "NJiANrRnbuvVEOgnqjiYaw"
#
# Or use the import block in configuration:
# import {
#   to = commander_non_shared_folder.example
#   id = "NJiANrRnbuvVEOgnqjiYaw"
# }