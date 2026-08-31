// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package sshkeys

import (
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// SharedAttributes returns the sshKeys resource attribute map shared between classic
// and new resources. Callers add any share-extension attribute separately.
func SharedAttributes() map[string]schema.Attribute {
	return utils.MergeResourceAttributes(
		commonrecordsutils.BaseRecordAttributes(),
		map[string]schema.Attribute{
			"login":       commonrecordsutils.OptionalLoginField(),
			"passphrase":  commonrecordsutils.OptionalPassphraseField(),
			"hostname":    commonrecordsutils.OptionalHostnameField(),
			"port":        commonrecordsutils.OptionalPortField(),
			"public_key":  commonrecordsutils.OptionalPublicKeyField(),
			"private_key": commonrecordsutils.OptionalPrivateKeyField(),
			"custom":      commonrecordsutils.CustomFieldAttributeSchema(),
		},
	)
}

// SharedDataSourceAttributes returns computed sshKeys data source attributes shared
// between classic and new data sources.
func SharedDataSourceAttributes() map[string]dschema.Attribute {
	return utils.MergeDataSourceAttributes(
		commonrecordsutils.DataSourceBaseRecordAttributes(),
		map[string]dschema.Attribute{
			"login":       commonrecordsutils.ComputedLoginField(),
			"passphrase":  commonrecordsutils.ComputedPassphraseField(),
			"hostname":    commonrecordsutils.ComputedHostnameField(),
			"port":        commonrecordsutils.ComputedPortField(),
			"public_key":  commonrecordsutils.ComputedPublicKeyField(),
			"private_key": commonrecordsutils.ComputedPrivateKeyField(),
			"custom":      commonrecordsutils.CustomFieldDataSourceAttributeSchema(),
		},
	)
}
