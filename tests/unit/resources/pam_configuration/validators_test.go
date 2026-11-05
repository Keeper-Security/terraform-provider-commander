// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamconfiguration_test

import (
	"context"
	"testing"

	pamconfiguration "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/pam_configuration"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// nestedObjectType returns the tftypes.Object type of a top-level block/attribute
// in the real resource schema, so test values match it exactly.
func nestedObjectType(t *testing.T, sch rschema.Schema, name string) tftypes.Object {
	t.Helper()
	objType, ok := sch.Type().TerraformType(context.Background()).(tftypes.Object)
	if !ok {
		t.Fatalf("expected tftypes.Object, got %T", sch.Type().TerraformType(context.Background()))
	}
	nested, ok := objType.AttributeTypes[name].(tftypes.Object)
	if !ok {
		t.Fatalf("expected %q to be an object type", name)
	}
	return nested
}

// ussObjectType returns the tftypes.Object type of the universal_secrets_sync block nested
// inside a cloud environment block's type (aws, azure, or gcp).
func ussObjectType(t *testing.T, envType tftypes.Object) tftypes.Object {
	t.Helper()
	nested, ok := envType.AttributeTypes["universal_secrets_sync"].(tftypes.Object)
	if !ok {
		t.Fatal("expected universal_secrets_sync to be an object type")
	}
	return nested
}

// nullObject builds a null-valued tftypes.Value for the given object type, with
// explicit field values overridden by vals.
func nullObject(objType tftypes.Object, vals map[string]tftypes.Value) tftypes.Value {
	out := map[string]tftypes.Value{}
	for name, attrType := range objType.AttributeTypes {
		if v, ok := vals[name]; ok {
			out[name] = v
			continue
		}
		out[name] = tftypes.NewValue(attrType, nil)
	}
	return tftypes.NewValue(objType, out)
}

// newValidateConfigRaw builds a full config tftypes.Value matching the real resource
// schema, with every attribute null except the ones in overrides.
func newValidateConfigRaw(t *testing.T, sch rschema.Schema, overrides map[string]tftypes.Value) tftypes.Value {
	t.Helper()
	objType, ok := sch.Type().TerraformType(context.Background()).(tftypes.Object)
	if !ok {
		t.Fatalf("expected tftypes.Object, got %T", sch.Type().TerraformType(context.Background()))
	}
	return nullObject(objType, overrides)
}

func runEnvironmentBlocksValidator(t *testing.T, sch rschema.Schema, raw tftypes.Value) resource.ValidateConfigResponse {
	t.Helper()
	r := pamconfiguration.NewPamConfigurationResource().(*pamconfiguration.PamConfigurationResource)
	req := resource.ValidateConfigRequest{Config: tfsdk.Config{Schema: sch, Raw: raw}}
	var resp resource.ValidateConfigResponse
	for _, v := range r.ConfigValidators(context.Background()) {
		v.ValidateResource(context.Background(), req, &resp)
	}
	return resp
}

func stringVal(s string) tftypes.Value { return tftypes.NewValue(tftypes.String, s) }
func boolVal(b bool) tftypes.Value     { return tftypes.NewValue(tftypes.Bool, b) }

// TestSchema_UniversalSecretsSyncNestedUnderCloudBlocksOnly confirms the schema shape
// directly: universal_secrets_sync exists under aws/azure/gcp (azure's variant also has
// vault_name, aws/gcp's don't), and domain/local_network have no such block at all - so
// these invalid combinations can't even be expressed in config, let alone need runtime
// validation.
func TestSchema_UniversalSecretsSyncNestedUnderCloudBlocksOnly(t *testing.T) {
	sch := getResSchema(t)

	awsUSS := ussObjectType(t, nestedObjectType(t, sch, "aws"))
	if _, ok := awsUSS.AttributeTypes["vault_name"]; ok {
		t.Error("expected aws.universal_secrets_sync to have no vault_name attribute")
	}
	if _, ok := awsUSS.AttributeTypes["enabled"]; !ok {
		t.Error("expected aws.universal_secrets_sync to have an enabled attribute")
	}

	gcpUSS := ussObjectType(t, nestedObjectType(t, sch, "gcp"))
	if _, ok := gcpUSS.AttributeTypes["vault_name"]; ok {
		t.Error("expected gcp.universal_secrets_sync to have no vault_name attribute")
	}

	azureUSS := ussObjectType(t, nestedObjectType(t, sch, "azure"))
	if _, ok := azureUSS.AttributeTypes["vault_name"]; !ok {
		t.Error("expected azure.universal_secrets_sync to have a vault_name attribute")
	}

	domainType := nestedObjectType(t, sch, "domain")
	if _, ok := domainType.AttributeTypes["universal_secrets_sync"]; ok {
		t.Error("expected domain to have no universal_secrets_sync block")
	}

	localType := nestedObjectType(t, sch, "local_network")
	if _, ok := localType.AttributeTypes["universal_secrets_sync"]; ok {
		t.Error("expected local_network to have no universal_secrets_sync block")
	}
}

// azureUniversalSecretsSyncObjectType is the attr.Type map for the
// azure.universal_secrets_sync block, matching commonUniversalSecretsSyncSchema() plus
// vault_name.
var azureUniversalSecretsSyncObjectType = map[string]attr.Type{
	"enabled":         types.BoolType,
	"folder_location": types.SetType{ElemType: types.StringType},
	"dry_run":         types.BoolType,
	"sync_identity":   types.StringType,
	"vault_name":      types.StringType,
}

func TestAzureUniversalSecretsSyncRequiredFieldsValidator_RequiresVaultName(t *testing.T) {
	v := pamconfiguration.AzureUniversalSecretsSyncRequiredFieldsValidator()
	obj, diags := types.ObjectValue(azureUniversalSecretsSyncObjectType, map[string]attr.Value{
		"enabled":         types.BoolValue(true),
		"folder_location": types.SetNull(types.StringType),
		"dry_run":         types.BoolValue(false),
		"sync_identity":   types.StringNull(),
		"vault_name":      types.StringNull(),
	})
	if diags.HasError() {
		t.Fatalf("building object: %v", diags)
	}

	var resp validator.ObjectResponse
	v.ValidateObject(context.Background(), validator.ObjectRequest{
		Path:        path.Root("azure").AtName("universal_secrets_sync"),
		ConfigValue: obj,
	}, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected an error when vault_name is null on a present universal_secrets_sync block")
	}
}

func TestAzureUniversalSecretsSyncRequiredFieldsValidator_AllowsVaultName(t *testing.T) {
	v := pamconfiguration.AzureUniversalSecretsSyncRequiredFieldsValidator()
	obj, diags := types.ObjectValue(azureUniversalSecretsSyncObjectType, map[string]attr.Value{
		"enabled":         types.BoolValue(true),
		"folder_location": types.SetNull(types.StringType),
		"dry_run":         types.BoolValue(false),
		"sync_identity":   types.StringNull(),
		"vault_name":      types.StringValue("my-vault"),
	})
	if diags.HasError() {
		t.Fatalf("building object: %v", diags)
	}

	var resp validator.ObjectResponse
	v.ValidateObject(context.Background(), validator.ObjectRequest{
		Path:        path.Root("azure").AtName("universal_secrets_sync"),
		ConfigValue: obj,
	}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("expected no error when vault_name is set, got %v", resp.Diagnostics)
	}
}

func TestAzureUniversalSecretsSyncRequiredFieldsValidator_AllowsAbsentBlock(t *testing.T) {
	v := pamconfiguration.AzureUniversalSecretsSyncRequiredFieldsValidator()

	var resp validator.ObjectResponse
	v.ValidateObject(context.Background(), validator.ObjectRequest{
		Path:        path.Root("azure").AtName("universal_secrets_sync"),
		ConfigValue: types.ObjectNull(azureUniversalSecretsSyncObjectType),
	}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("expected no error when universal_secrets_sync block is absent, got %v", resp.Diagnostics)
	}
}

func TestEnvironmentBlocksValidator_Azure_AllowsVaultName(t *testing.T) {
	sch := getResSchema(t)
	azureType := nestedObjectType(t, sch, "azure")
	ussType := ussObjectType(t, azureType)

	raw := newValidateConfigRaw(t, sch, map[string]tftypes.Value{
		"environment": stringVal("azure"),
		"azure": nullObject(azureType, map[string]tftypes.Value{
			"azure_id":        stringVal("az-1"),
			"client_id":       stringVal("client-1"),
			"client_secret":   stringVal("secret-1"),
			"subscription_id": stringVal("sub-1"),
			"tenant_id":       stringVal("tenant-1"),
			"universal_secrets_sync": nullObject(ussType, map[string]tftypes.Value{
				"enabled":    boolVal(true),
				"vault_name": stringVal("my-vault"),
			}),
		}),
	})

	resp := runEnvironmentBlocksValidator(t, sch, raw)
	if resp.Diagnostics.HasError() {
		t.Fatalf("expected no error when azure.universal_secrets_sync.vault_name is set, got %v", resp.Diagnostics)
	}
}

func TestEnvironmentBlocksValidator_Azure_AllowsAbsentUniversalSecretsSync(t *testing.T) {
	sch := getResSchema(t)
	azureType := nestedObjectType(t, sch, "azure")

	raw := newValidateConfigRaw(t, sch, map[string]tftypes.Value{
		"environment": stringVal("azure"),
		"azure": nullObject(azureType, map[string]tftypes.Value{
			"azure_id":        stringVal("az-1"),
			"client_id":       stringVal("client-1"),
			"client_secret":   stringVal("secret-1"),
			"subscription_id": stringVal("sub-1"),
			"tenant_id":       stringVal("tenant-1"),
		}),
	})

	resp := runEnvironmentBlocksValidator(t, sch, raw)
	if resp.Diagnostics.HasError() {
		t.Fatalf("expected no error when azure.universal_secrets_sync is entirely absent, got %v", resp.Diagnostics)
	}
}

func TestEnvironmentBlocksValidator_Aws_AllowsUniversalSecretsSyncWithoutVaultName(t *testing.T) {
	sch := getResSchema(t)
	awsType := nestedObjectType(t, sch, "aws")
	ussType := ussObjectType(t, awsType)

	raw := newValidateConfigRaw(t, sch, map[string]tftypes.Value{
		"environment": stringVal("aws"),
		"aws": nullObject(awsType, map[string]tftypes.Value{
			"aws_id": stringVal("aws-1"),
			"universal_secrets_sync": nullObject(ussType, map[string]tftypes.Value{
				"enabled": boolVal(true),
			}),
		}),
	})

	resp := runEnvironmentBlocksValidator(t, sch, raw)
	if resp.Diagnostics.HasError() {
		t.Fatalf("expected no error for aws.universal_secrets_sync without vault_name, got %v", resp.Diagnostics)
	}
}

func TestUniversalSecretsSyncNotEmptyValidator(t *testing.T) {
	v := pamconfiguration.UniversalSecretsSyncNotEmptyValidator()
	allNull := map[string]attr.Value{
		"enabled":         types.BoolNull(),
		"folder_location": types.SetNull(types.StringType),
		"dry_run":         types.BoolNull(),
		"sync_identity":   types.StringNull(),
		"vault_name":      types.StringNull(),
	}
	with := func(name string, val attr.Value) map[string]attr.Value {
		m := map[string]attr.Value{}
		for k, v := range allNull {
			m[k] = v
		}
		m[name] = val
		return m
	}
	mustObj := func(vals map[string]attr.Value) types.Object {
		obj, diags := types.ObjectValue(azureUniversalSecretsSyncObjectType, vals)
		if diags.HasError() {
			t.Fatalf("building object: %v", diags)
		}
		return obj
	}

	tests := []struct {
		name      string
		value     types.Object
		wantError bool
	}{
		{"empty block", mustObj(allNull), true},
		{"enabled set", mustObj(with("enabled", types.BoolValue(false))), false},
		{"only dry_run set", mustObj(with("dry_run", types.BoolValue(true))), false},
		{"only vault_name set", mustObj(with("vault_name", types.StringValue("kv"))), false},
		{"unknown attribute counts as set", mustObj(with("sync_identity", types.StringUnknown())), false},
		{"absent block", types.ObjectNull(azureUniversalSecretsSyncObjectType), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var resp validator.ObjectResponse
			v.ValidateObject(context.Background(), validator.ObjectRequest{
				Path:        path.Root("aws").AtName("universal_secrets_sync"),
				ConfigValue: tt.value,
			}, &resp)
			if got := resp.Diagnostics.HasError(); got != tt.wantError {
				t.Errorf("HasError = %v, want %v; diagnostics: %v", got, tt.wantError, resp.Diagnostics)
			}
		})
	}
}
