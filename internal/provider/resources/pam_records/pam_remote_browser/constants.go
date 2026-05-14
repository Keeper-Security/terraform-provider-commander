// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamremotebrowser

// Commander CLI commands.
const (
	CmdPamRbiEdit = "pam rbi edit"
)

// Commander CLI command flags.
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

// Error summaries (first argument to AddError).
const (
	ErrSummaryAddPamRemoteBrowserRecordFailed    = "Failed to add PAM remote browser record"
	ErrSummaryPamRemoteBrowserRecordUpdateFailed = "Failed to update PAM remote browser record"
	ErrSummaryPamRemoteBrowserReadFailed         = "Failed to read PAM remote browser record"
	ErrSummaryPamRbiEditFailed                   = "Failed to update PAM remote browser settings"
)

// Error details operation messages (second argument to ExecuteCommand and AddError; short description for logs).
const (
	ErrDetailAddPamRemoteBrowserRecordFailed    = "Unable to add PAM remote browser record"
	ErrDetailPamRemoteBrowserRecordUpdateFailed = "Unable to update PAM remote browser record"
	ErrDetailPamRemoteBrowserReadFailed         = "Unable to read PAM remote browser record"
	ErrDetailPamRbiEditFailed                   = "Unable to apply PAM remote browser settings"
)
