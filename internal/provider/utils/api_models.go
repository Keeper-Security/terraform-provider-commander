package utils

// NodeInfo represents a node from the enterprise-info API response
type NodeInfo struct {
	NodeId         int    `json:"node_id"`
	Name           string `json:"name"`
	ParentNodeName string `json:"parent_node"`
	ParentNodeId   int    `json:"parent_id"`
}
