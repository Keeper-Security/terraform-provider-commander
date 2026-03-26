# Import is supported. Use either the user email/ID only (logged-in account) or comma-separated managed_company,user.
# User only (logged-in account):
terraform import commander_enterprise_user.example "user@example.com"
terraform import commander_enterprise_user.example 1326075447607317
# With managed company (managed_company_name_or_id,user_email_or_id):
terraform import commander_enterprise_user.example "Test Company,1326075447607317"
terraform import commander_enterprise_user.example "1169425105420462,user@example.com"

# Or use the import block in configuration:
# import {
#   to = commander_enterprise_user.example
#   id = "user@example.com"
# }
# import {
#   to = commander_enterprise_user.example
#   id = "Test Company,1326075447607317"
# }
