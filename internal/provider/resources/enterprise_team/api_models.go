package enterpriseteam

// TeamInfo represents a team from the API response
type TeamInfo struct {
	TeamUid string `json:"team_uid"`
	Name    string `json:"name"`
}

// EnterpriseTeamReadResponse represents the team information from the read API response
type EnterpriseTeamReadResponse struct {
	TeamUid   string   `json:"team_uid"`
	Name      string   `json:"name"`
	Restricts string   `json:"restricts"`
	Node      string   `json:"node"`
	Users     []string `json:"users"`
	Roles     []string `json:"roles"`
}
