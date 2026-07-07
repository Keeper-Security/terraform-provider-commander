// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package sshkeys

const (
	FlagLogin      = "login"
	FlagPassphrase = "f.password.passphrase"
	FlagHost       = "f.host"
	FlagKeyPair    = "f.keyPair"

	PassphraseFieldLabel = "passphrase"
)

const (
	LoginDescription         = "SSH login"
	LoginMarkdownDescription = "SSH **login**"

	PassphraseDescription         = "Passphrase"
	PassphraseMarkdownDescription = "**Passphrase**"

	HostnameDescription         = "SSH Hostname or IP address."
	HostnameMarkdownDescription = "SSH **Hostname** or **IP address**."

	PortDescription         = "SSH port."
	PortMarkdownDescription = "SSH **port**."

	PublicKeyDescription         = "Public key."
	PublicKeyMarkdownDescription = "**Public key**"

	PrivateKeyDescription         = "Private key."
	PrivateKeyMarkdownDescription = "**Private key**."

	DSNotesDescription         = "Notes on the record, if any."
	DSNotesMarkdownDescription = "**Notes** on the record, if any."

	DSFolderDescription         = "Folder path where the record is stored."
	DSFolderMarkdownDescription = "**Folder path** where the record is stored."
)
