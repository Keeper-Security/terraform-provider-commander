// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package secretsmanager

type CreateAppResponse struct {
	AppName   string `json:"app_name"`
	AppUID    string `json:"app_uid"`
	CreatedAt string `json:"created_at"`
	Message   string `json:"message"`
}

type ShareResponse struct {
	UID       string `json:"uid"`
	Editable  bool   `json:"editable"`
	ShareType string `json:"share_type"`
	Title     string `json:"title"`
	Type      string `json:"type"`
}

type UserResponse struct {
	Username   string `json:"username"`
	Role       string `json:"role"`
	Editable   bool   `json:"editable"`
	ShareAdmin bool   `json:"share_admin"`
	Shareable  bool   `json:"shareable"`
}

type GetAppResponse struct {
	AppName       string          `json:"app_name"`
	AppUID        string          `json:"app_uid"`
	ClientDevices []interface{}   `json:"client_devices"`
	Shares        []ShareResponse `json:"shares"`
	Users         []UserResponse  `json:"users"`
}
