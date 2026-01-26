// THIS FILES STORE API RESPONSE STRUCTS
package enterpriserole

// RoleResponse represents a role from the enterprise-info API response
// Used in the Read operation to unmarshal role data returned by the Commander CLI
type RoleResponse struct {
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

// NodeInfo represents a node from the enterprise-info API response
type NodeInfo struct {
	NodeId int    `json:"node_id"`
	Name   string `json:"name"`
}
