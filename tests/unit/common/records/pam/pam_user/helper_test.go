// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamuser_test

import (
	"encoding/json"
	"testing"

	commonpamuser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam/pam_user"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestMapRotationSettingsToState_NilRotInfoPreservesExisting(t *testing.T) {
	t.Parallel()

	state := commonpamuser.PamUserSharedModel{
		RotationSettings: &commonpamuser.PamUserRotationSettings{
			RotationProfile: types.StringValue(commonpamuser.RotProfileGeneral),
			Configuration:   types.StringValue("existing-config"),
		},
	}

	commonpamuser.MapRotationSettingsToState(nil, nil, state.RotationSettings, &state)

	if state.RotationSettings == nil {
		t.Fatal("expected rotation settings to be preserved")
	}
	if !state.RotationSettings.RotationProfile.Equal(types.StringValue(commonpamuser.RotProfileGeneral)) {
		t.Fatalf("rotation_profile = %v, want general", state.RotationSettings.RotationProfile)
	}
	if !state.RotationSettings.Configuration.Equal(types.StringValue("existing-config")) {
		t.Fatalf("configuration = %v, want existing-config", state.RotationSettings.Configuration)
	}
}

func TestMapRotationSettingsToState_RotInfoWithEmptyVaultRecord(t *testing.T) {
	t.Parallel()

	state := commonpamuser.PamUserSharedModel{}
	rotInfo := &commonpamuser.PamRotationInfoResponse{
		PamConfigUID: "config-uid-123",
		ScheduleType: commonpamuser.RotProfileScheduleTypeManual,
		Disable:      false,
	}
	rec := &utils.VaultRecordGetResponse{
		DagDebug: &utils.DagDebugResponse{},
	}

	commonpamuser.MapRotationSettingsToState(rotInfo, rec, nil, &state)

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
	if !state.RotationSettings.RotationProfile.IsNull() {
		t.Fatalf("rotation_profile = %v, want null without dagDebug", state.RotationSettings.RotationProfile)
	}
}

func TestMapRotationSettingsToState_ConfigurationPrefersRotInfoOverVaultRecord(t *testing.T) {
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

func TestMapRotationSettingsToState_ConfigurationFromVaultRecordWhenRotInfoEmpty(t *testing.T) {
	t.Parallel()

	state := commonpamuser.PamUserSharedModel{}
	rotInfo := &commonpamuser.PamRotationInfoResponse{Disable: false}
	rec := &utils.VaultRecordGetResponse{
		PamConfigurationUID: "vault-config-uid",
		DagDebug:            &utils.DagDebugResponse{},
	}

	commonpamuser.MapRotationSettingsToState(rotInfo, rec, nil, &state)

	if !state.RotationSettings.Configuration.Equal(types.StringValue("vault-config-uid")) {
		t.Fatalf("configuration = %v, want vault-config-uid", state.RotationSettings.Configuration)
	}
}

func TestMapRotationSettingsToState_GeneralProfileFromVaultRecord(t *testing.T) {
	t.Parallel()

	state := commonpamuser.PamUserSharedModel{}
	rotInfo := &commonpamuser.PamRotationInfoResponse{Disable: false}
	rec := &utils.VaultRecordGetResponse{
		DagDebug: &utils.DagDebugResponse{},
		RotationProfile: &utils.RotationProfileResponse{
			Type:             commonpamuser.RotProfileGeneral,
			ResourceUID:      "resource-uid",
			ConfigurationUID: "config-uid",
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

func TestMapRotationSettingsToState_IAMUserProfileFromVaultRecord(t *testing.T) {
	t.Parallel()

	state := commonpamuser.PamUserSharedModel{}
	rotInfo := &commonpamuser.PamRotationInfoResponse{
		PamConfigUID: "iam-config-uid",
		Disable:      false,
	}
	rec := &utils.VaultRecordGetResponse{
		DagDebug: &utils.DagDebugResponse{},
		RotationProfile: &utils.RotationProfileResponse{
			Type: commonpamuser.RotProfileIAMUser,
		},
	}

	commonpamuser.MapRotationSettingsToState(rotInfo, rec, nil, &state)

	if !state.RotationSettings.RotationProfile.Equal(types.StringValue(commonpamuser.RotProfileIAMUser)) {
		t.Fatalf("rotation_profile = %v, want iam_user", state.RotationSettings.RotationProfile)
	}
	if !state.RotationSettings.Resource.IsNull() {
		t.Fatalf("resource = %v, want null for iam_user", state.RotationSettings.Resource)
	}
}

func TestMapRotationSettingsToState_DisabledRotationParsesScheduleAfterClear(t *testing.T) {
	t.Parallel()

	state := commonpamuser.PamUserSharedModel{}
	rotInfo := &commonpamuser.PamRotationInfoResponse{
		PamConfigUID: "config-uid",
		Disable:      true,
		ScheduleType: commonpamuser.RotProfileScheduleTypeManual,
		ScheduleData: `{"type":"CRON","cron":"0 0 4 * * ?","tz":"Etc/UTC"}`,
	}
	rec := &utils.VaultRecordGetResponse{
		DagDebug: &utils.DagDebugResponse{},
		RotationProfile: &utils.RotationProfileResponse{
			Type: commonpamuser.RotProfileGeneral,
		},
	}

	commonpamuser.MapRotationSettingsToState(rotInfo, rec, nil, &state)

	if !state.RotationSettings.Enabled.Equal(types.BoolValue(false)) {
		t.Fatalf("enabled = %v, want false", state.RotationSettings.Enabled)
	}
	if !state.RotationSettings.OnDemand.IsNull() {
		t.Fatalf("on_demand = %v, want null when disabled", state.RotationSettings.OnDemand)
	}
	// parseScheduleValue runs after the disabled branch and repopulates schedule fields.
	if !state.RotationSettings.ScheduleCron.Equal(types.StringValue("0 0 4 * * ?")) {
		t.Fatalf("schedule_cron = %v, want 0 0 4 * * ?", state.RotationSettings.ScheduleCron)
	}
}

func TestMapRotationSettingsToState_CronSchedule(t *testing.T) {
	t.Parallel()

	state := commonpamuser.PamUserSharedModel{}
	rotInfo := &commonpamuser.PamRotationInfoResponse{
		Disable:      false,
		ScheduleType: "cron",
		ScheduleData: `{"type":"CRON","cron":"0 0 4 * * ?","tz":"Etc/UTC"}`,
	}
	rec := &utils.VaultRecordGetResponse{
		DagDebug: &utils.DagDebugResponse{},
		RotationProfile: &utils.RotationProfileResponse{
			Type: commonpamuser.RotProfileGeneral,
		},
	}

	commonpamuser.MapRotationSettingsToState(rotInfo, rec, nil, &state)

	if !state.RotationSettings.ScheduleCron.Equal(types.StringValue("0 0 4 * * ?")) {
		t.Fatalf("schedule_cron = %v, want 0 0 4 * * ?", state.RotationSettings.ScheduleCron)
	}
	if !state.RotationSettings.OnDemand.IsNull() {
		t.Fatalf("on_demand = %v, want null for cron schedule", state.RotationSettings.OnDemand)
	}
	if !state.RotationSettings.ScheduleJSON.IsNull() {
		t.Fatalf("schedule_json = %v, want null for cron schedule", state.RotationSettings.ScheduleJSON)
	}
}

func TestMapRotationSettingsToState_WeeklyScheduleJSON(t *testing.T) {
	t.Parallel()

	const scheduleJSON = `{"type":"WEEKLY","utcTime":"00:00","weekday":"SATURDAY","intervalCount":1}`

	state := commonpamuser.PamUserSharedModel{}
	rotInfo := &commonpamuser.PamRotationInfoResponse{
		Disable:      false,
		ScheduleType: "weekly",
		ScheduleData: scheduleJSON,
	}
	rec := &utils.VaultRecordGetResponse{
		DagDebug: &utils.DagDebugResponse{},
		RotationProfile: &utils.RotationProfileResponse{
			Type: commonpamuser.RotProfileGeneral,
		},
	}

	commonpamuser.MapRotationSettingsToState(rotInfo, rec, nil, &state)

	if !state.RotationSettings.ScheduleJSON.Equal(types.StringValue(scheduleJSON)) {
		t.Fatalf("schedule_json = %v, want %s", state.RotationSettings.ScheduleJSON, scheduleJSON)
	}
	if !state.RotationSettings.ScheduleCron.IsNull() {
		t.Fatalf("schedule_cron = %v, want null for weekly schedule", state.RotationSettings.ScheduleCron)
	}
}

func TestMapRotationSettingsToState_Complexity(t *testing.T) {
	t.Parallel()

	state := commonpamuser.PamUserSharedModel{}
	rotInfo := &commonpamuser.PamRotationInfoResponse{
		Disable: false,
		PasswordComplexityDetails: &commonpamuser.PamRotationInfoPasswordComplexityDetailsDataResponse{
			Length:    32,
			Capital:   5,
			Lowercase: 1,
			Digits:    1,
			Special:   2,
		},
	}
	rec := &utils.VaultRecordGetResponse{
		DagDebug: &utils.DagDebugResponse{},
	}

	commonpamuser.MapRotationSettingsToState(rotInfo, rec, nil, &state)

	if !state.RotationSettings.Complexity.Equal(types.StringValue("32,5,1,1,2")) {
		t.Fatalf("complexity = %v, want 32,5,1,1,2", state.RotationSettings.Complexity)
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
			ParentAclEdge: &utils.DagDebugParentAclEdgeResponse{
				Content: &utils.DagDebugParentAclEdgeContentResponse{
					RotationSettings: &utils.DagDebugParentAclEdgeContentRotationSettingsResponse{
						Noop:              true,
						SaaSRecordUIDList: []string{"saas-config-uid"},
					},
				},
			},
		},
		RotationProfile: &utils.RotationProfileResponse{
			ConfigurationUID: "pam-config-uid",
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
	if !state.RotationSettings.Resource.IsNull() {
		t.Fatalf("resource = %v, want null for saas", state.RotationSettings.Resource)
	}
}

func TestMapRotationSettingsToState_SaaSWithoutVaultRotationProfile(t *testing.T) {
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
		PamConfigUID: "DZcSrlHAeUPfCWkljaKJ-g",
		Disable:      true,
		PasswordComplexityDetails: &commonpamuser.PamRotationInfoPasswordComplexityDetailsDataResponse{
			Length:    32,
			Capital:   5,
			Lowercase: 1,
			Digits:    1,
			Special:   2,
		},
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
	if !state.RotationSettings.Configuration.Equal(types.StringValue("DZcSrlHAeUPfCWkljaKJ-g")) {
		t.Fatalf("configuration = %v, want DZcSrlHAeUPfCWkljaKJ-g", state.RotationSettings.Configuration)
	}
	if !state.RotationSettings.SaaSConfig.Equal(types.StringValue("rt9LG5vZJCO1a2-Sg2hk3A")) {
		t.Fatalf("saas_config = %v, want rt9LG5vZJCO1a2-Sg2hk3A", state.RotationSettings.SaaSConfig)
	}
	if !state.RotationSettings.Enabled.Equal(types.BoolValue(false)) {
		t.Fatalf("enabled = %v, want false", state.RotationSettings.Enabled)
	}
	if !state.RotationSettings.Complexity.Equal(types.StringValue("32,5,1,1,2")) {
		t.Fatalf("complexity = %v, want 32,5,1,1,2", state.RotationSettings.Complexity)
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

	dagDebug := &utils.DagDebugResponse{}
	rec.DagDebug = dagDebug
	if got := commonpamuser.DagDebugResponseFromVaultRecord(rec); got != dagDebug {
		t.Fatalf("got %v, want %v", got, dagDebug)
	}
}
