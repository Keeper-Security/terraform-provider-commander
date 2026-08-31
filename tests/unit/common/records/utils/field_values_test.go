// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package utils_test

import (
	"encoding/json"
	"testing"

	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestStringUnset(t *testing.T) {
	t.Parallel()

	if !commonrecordsutils.StringUnset(types.StringNull()) {
		t.Fatal("expected null to be unset")
	}
	if !commonrecordsutils.StringUnset(types.StringValue("  ")) {
		t.Fatal("expected blank to be unset")
	}
	if commonrecordsutils.StringUnset(types.StringValue("x")) {
		t.Fatal("expected non-blank to be set")
	}
}

func TestHostValue_IsEmpty(t *testing.T) {
	t.Parallel()

	var nilHost *commonrecordsutils.HostValue
	if !nilHost.IsEmpty() {
		t.Fatal("expected nil host to be empty")
	}

	empty := &commonrecordsutils.HostValue{
		HostName: types.StringNull(),
		Port:     types.StringNull(),
	}
	if !empty.IsEmpty() {
		t.Fatal("expected null host fields to be empty")
	}

	withHost := &commonrecordsutils.HostValue{
		HostName: types.StringValue("localhost"),
		Port:     types.StringNull(),
	}
	if withHost.IsEmpty() {
		t.Fatal("expected hostname-only host to be non-empty")
	}
}

func TestHostFlatToJSON(t *testing.T) {
	t.Parallel()

	if _, ok := commonrecordsutils.HostFlatToJSON(types.StringNull(), types.StringNull()); ok {
		t.Fatal("expected empty host to return ok=false")
	}

	j, ok := commonrecordsutils.HostFlatToJSON(types.StringValue("localhost"), types.StringValue("22"))
	if !ok {
		t.Fatal("expected host JSON")
	}
	if j != `{"hostName":"localhost","port":"22"}` {
		t.Fatalf("host JSON = %q", j)
	}
}

func TestKeyPairFlatToJSON(t *testing.T) {
	t.Parallel()

	if _, ok := commonrecordsutils.KeyPairFlatToJSON(types.StringNull(), types.StringNull()); ok {
		t.Fatal("expected empty key pair to return ok=false")
	}

	j, ok := commonrecordsutils.KeyPairFlatToJSON(types.StringValue("pub"), types.StringValue("priv"))
	if !ok {
		t.Fatal("expected key pair JSON")
	}
	if j != `{"publicKey":"pub","privateKey":"priv"}` {
		t.Fatalf("key pair JSON = %q", j)
	}
}

func TestFlatHostFromFields(t *testing.T) {
	t.Parallel()

	fields := []utils.VaultRecordFieldResponse{
		{Type: "host", Value: json.RawMessage(`[{"hostName":"localhost","port":"22"}]`)},
	}
	hostname, port := commonrecordsutils.FlatHostFromFields(fields)
	if hostname.ValueString() != "localhost" || port.ValueString() != "22" {
		t.Fatalf("host = %q:%q", hostname.ValueString(), port.ValueString())
	}

	hostname, port = commonrecordsutils.FlatHostFromFields(nil)
	if !hostname.IsNull() || !port.IsNull() {
		t.Fatalf("expected null host, got %q:%q", hostname.ValueString(), port.ValueString())
	}
}

func TestFlatKeyPairFromFields(t *testing.T) {
	t.Parallel()

	fields := []utils.VaultRecordFieldResponse{
		{Type: "keyPair", Value: json.RawMessage(`[{"publicKey":"pub","privateKey":"priv"}]`)},
	}
	publicKey, privateKey := commonrecordsutils.FlatKeyPairFromFields(fields)
	if publicKey.ValueString() != "pub" || privateKey.ValueString() != "priv" {
		t.Fatalf("keys = %q:%q", publicKey.ValueString(), privateKey.ValueString())
	}

	publicKey, privateKey = commonrecordsutils.FlatKeyPairFromFields(nil)
	if !publicKey.IsNull() || !privateKey.IsNull() {
		t.Fatalf("expected null keys, got %q:%q", publicKey.ValueString(), privateKey.ValueString())
	}
}

func TestKeyPairValue_IsEmpty(t *testing.T) {
	t.Parallel()

	var nilKP *commonrecordsutils.KeyPairValue
	if !nilKP.IsEmpty() {
		t.Fatal("expected nil key pair to be empty")
	}

	kp := &commonrecordsutils.KeyPairValue{
		PublicKey:  types.StringValue("pub"),
		PrivateKey: types.StringNull(),
	}
	if kp.IsEmpty() {
		t.Fatal("expected public key only to be non-empty")
	}
}
