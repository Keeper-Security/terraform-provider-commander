// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package bankaccount

import (
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// SharedAttributes returns the Bank Account resource attribute map shared between
// classic and new resources. Callers add any share-extension attribute separately.
func SharedAttributes() map[string]schema.Attribute {
	return utils.MergeResourceAttributes(
		commonrecordsutils.BaseRecordAttributes(),
		map[string]schema.Attribute{
			"bank_account":    commonrecordsutils.BankAccountNestedSchema(),
			"name":            commonrecordsutils.NameNestedSchema(true),
			"login":           commonrecordsutils.OptionalLoginField(),
			"password":        commonrecordsutils.OptionalPasswordField(),
			"website_address": commonrecordsutils.OptionalWebsiteAddressField(),
			"card_ref":        commonrecordsutils.OptionalCardRefField(),
			"custom":          commonrecordsutils.CustomFieldAttributeSchema(),
		},
	)
}

// SharedDataSourceAttributes returns computed Bank Account data source attributes
// shared between classic and new data sources. Callers add the lookup key (e.g.
// account) and share-extension attributes separately.
func SharedDataSourceAttributes() map[string]dschema.Attribute {
	return utils.MergeDataSourceAttributes(
		commonrecordsutils.DataSourceBaseRecordAttributes(),
		map[string]dschema.Attribute{
			"bank_account":    commonrecordsutils.BankAccountDataSourceNestedSchema(),
			"name":            commonrecordsutils.NameDataSourceNestedSchema(),
			"login":           commonrecordsutils.ComputedLoginField(),
			"password":        commonrecordsutils.ComputedPasswordField(),
			"website_address": commonrecordsutils.ComputedWebsiteAddressField(),
			"card_ref":        commonrecordsutils.ComputedCardRefField(),
			"custom":          commonrecordsutils.CustomFieldDataSourceAttributeSchema(),
		},
	)
}
