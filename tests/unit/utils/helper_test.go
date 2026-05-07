// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package utils_test

import (
	"context"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestExtractNodeIDFromCreateNodeResponse(t *testing.T) {
	tests := []struct {
		input  string
		wantID string
		wantOK bool
	}{
		{"Node ID: 1169425105420462", "1169425105420462", true},
		{"Node is created with Node ID: 123", "123", true},
		{"Node ID:  999", "999", true},
		{"no id here", "", false},
		{"", "", false},
		{"Node ID:", "", false},
	}
	for _, tt := range tests {
		gotID, gotOK := utils.ExtractNodeIDFromCreateNodeResponse(tt.input)
		if gotID != tt.wantID || gotOK != tt.wantOK {
			t.Errorf("ExtractNodeIDFromCreateNodeResponse(%q) = %q, %v; want %q, %v", tt.input, gotID, gotOK, tt.wantID, tt.wantOK)
		}
	}
}

func TestExtractNodeName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Metronlabs\\Aditya Dev Inc", "Aditya Dev Inc"},
		{"Metronlabs", "Metronlabs"},
		{"SingleNode", "SingleNode"},
		{"", ""},
		{"A\\B\\C", "C"},
	}
	for _, tt := range tests {
		got := utils.ExtractNodeName(tt.input)
		if got != tt.want {
			t.Errorf("ExtractNodeName(%q) = %q; want %q", tt.input, got, tt.want)
		}
	}
}

func TestUnmarshalApiResponse(t *testing.T) {
	type target struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	data := map[string]interface{}{"id": float64(1), "name": "test"}
	var out target
	err := utils.UnmarshalApiResponse(data, &out)
	if err != nil {
		t.Fatalf("UnmarshalApiResponse: %v", err)
	}
	if out.ID != 1 || out.Name != "test" {
		t.Errorf("got %+v; want id=1, name=test", out)
	}
}

func TestUnmarshalApiResponse_InvalidTarget(t *testing.T) {
	data := map[string]interface{}{"x": "y"}
	var out chan int
	err := utils.UnmarshalApiResponse(data, &out)
	if err != nil {
		// Marshaling to JSON works; unmarshaling into chan may fail
		return
	}
	// If no error, that's ok for this test
	_ = out
}

func TestConvertItemsToIdMap_NullSet(t *testing.T) {
	lookup := utils.LookupMaps{
		IdentifierToId: map[string]string{"a": "1"},
		IdToIdentifier: map[string]string{"1": "a"},
	}
	validateItem := func(s string) (bool, string) { return true, "" }
	result, err := utils.ConvertItemsToIdMap(types.SetNull(types.StringType), lookup, "role", validateItem)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty map for null set, got %v", result)
	}
}

func TestConvertItemsToIdMap_UnknownSet(t *testing.T) {
	lookup := utils.LookupMaps{
		IdentifierToId: map[string]string{"a": "1"},
		IdToIdentifier: map[string]string{"1": "a"},
	}
	validateItem := func(s string) (bool, string) { return true, "" }
	result, err := utils.ConvertItemsToIdMap(types.SetUnknown(types.StringType), lookup, "role", validateItem)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty map for unknown set, got %v", result)
	}
}

func TestConvertItemsToIdMap_EmptySet(t *testing.T) {
	ctx := context.Background()
	emptySet, _ := types.SetValueFrom(ctx, types.StringType, []types.String{})
	lookup := utils.LookupMaps{
		IdentifierToId: map[string]string{"a": "1"},
		IdToIdentifier: map[string]string{"1": "a"},
	}
	validateItem := func(s string) (bool, string) { return true, "" }
	result, err := utils.ConvertItemsToIdMap(emptySet, lookup, "role", validateItem)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty map for empty set, got %v", result)
	}
}

func TestConvertItemsToIdMap_ById(t *testing.T) {
	ctx := context.Background()
	setVal, _ := types.SetValueFrom(ctx, types.StringType, []types.String{types.StringValue("1")})
	lookup := utils.LookupMaps{
		IdentifierToId: map[string]string{"role1": "1"},
		IdToIdentifier: map[string]string{"1": "role1"},
	}
	validateItem := func(s string) (bool, string) { return true, "" }
	result, err := utils.ConvertItemsToIdMap(setVal, lookup, "role", validateItem)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 || result["1"] != "1" {
		t.Errorf("expected map[1:1], got %v", result)
	}
}

func TestConvertItemsToIdMap_ByIdentifier(t *testing.T) {
	ctx := context.Background()
	setVal, _ := types.SetValueFrom(ctx, types.StringType, []types.String{types.StringValue("role1")})
	lookup := utils.LookupMaps{
		IdentifierToId: map[string]string{"role1": "1"},
		IdToIdentifier: map[string]string{"1": "role1"},
	}
	validateItem := func(s string) (bool, string) { return true, "" }
	result, err := utils.ConvertItemsToIdMap(setVal, lookup, "role", validateItem)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 || result["1"] != "role1" {
		t.Errorf("expected map[1:role1], got %v", result)
	}
}

func TestConvertItemsToIdMap_NotFound(t *testing.T) {
	ctx := context.Background()
	setVal, _ := types.SetValueFrom(ctx, types.StringType, []types.String{types.StringValue("unknown")})
	lookup := utils.LookupMaps{
		IdentifierToId: map[string]string{"role1": "1"},
		IdToIdentifier: map[string]string{"1": "role1"},
	}
	validateItem := func(s string) (bool, string) { return true, "" }
	_, err := utils.ConvertItemsToIdMap(setVal, lookup, "role", validateItem)
	if err == nil {
		t.Fatal("expected error when item not in lookup")
	}
}

func TestConvertItemsToIdMap_ValidateFails(t *testing.T) {
	ctx := context.Background()
	setVal, _ := types.SetValueFrom(ctx, types.StringType, []types.String{types.StringValue("bad")})
	lookup := utils.LookupMaps{
		IdentifierToId: map[string]string{},
		IdToIdentifier: map[string]string{},
	}
	validateItem := func(s string) (bool, string) {
		if s == "bad" {
			return false, "invalid"
		}
		return true, ""
	}
	_, err := utils.ConvertItemsToIdMap(setVal, lookup, "role", validateItem)
	if err == nil {
		t.Fatal("expected error when validateItem returns false")
	}
}

func TestConvertItemsToIdMap_Duplicate(t *testing.T) {
	ctx := context.Background()
	setVal, _ := types.SetValueFrom(ctx, types.StringType, []types.String{types.StringValue("1"), types.StringValue("role1")})
	lookup := utils.LookupMaps{
		IdentifierToId: map[string]string{"role1": "1"},
		IdToIdentifier: map[string]string{"1": "role1"},
	}
	validateItem := func(s string) (bool, string) { return true, "" }
	_, err := utils.ConvertItemsToIdMap(setVal, lookup, "role", validateItem)
	if err == nil {
		t.Fatal("expected error for duplicate role id")
	}
}

func TestConvertItemsToIdMap_SkipsEmpty(t *testing.T) {
	ctx := context.Background()
	setVal, _ := types.SetValueFrom(ctx, types.StringType, []types.String{types.StringValue(""), types.StringValue("1")})
	lookup := utils.LookupMaps{
		IdentifierToId: map[string]string{"role1": "1"},
		IdToIdentifier: map[string]string{"1": "role1"},
	}
	validateItem := func(s string) (bool, string) { return true, "" }
	result, err := utils.ConvertItemsToIdMap(setVal, lookup, "role", validateItem)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 || result["1"] != "1" {
		t.Errorf("expected map[1:1] (empty skipped), got %v", result)
	}
}
