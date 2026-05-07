// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package utils_test

import (
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
)

func TestValidateRecordIdentifiersAgainstList(t *testing.T) {
	entries := []utils.VaultRecordListEntry{
		{RecordUID: "QJz-Nl64lQQhLWX2TaoNag", Title: "Remote Browser", Type: "pamRemoteBrowser"},
		{RecordUID: "zOfPoHpn6EcTqyeB8AJQhQ", Title: "MySQL Admin User", Type: "pamUser"},
	}

	t.Run("empty identifiers", func(t *testing.T) {
		if err := utils.ValidateRecordIdentifiersAgainstList(nil, entries); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("valid by uid", func(t *testing.T) {
		err := utils.ValidateRecordIdentifiersAgainstList([]string{"QJz-Nl64lQQhLWX2TaoNag"}, entries)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("valid by title", func(t *testing.T) {
		err := utils.ValidateRecordIdentifiersAgainstList([]string{"MySQL Admin User"}, entries)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		err := utils.ValidateRecordIdentifiersAgainstList([]string{"nope", "also-bad"}, entries)
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("dedupes same key", func(t *testing.T) {
		err := utils.ValidateRecordIdentifiersAgainstList([]string{"nope", "nope"}, entries)
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("skips blank", func(t *testing.T) {
		err := utils.ValidateRecordIdentifiersAgainstList([]string{"  ", "\t"}, entries)
		if err != nil {
			t.Fatal(err)
		}
	})
}
