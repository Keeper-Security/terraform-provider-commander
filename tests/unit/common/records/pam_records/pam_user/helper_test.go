// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamuser_test

import (
	"testing"

	commonpamuser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records/pam_user"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestMapRotationSettingsToState_NilRotationProfile(t *testing.T) {
	t.Parallel()

	state := commonpamuser.PamUserSharedModel{}
	rotInfo := &commonpamuser.PamRotationInfoResponse{
		PamConfigUID: "config-uid-123",
		ScheduleType: commonpamuser.RotProfileScheduleTypeManual,
		Disable:      false,
	}

	commonpamuser.MapRotationSettingsToState(rotInfo, nil, nil, &state)

	if state.RotationSettings == nil {
		t.Fatal("expected rotation settings to be set")
	}
	if !state.RotationSettings.Configuration.Equal(types.StringValue("config-uid-123")) {
		t.Fatalf("configuration = %v, want config-uid-123", state.RotationSettings.Configuration)
	}
	if !state.RotationSettings.OnDemand.Equal(types.BoolValue(true)) {
		t.Fatalf("on_demand = %v, want true", state.RotationSettings.OnDemand)
	}
	if !state.RotationSettings.Enabled.Equal(types.BoolValue(true)) {
		t.Fatalf("enabled = %v, want true", state.RotationSettings.Enabled)
	}
}

func TestRotationProfileFromVaultRecord(t *testing.T) {
	t.Parallel()

	if got := commonpamuser.RotationProfileFromVaultRecord(nil); got != nil {
		t.Fatalf("nil record: got %v, want nil", got)
	}

	rec := &utils.VaultRecordGetResponse{}
	if got := commonpamuser.RotationProfileFromVaultRecord(rec); got != nil {
		t.Fatalf("nil dagDebug: got %v, want nil", got)
	}

	profile := &utils.DagDebugRotationProfileResponse{Type: commonpamuser.RotProfileGeneral}
	rec.DagDebug = &utils.DagDebugResponse{RotationProfile: profile}
	if got := commonpamuser.RotationProfileFromVaultRecord(rec); got != profile {
		t.Fatalf("got %v, want %v", got, profile)
	}
}
