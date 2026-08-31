# Import is supported. The import ID is the vault record UID of the login record.
# Find it via: `record-list -t login` in the Commander CLI, or in the Keeper Vault UI.

terraform import commander_classic_login.imported_login "Tk835WOSf4LIz1vuyrdfEg"

# Or use an import block in configuration:
# import {
#   to = commander_classic_login.imported_login
#   id = "Tk835WOSf4LIz1vuyrdfEg"
# }
#
# After import, run `terraform plan` and align your configuration with the
# remote state (title, login, website_address, folder_location, notes, etc.).
# The sensitive `password` field is re-read from the vault and stored in state,
# but set it explicitly in configuration if you want Terraform to manage its
# value going forward.
