// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamuser

type PamRotationInfoResponse struct {
	PamConfigUID              string                                                `json:"pam_config_uid"`
	AdminResourceUID          string                                                `json:"admin_resource_uid"`
	ScheduleType              string                                                `json:"schedule_type"`
	ScheduleData              string                                                `json:"schedule_data"` // We will use this to set ScheduleCron, ScheduleJSON
	PasswordComplexityDetails *PamRotationInfoPasswordComplexityDetailsDataResponse `json:"password_complexity_detail"`
	Disable                   bool                                                  `json:"disabled"`
}

type PamRotationInfoPasswordComplexityDetailsDataResponse struct {
	Length    int `json:"length"`
	Capital   int `json:"caps"`
	Lowercase int `json:"lowercase"`
	Digits    int `json:"digits"`
	Special   int `json:"special"`
}
