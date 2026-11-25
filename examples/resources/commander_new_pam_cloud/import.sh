# Import is supported. Import ID is the vault record UID of the PAM cloud record.
terraform import commander_new_pam_cloud.full_fields "AbCdEfGhIjKlMnOpQrStUw"

# Or use the import block in configuration:
# import {
#   to = commander_new_pam_cloud.full_fields
#   id = "AbCdEfGhIjKlMnOpQrStUw"
# }
#
# After import, run terraform plan and align configuration with remote state
# (title, url, account_identifier, notes, share, etc.).
