# Provider configuration for Commander Service Mode.
# See: https://docs.keeper.io/en/keeperpam/commander-cli/service-mode-rest-api

provider "commander" {
  service_mode_url     = "https://your-commander-service.example.com"
  service_mode_api_key = "your-api-key"
  timeout              = 60 # optional; defaults to 60 seconds for HTTP Client and API Call timeout
}
