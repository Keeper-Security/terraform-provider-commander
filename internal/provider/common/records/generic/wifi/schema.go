// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package wifi

import (
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// SharedAttributes returns the WiFi resource attribute map shared between classic
// and new resources. Callers add any share-extension attribute separately.
func SharedAttributes() map[string]schema.Attribute {
	return utils.MergeResourceAttributes(
		commonrecordsutils.BaseRecordAttributes(),
		map[string]schema.Attribute{
			"ssid":           commonrecordsutils.RequiredSSIDField(),
			"password":       commonrecordsutils.OptionalPasswordField(),
			"encryption":     commonrecordsutils.OptionalWifiEncryptionField(AllowedEncryptions),
			"is_ssid_hidden": commonrecordsutils.OptionalSSIDHiddenField(),
			"custom":         commonrecordsutils.CustomFieldAttributeSchema(),
		},
	)
}

// SharedDataSourceAttributes returns computed WiFi data source attributes shared
// between classic and new data sources. Callers add the lookup key (e.g. wifi)
// and share-extension attributes separately.
func SharedDataSourceAttributes() map[string]dschema.Attribute {
	return utils.MergeDataSourceAttributes(
		commonrecordsutils.DataSourceBaseRecordAttributes(),
		map[string]dschema.Attribute{
			"ssid":           commonrecordsutils.ComputedSSIDField(),
			"password":       commonrecordsutils.ComputedPasswordField(),
			"encryption":     commonrecordsutils.ComputedWifiEncryptionField(),
			"is_ssid_hidden": commonrecordsutils.ComputedSSIDHiddenField(),
			"custom":         commonrecordsutils.CustomFieldDataSourceAttributeSchema(),
		},
	)
}
