// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package utils

// NodeInfo represents a node from the enterprise-info API response.
type EnterpriseNodeResponse struct {
	NodeId         int    `json:"node_id"`
	Name           string `json:"name"`
	ParentNodeName string `json:"parent_node"`
	ParentNodeId   int    `json:"parent_id"`
	Isolated       bool   `json:"isolated"`
}

type ManagedCompanyResponse struct {
	CompanyId   int      `json:"company_id"`
	CompanyName string   `json:"company_name"`
	Node        string   `json:"node"`
	NodeName    string   `json:"node_name"`
	Plan        string   `json:"plan"`
	Storage     string   `json:"storage"`
	Addons      []string `json:"addons"`
	Allocated   int      `json:"allocated"`
}

// EnterpriseTeamResponse represents the team information from the read API response.
type EnterpriseTeamResponse struct {
	TeamUid   string   `json:"team_uid"`
	Name      string   `json:"name"`
	Restricts string   `json:"restricts"`
	Node      string   `json:"node"`
	Users     []string `json:"users"`
	Roles     []string `json:"roles"`
}

// ManagedNodePermission represents one entry in the managed_nodes_permissions array from the API.
type ManagedNodePermission struct {
	NodeName   string   `json:"node_name"`
	NodeId     int64    `json:"node_id"`
	Cascade    bool     `json:"cascade"`
	Privileges []string `json:"privileges"`
}

// EnterpriseRoleResponse represents a role from the enterprise-info API response.
type EnterpriseRoleResponse struct {
	RoleId                  int                     `json:"role_id"`
	Name                    string                  `json:"name"`
	Node                    string                  `json:"node"`
	VisibleBelow            bool                    `json:"visible_below"`
	DefaultRole             bool                    `json:"default_role"`
	Admin                   bool                    `json:"admin"`
	Users                   []string                `json:"users"`
	Teams                   []string                `json:"teams"`
	Enforcements            map[string]string       `json:"enforcements"`
	ManagedNodesPermissions []ManagedNodePermission `json:"managed_nodes_permissions"`
}

// EnterpriseUserResponse represents a user from the API response.
type EnterpriseUserResponse struct {
	UserId   int      `json:"user_id"`
	Email    string   `json:"email"`
	Status   string   `json:"status"`
	Name     string   `json:"name"`
	JobTitle string   `json:"job_title"`
	Roles    []string `json:"roles"`
	Teams    []string `json:"teams"`
	Node     string   `json:"node"`
}

// EnterpriseScimResponse represents the SCIM configuration from the read API response.
type EnterpriseScimResponse struct {
	ScimID            int    `json:"scim_id"`
	ScimURL           string `json:"scim_url"`
	NodeID            int    `json:"node_id"`
	NodeName          string `json:"node_name"`
	Status            string `json:"status"`
	Prefix            string `json:"prefix"`
	UniqueGroups      bool   `json:"unique_groups"`
	ProvisioningToken string `json:"provisioning_token"`
}

// SharedFolderRecordEntry is one element of the records array from get shared folder --format json.
type SharedFolderRecordEntry struct {
	RecordUID string `json:"record_uid"`
	CanShare  bool   `json:"can_share"`
	CanEdit   bool   `json:"can_edit"`
}

// SharedFolderUserEntry is one element of the users array from get shared folder --format json.
type SharedFolderUserEntry struct {
	Username      string `json:"username"`
	UserID        string `json:"user_id"`
	ManageUsers   bool   `json:"manage_users"`
	ManageRecords bool   `json:"manage_records"`
	Expiration    string `json:"expiration"`
}

// SharedFolderResponse is the data payload from get SHARED_FOLDER_ID --format json.
type SharedFolderResponse struct {
	SharedFolderUID string                    `json:"shared_folder_uid"`
	Path            string                    `json:"path"`
	ManageUsers     bool                      `json:"manage_users"`
	ManageRecords   bool                      `json:"manage_records"`
	CanShare        bool                      `json:"can_share"`
	CanEdit         bool                      `json:"can_edit"`
	Records         []SharedFolderRecordEntry `json:"records"`
	Users           []SharedFolderUserEntry   `json:"users"`
}

// PamConfigListSharedFolder is the shared_folder object from pam config list --format json.
type PamConfigListSharedFolder struct {
	Name string `json:"name"`
	UID  string `json:"uid"`
}

// PamConfigListResponse is the data payload from pam config list --config ID --format json.
type PamConfigListResponse struct {
	UID          string                    `json:"uid"`
	Name         string                    `json:"name"`
	ConfigType   string                    `json:"config_type"`
	SharedFolder PamConfigListSharedFolder `json:"shared_folder"`
	GatewayUID   string                    `json:"gateway_uid"`
}
