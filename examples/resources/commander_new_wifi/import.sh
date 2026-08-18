# Import is supported. The import ID is the vault record UID of the wifiCredentials record.
# Find it via: `record-list -t wifiCredentials` in the Commander CLI, or in the Keeper Vault UI.

terraform import commander_new_wifi.imported_wifi "Tk835WOSf4LIz1vuyrdfEg"

# Or use an import block in configuration:
# import {
#   to = commander_new_wifi.imported_wifi
#   id = "Tk835WOSf4LIz1vuyrdfEg"
# }
#
# After import, run `terraform plan` and align your configuration with the
# remote state (title, ssid, encryption, is_ssid_hidden, folder, etc.).
# The sensitive `password` field is re-read from the vault and stored in
# state, but be sure to set it explicitly in configuration if you want
# Terraform to manage its value going forward.
