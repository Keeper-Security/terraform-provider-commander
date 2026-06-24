// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamuser_test

import (
	"encoding/json"
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

func TestMapRotationSettingsToState_ConfigurationFromVaultRecordWhenDagDebugPresent(t *testing.T) {
	t.Parallel()

	state := commonpamuser.PamUserSharedModel{}
	rotInfo := &commonpamuser.PamRotationInfoResponse{
		PamConfigUID: "rot-info-config-uid",
		Disable:      false,
	}
	rec := &utils.VaultRecordGetResponse{
		PamConfigurationUID: "vault-config-uid",
		DagDebug:            &utils.DagDebugResponse{},
	}

	commonpamuser.MapRotationSettingsToState(rotInfo, rec, nil, &state)

	if !state.RotationSettings.Configuration.Equal(types.StringValue("rot-info-config-uid")) {
		t.Fatalf("configuration = %v, want rot-info-config-uid", state.RotationSettings.Configuration)
	}
}

func TestMapRotationSettingsToState_GeneralProfileFromDagDebug(t *testing.T) {
	t.Parallel()

	state := commonpamuser.PamUserSharedModel{}
	rotInfo := &commonpamuser.PamRotationInfoResponse{Disable: false}
	rec := &utils.VaultRecordGetResponse{
		DagDebug: &utils.DagDebugResponse{
			RotationProfile: &utils.DagDebugRotationProfileResponse{
				Type:             commonpamuser.RotProfileGeneral,
				ResourceUID:      "resource-uid",
				ConfigurationUID: "config-uid",
			},
		},
	}

	commonpamuser.MapRotationSettingsToState(rotInfo, rec, nil, &state)

	if !state.RotationSettings.RotationProfile.Equal(types.StringValue(commonpamuser.RotProfileGeneral)) {
		t.Fatalf("rotation_profile = %v, want general", state.RotationSettings.RotationProfile)
	}
	if !state.RotationSettings.Resource.Equal(types.StringValue("resource-uid")) {
		t.Fatalf("resource = %v, want resource-uid", state.RotationSettings.Resource)
	}
}

func TestMapRotationSettingsToState_SaaSFromParentAclEdge(t *testing.T) {
	t.Parallel()

	state := commonpamuser.PamUserSharedModel{}
	rotInfo := &commonpamuser.PamRotationInfoResponse{
		PamConfigUID: "pam-config-from-rot-info",
		Disable:      false,
	}
	rec := &utils.VaultRecordGetResponse{
		DagDebug: &utils.DagDebugResponse{
			RotationProfile: &utils.DagDebugRotationProfileResponse{
				ConfigurationUID: "pam-config-uid",
			},
			ParentAclEdge: &utils.DagDebugParentAclEdgeResponse{
				Content: &utils.DagDebugParentAclEdgeContentResponse{
					RotationSettings: &utils.DagDebugParentAclEdgeContentRotationSettingsResponse{
						Noop:              true,
						SaaSRecordUIDList: []string{"saas-config-uid"},
					},
				},
			},
		},
	}

	commonpamuser.MapRotationSettingsToState(rotInfo, rec, nil, &state)

	if state.RotationSettings == nil {
		t.Fatal("expected rotation settings to be set")
	}
	if !state.RotationSettings.RotationProfile.Equal(types.StringValue(commonpamuser.RotProfileSaaS)) {
		t.Fatalf("rotation_profile = %v, want saas", state.RotationSettings.RotationProfile)
	}
	if !state.RotationSettings.Configuration.Equal(types.StringValue("pam-config-from-rot-info")) {
		t.Fatalf("configuration = %v, want pam-config-from-rot-info", state.RotationSettings.Configuration)
	}
	if !state.RotationSettings.SaaSConfig.Equal(types.StringValue("saas-config-uid")) {
		t.Fatalf("saas_config = %v, want saas-config-uid", state.RotationSettings.SaaSConfig)
	}
}

func TestMapRotationSettingsToState_SaaSFromCamelCaseParentAclEdge(t *testing.T) {
	t.Parallel()

	state := commonpamuser.PamUserSharedModel{}
	rotInfo := &commonpamuser.PamRotationInfoResponse{Disable: false}
	rec := &utils.VaultRecordGetResponse{
		DagDebug: &utils.DagDebugResponse{
			ParentAclEdge: &utils.DagDebugParentAclEdgeResponse{
				Content: &utils.DagDebugParentAclEdgeContentResponse{
					RotationSettings: &utils.DagDebugParentAclEdgeContentRotationSettingsResponse{
						Noop:              true,
						SaaSRecordUIDList: []string{"saas-config-uid"},
					},
				},
			},
		},
	}

	commonpamuser.MapRotationSettingsToState(rotInfo, rec, nil, &state)

	if !state.RotationSettings.RotationProfile.Equal(types.StringValue(commonpamuser.RotProfileSaaS)) {
		t.Fatalf("rotation_profile = %v, want saas", state.RotationSettings.RotationProfile)
	}
	if !state.RotationSettings.SaaSConfig.Equal(types.StringValue("saas-config-uid")) {
		t.Fatalf("saas_config = %v, want saas-config-uid", state.RotationSettings.SaaSConfig)
	}
}

func TestMapRotationSettingsToState_SaaSWithoutRotationProfileDoesNotPanic(t *testing.T) {
	t.Parallel()

	state := commonpamuser.PamUserSharedModel{}
	rotInfo := &commonpamuser.PamRotationInfoResponse{
		PamConfigUID: "z_GCs8J-JNwgg-0k1UkIdg",
		Disable:      false,
	}
	rec := &utils.VaultRecordGetResponse{
		DagDebug: &utils.DagDebugResponse{
			ParentAclEdge: &utils.DagDebugParentAclEdgeResponse{
				Content: &utils.DagDebugParentAclEdgeContentResponse{
					RotationSettings: &utils.DagDebugParentAclEdgeContentRotationSettingsResponse{
						Noop:              true,
						SaaSRecordUIDList: []string{"rt9LG5vZJCO1a2-Sg2hk3A"},
					},
				},
			},
		},
	}

	commonpamuser.MapRotationSettingsToState(rotInfo, rec, nil, &state)

	if !state.RotationSettings.RotationProfile.Equal(types.StringValue(commonpamuser.RotProfileSaaS)) {
		t.Fatalf("rotation_profile = %v, want saas", state.RotationSettings.RotationProfile)
	}
	if !state.RotationSettings.Configuration.Equal(types.StringValue("z_GCs8J-JNwgg-0k1UkIdg")) {
		t.Fatalf("configuration = %v, want z_GCs8J-JNwgg-0k1UkIdg", state.RotationSettings.Configuration)
	}
	if !state.RotationSettings.SaaSConfig.Equal(types.StringValue("rt9LG5vZJCO1a2-Sg2hk3A")) {
		t.Fatalf("saas_config = %v, want rt9LG5vZJCO1a2-Sg2hk3A", state.RotationSettings.SaaSConfig)
	}
}

func TestMapRotationSettingsToState_RealSaaSDagDebugAPI(t *testing.T) {
	t.Parallel()

	const dagDebugJSON = `{
		"all_edges": [
			{"type": "acl", "head_uid": "z_GCs8J-JNwgg-0k1UkIdg"},
			{"type": "acl", "head_uid": "DZcSrlHAeUPfCWkljaKJ-g"},
			{"type": "data", "head_uid": "2rAkAzKWrJ6BTazfHTHIew"}
		],
		"parentAclEdge": {
			"parent_uid": "DZcSrlHAeUPfCWkljaKJ-g",
			"parent_type": "pam_network",
			"content": {
				"rotation_settings": {
					"noop": true,
					"saas_record_uid_list": ["rt9LG5vZJCO1a2-Sg2hk3A"]
				}
			}
		}
	}`

	var dagDebug utils.DagDebugResponse
	if err := json.Unmarshal([]byte(dagDebugJSON), &dagDebug); err != nil {
		t.Fatalf("unmarshal dagDebug: %v", err)
	}

	state := commonpamuser.PamUserSharedModel{}
	rotInfo := &commonpamuser.PamRotationInfoResponse{
		Disable:      false,
		ScheduleData: `[{"type":"CRON","cron":"0 0 4 * * ?","tz":"Pacific/Honolulu"}]`,
	}
	rec := &utils.VaultRecordGetResponse{
		PamConfigurationUID: "z_GCs8J-JNwgg-0k1UkIdg",
		DagDebug:            &dagDebug,
	}

	commonpamuser.MapRotationSettingsToState(rotInfo, rec, nil, &state)

	if state.RotationSettings == nil {
		t.Fatal("expected rotation settings to be set")
	}
	if !state.RotationSettings.RotationProfile.Equal(types.StringValue(commonpamuser.RotProfileSaaS)) {
		t.Fatalf("rotation_profile = %v, want saas", state.RotationSettings.RotationProfile)
	}
	if !state.RotationSettings.Configuration.Equal(types.StringValue("z_GCs8J-JNwgg-0k1UkIdg")) {
		t.Fatalf("configuration = %v, want z_GCs8J-JNwgg-0k1UkIdg", state.RotationSettings.Configuration)
	}
	if !state.RotationSettings.SaaSConfig.Equal(types.StringValue("rt9LG5vZJCO1a2-Sg2hk3A")) {
		t.Fatalf("saas_config = %v, want rt9LG5vZJCO1a2-Sg2hk3A", state.RotationSettings.SaaSConfig)
	}
}

func TestDagDebugResponseFromVaultRecord(t *testing.T) {
	t.Parallel()

	if got := commonpamuser.DagDebugResponseFromVaultRecord(nil); got != nil {
		t.Fatalf("nil record: got %v, want nil", got)
	}

	rec := &utils.VaultRecordGetResponse{}
	if got := commonpamuser.DagDebugResponseFromVaultRecord(rec); got != nil {
		t.Fatalf("nil dagDebug: got %v, want nil", got)
	}

	dagDebug := &utils.DagDebugResponse{
		RotationProfile: &utils.DagDebugRotationProfileResponse{Type: commonpamuser.RotProfileGeneral},
	}
	rec.DagDebug = dagDebug
	if got := commonpamuser.DagDebugResponseFromVaultRecord(rec); got != dagDebug {
		t.Fatalf("got %v, want %v", got, dagDebug)
	}
}
