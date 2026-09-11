# Import is supported. The import ID is the vault record UID of the New (NSF) SSH Keys record.

terraform import commander_new_ssh_keys.imported_ssh_keys "Tk835WOSf4LIz1vuyrdfEg"

# Or use an import block in configuration:
# import {
#   to = commander_new_ssh_keys.imported_ssh_keys
#   id = "Tk835WOSf4LIz1vuyrdfEg"
# }
#
# After import, run `terraform plan` and align your configuration with the
# remote state (title, login, hostname, port, public_key, private_key, etc.).
# Sensitive fields are re-read from the vault; set them explicitly in
# configuration if you want Terraform to manage their values going forward.
