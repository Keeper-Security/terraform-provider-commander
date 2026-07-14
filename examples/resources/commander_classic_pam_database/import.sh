# Import is supported. Import ID is the vault record UID of the PAM database record.
terraform import commander_classic_pam_database.example "AbCdEfGhIjKlMnOpQrStUw"

# Or use the import block in configuration:
# import {
#   to = commander_classic_pam_database.example
#   id = "AbCdEfGhIjKlMnOpQrStUw"
# }
#
# After import, run terraform plan and align configuration with remote state
# (title, hostname_or_ip, database_type, pam_settings, etc.).
