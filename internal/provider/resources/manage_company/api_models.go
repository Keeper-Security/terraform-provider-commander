package managecompany

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
