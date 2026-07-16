# Import is supported. The import ID is the vault record UID of the pamUser record.
# Find it via: `record-list -t pamUser` in the Commander CLI, or in the Keeper Vault UI.

terraform import commander_new_pam_user.imported_pam_user "KshAnhsoDL4Q1hA9akk4sg"

# Or use an import block in configuration:
# import {
#   to = commander_new_pam_user.imported_pam_user
#   id = "KshAnhsoDL4Q1hA9akk4sg"
# }
#
# After import, run `terraform plan` and align your configuration with the
# remote state (title, login, folder, rotation_settings, share, etc.).
# Sensitive fields like `password` and `private_pem_key` are not re-emitted;
# set them explicitly in configuration if you want Terraform to manage them.
