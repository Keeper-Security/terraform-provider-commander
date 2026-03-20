# Import is supported. Use either the team name/UID only (logged-in account) or comma-separated managed_company,team.
# Team only (logged-in account):
terraform import commander_enterprise_team.example "Engineering"
terraform import commander_enterprise_team.example 1234567890123456
# With managed company (managed_company_name_or_id,team_name_or_uid):
terraform import commander_enterprise_team.example "Test Company,Engineering"
terraform import commander_enterprise_team.example "1169425105420462,1234567890123456"

# Or use the import block in configuration:
# import {
#   to = commander_enterprise_team.example
#   id = "Engineering"
# }
# import {
#   to = commander_enterprise_team.example
#   id = "Test Company,Engineering"
# }
