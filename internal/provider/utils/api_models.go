// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package utils

// NodeInfo represents a node from the enterprise-info API response
type EnterpriseNodeResponse struct {
	NodeId         int    `json:"node_id"`
	Name           string `json:"name"`
	ParentNodeName string `json:"parent_node"`
	ParentNodeId   int    `json:"parent_id"`
}

type ManageCompanyResponse struct {
	CompanyId   int      `json:"company_id"`
	CompanyName string   `json:"company_name"`
	Node        string   `json:"node"`
	NodeName    string   `json:"node_name"`
	Plan        string   `json:"plan"`
	Storage     string   `json:"storage"`
	Addons      []string `json:"addons"`
	Allocated   int      `json:"allocated"`
}

// EnterpriseTeamResponse represents the team information from the read API response
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

// EnterpriseRoleResponse represents a role from the enterprise-info API response
type EnterpriseRoleResponse struct {
	RoleId                  int                     `json:"role_id"`
	Name                    string                  `json:"name"`
	Node                    string                  `json:"node"`
	VisibleBelow            bool                    `json:"visible_below"`
	DefaultRole             bool                    `json:"default_role"`
	Admin                   bool                    `json:"admin"`
	Users                   []string                `json:"users"`
	Teams                   []string                `json:"teams"`
	Enforcements            []string                `json:"enforcements"`
	ManagedNodesPermissions []ManagedNodePermission `json:"managed_nodes_permissions"`
}

// EnterpriseUserResponse represents a user from the API response
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
