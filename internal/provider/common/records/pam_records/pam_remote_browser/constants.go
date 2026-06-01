// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamremotebrowser

// Terraform schema descriptions (plain and Markdown) for registry / docs.
const (
	IDDescription         = "The PAM remote browser record UID assigned by Keeper after create."
	IDMarkdownDescription = "The PAM remote browser record **UID** assigned by Keeper after create."

	TitleDescription         = "Title of the PAM remote browser record."
	TitleMarkdownDescription = "**Title** of the PAM remote browser record."

	URLDescription         = "Target URL for the PAM remote browser session."
	URLMarkdownDescription = "**Target URL** for the PAM remote browser session."

	NotesDescription         = "Optional notes for this PAM remote browser record."
	NotesMarkdownDescription = "Optional **notes** for this PAM remote browser record."

	FolderDescription         = "Folder UID or path to store PAM remote browser record in your Keeper vault. If not provided, the record will be stored in the root path of vault."
	FolderMarkdownDescription = "Folder **UID** or path to store PAM remote browser record in your Keeper vault. If not provided, the record will be stored in the root path of vault."

	PamRemoteBrowserSettingsDescription         = "PAM settings for the PAM remote browser record."
	PamRemoteBrowserSettingsMarkdownDescription = "PAM **settings** for the PAM remote browser record."

	SettingsConfigurationDescription         = "PAM Configuration UID for PAM remote browser settings."
	SettingsConfigurationMarkdownDescription = "**PAM Configuration UID** for PAM remote browser settings."

	SettingsRemoteBrowserIsolationDescription         = "Enable remote browser isolation."
	SettingsRemoteBrowserIsolationMarkdownDescription = "Enable **remote browser isolation**."

	SettingsConnectionsRecordingDescription         = "Manage graphical session recording."
	SettingsConnectionsRecordingMarkdownDescription = "**Manage graphical session recording**."

	SettingsKeyEventsDescription         = "Manage key events for session recording."
	SettingsKeyEventsMarkdownDescription = "**Manage key events for session recording**."

	SettingsAllowURLNavigationDescription         = "Allow navigation via direct URL manipulation."
	SettingsAllowURLNavigationMarkdownDescription = "Allow **navigation** via direct URL manipulation."

	SettingsIgnoreServerCertDescription         = "Ignore Server Certificate."
	SettingsIgnoreServerCertMarkdownDescription = "**Ignore Server Certificate**."

	SettingsAllowedURLsDescription         = "Allowed URL patterns."
	SettingsAllowedURLsMarkdownDescription = "**Allowed URL patterns.**"

	SettingsAllowedResourceURLsDescription         = "Allowed resource URL patterns."
	SettingsAllowedResourceURLsMarkdownDescription = "**Allowed resource URL patterns.**"

	SettingsAutoFillTargetsDescription         = "Browser autofill targets."
	SettingsAutoFillTargetsMarkdownDescription = "**Browser autofill targets.**"

	SettingsAutoFillCredentialsDescription         = "Record UID of Credentials attached to the PAM configuration."
	SettingsAutoFillCredentialsMarkdownDescription = "Record UID of **Credentials** attached to the PAM configuration."

	SettingsAllowCopyDescription         = "Can copy to clipboard."
	SettingsAllowCopyMarkdownDescription = "Can **copy** to clipboard."

	SettingsAllowPasteDescription         = "Can paste from clipboard."
	SettingsAllowPasteMarkdownDescription = "Can **paste** from clipboard."

	SettingsDisableAudioDescription         = "Disable audio."
	SettingsDisableAudioMarkdownDescription = "**Disable audio**."

	SettingsAudioChannelsDescription         = "Number of audio channels; must be 1 for mono or 2 for stereo."
	SettingsAudioChannelsMarkdownDescription = "Number of **audio channels**; must be `1` for **mono** or `2` for **stereo**."

	SettingsAudioBitDepthDescription         = "Audio bit depth; must be 8 for 8-bit or 16 for 16-bit."
	SettingsAudioBitDepthMarkdownDescription = "Audio **bit depth**; must be `8` for **8-bit** or `16` for **16-bit**."

	SettingsAudioSampleRateDescription         = "Audio sample rate in Hz (for example 48000)."
	SettingsAudioSampleRateMarkdownDescription = "Audio **sample rate** in Hz (for example `48000`)."
)

// Commander CLI commands.
const (
	CmdPamRbiEdit = "pam rbi edit"
)

// Commander CLI command flags for `pam rbi edit`.
const (
	FlagRecord                 = "--record"
	FlagConfiguration          = "--configuration"
	FlagRemoteBrowserIsolation = "--remote-browser-isolation"
	FlagConnectionsRecording   = "--connections-recording"
	FlagKeyEvents              = "--key-events"
	FlagAllowURLNavigation     = "--allow-url-navigation"
	FlagIgnoreServerCert       = "--ignore-server-cert"
	FlagAllowedURLs            = "--allowed-urls"
	FlagAllowedResourceURLs    = "--allowed-resource-urls"
	FlagAutofillCredentials    = "--autofill-credentials"
	FlagAutofillTargets        = "--autofill-targets"
	FlagAllowCopy              = "--allow-copy"
	FlagAllowPaste             = "--allow-paste"
	FlagDisableAudio           = "--disable-audio"
	FlagAudioChannels          = "--audio-channels"
	FlagAudioBitDepth          = "--audio-bit-depth"
	FlagAudioSampleRate        = "--audio-sample-rate"
)
