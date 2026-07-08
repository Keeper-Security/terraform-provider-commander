// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package wifi

const (
	FlagSSID         = "text.SSID"
	FlagPassword     = "password"
	FlagEncryption   = "wifiEncryption"
	FlagIsSSIDHidden = "isSSIDHidden"
)

// AllowedEncryptions lists the supported wifiEncryption values accepted by Keeper.
var AllowedEncryptions = []string{"wep", "wpa", "noEncryption"}
