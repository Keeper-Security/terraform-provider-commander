resource "commander_new_ssh_keys" "example" {
  title = "Production SSH Key"

  login       = "alice@xyz.com"
  passphrase  = "ExamplePassphrase123!"
  hostname    = "12.0.0.1"
  port        = "22"
  public_key  = "ssh-rsa AAAA..."
  private_key = "ssh-rsa AAAA..."

  notes           = "SSH key for production bastion access."
  folder_location = "_REPLACE_WITH_FOLDER_PATH_OR_UID_"

  share = {
    "alice@example.com" = "full-manager"
  }
}
