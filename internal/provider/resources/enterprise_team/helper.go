// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package enterpriseteam

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// parseNodesResponse is used by convertNodeToName for ID-to-name resolution; kept for potential reuse.
//
//nolint:unused
func parseNodesResponse(data interface{}) ([]utils.EnterpriseNodeResponse, error) {
	var nodes []utils.EnterpriseNodeResponse

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("unable to process the response from Keeper Commander Service Mode API: %w", err)
	}

	if err := json.Unmarshal(dataBytes, &nodes); err != nil {
		return nil, fmt.Errorf("unable to parse enterprise nodes list from Service Mode API response: %w", err)
	}

	return nodes, nil
}

// convertNodeToName converts node from API response to node name.
// API may return node as name or ID. If it's an ID, we fetch nodes and convert to name.
// Final state always stores node name (not ID).
//
//nolint:unused
func convertNodeToName(ctx context.Context, apiManager *api.ApiManager, nodeFromApi string) (types.String, error) {
	// Handle empty/null cases
	nodeFromApi = strings.TrimSpace(nodeFromApi)
	if nodeFromApi == "" {
		return types.StringNull(), nil
	}

	// Check if nodeFromApi is numeric (likely an ID)
	// If it's not numeric, assume it's already a name and use it directly
	nodeIdInt, atoiErr := strconv.Atoi(nodeFromApi)
	if atoiErr != nil {
		// Not numeric - assume it's a name, use it directly (atoi error is expected, not propagated).
		//nolint:nilerr
		return types.StringValue(nodeFromApi), nil
	}

	// It's numeric - need to convert ID to name
	// Fetch all nodes to build lookup
	nodesResp, err := apiManager.ExecuteCommand(ctx, "enterprise-info -n --format json -v -q", "Unable to fetch enterprise nodes")
	if err != nil {
		return types.StringNull(), fmt.Errorf("failed to fetch nodes: %w", err)
	}

	nodes, err := parseNodesResponse(nodesResp.Data)
	if err != nil {
		return types.StringNull(), fmt.Errorf("failed to parse nodes: %w", err)
	}

	// Find node by ID
	for _, node := range nodes {
		if node.NodeId == nodeIdInt {
			if node.Name != "" {
				nodeName, ok := utils.ExtractNodeIDFromCreateNodeResponse(node.Name)
				if ok {
					return types.StringValue(nodeName), nil
				}
				return types.StringValue(node.Name), nil
			}
		}
	}

	// Node ID not found - return null (node may have been deleted)
	// This handles edge case where node was deleted but team still references it
	return types.StringNull(), nil
}

// Note: After creating a node, service mode api returns message like: "Node is created with Node ID: 1169425105420462"
// This function extracts the node id from the response.
func extractTeamIdFromCreateTeamResponse(s string) (string, bool) {
	_, after, ok := strings.Cut(s, "Team ID:")
	if !ok {
		return "", false
	}
	id := strings.TrimSpace(after)
	if id == "" {
		return "", false
	}
	return id, true
}
