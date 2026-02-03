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

// EnterpriseRoleResponse represents a role from the enterprise-info API response
type EnterpriseRoleResponse struct {
	RoleId       int      `json:"role_id"`
	Name         string   `json:"name"`
	Node         string   `json:"node"`
	VisibleBelow bool     `json:"visible_below"`
	DefaultRole  bool     `json:"default_role"`
	Admin        bool     `json:"admin"`
	UserCount    int      `json:"user_count"`
	Users        []string `json:"users"`
	TeamCount    int      `json:"team_count"`
	Teams        []string `json:"teams"`
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
