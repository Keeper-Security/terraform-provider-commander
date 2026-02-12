# Look up a managed company by name or ID. Returns id, name, node, plan, file_plan.

data "commander_manage_company" "example" {
  managed_company = "Acme Corp"
}

output "managed_company_id" {
  value = data.commander_manage_company.example.id
}

output "managed_company_name" {
  value = data.commander_manage_company.example.name
}

output "managed_company_plan" {
  value = data.commander_manage_company.example.plan
}

output "managed_company_file_plan" {
  value = data.commander_manage_company.example.file_plan
}
