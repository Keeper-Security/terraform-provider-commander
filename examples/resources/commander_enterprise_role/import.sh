# Import is supported. Use either the role name/ID only (logged-in account) or comma-separated managed_company,role.
# Role only (logged-in account):
terraform import commander_enterprise_role.example "Admin"
terraform import commander_enterprise_role.example 1234567890
# With managed company (managed_company_name_or_id,role_name_or_id):
terraform import commander_enterprise_role.example "Test Company,Admin"
terraform import commander_enterprise_role.example "1169425105420462,1234567890"

# Or use the import block in configuration:
# import {
#   to = commander_enterprise_role.example
#   id = "Admin"
# }
# import {
#   to = commander_enterprise_role.example
#   id = "Test Company,Admin"
# }
