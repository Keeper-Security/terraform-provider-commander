# Import by classic shared folder UID (required for correct state; Read refreshes name and permissions from Commander).
# Get the UID from Commander or from a prior data.commander_shared_folder.*.id / resource id.
#
# terraform import commander_shared_folder.example "BTbjhOmqw9iYal3OQJ9UAQ"
#
# Or use the import block in configuration:
# import {
#   to = commander_shared_folder.example
#   id = "BTbjhOmqw9iYal3OQJ9UAQ"
# }
