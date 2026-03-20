// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package utils

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ParseNodesResponse parses the JSON response from enterprise-info -n command.
func ParseNodesResponse(data interface{}) ([]EnterpriseNodeResponse, error) {
	var nodes []EnterpriseNodeResponse

	if err := UnmarshalApiResponse(data, &nodes); err != nil {
		return nil, fmt.Errorf("unable to parse enterprise nodes list from API response: %w", err)
	}

	return nodes, nil
}

// BuildNodeLookupMaps creates lookup maps from API response.
// Uses ExtractNodeName for names so path format (e.g. "Parent\\Node Name") matches the short name.
func BuildNodeLookupMaps(nodesRespData []EnterpriseNodeResponse) LookupMaps {
	identifierToId := make(map[string]string)
	idToIdentifier := make(map[string]string)

	for _, node := range nodesRespData {
		if node.NodeId > 0 && node.Name != "" {
			nodeIdStr := strconv.Itoa(node.NodeId)
			nodeName := ExtractNodeName(node.Name)
			identifierToId[nodeName] = nodeIdStr
			identifierToId[node.Name] = nodeIdStr // support both path and short name
			idToIdentifier[nodeIdStr] = nodeName
		}
	}

	return LookupMaps{
		IdentifierToId: identifierToId,
		IdToIdentifier: idToIdentifier,
	}
}

// RestoreUserInputFormatForNode converts the node value from API response back to the format
// that the user originally provided in their Terraform configuration.
//
// This function preserves the original user input format to prevent false diffs in Terraform plans.
// If a user specified node by ID (e.g., "123"), the function will return the ID. If they specified
// by name (e.g., "Company Inc"), it will return the name.
//
// Parameters:
//   - apiNodeValue: Node name (or path) returned by the API (from enterprise-info command)
//   - currentStateNode: Current Terraform state containing node (what user originally provided)
//
// Returns:
//   - types.String: Node in the original user input format (name or node ID)
//   - error: Error if fetching nodes or building lookup maps fails
//
// Example:
//
//	User config: node = "123"
//	API returns: "Company Inc"
//	Function returns: "123" (preserves original ID)
func RestoreUserInputFormatForNode(ctx context.Context, apiManager *api.ApiManager, apiNodeValue string, currentStateNode types.String) (types.String, error) {
	apiNodeValue = strings.TrimSpace(apiNodeValue)
	if apiNodeValue == "" {
		return types.StringNull(), nil
	}

	// Fetch all nodes to build lookup maps
	nodesResp, err := apiManager.ExecuteCommand(ctx, "enterprise-info -n --format json -v -q", "Unable to fetch enterprise nodes")
	if err != nil {
		return types.StringNull(), fmt.Errorf("failed to fetch nodes: %w", err)
	}

	nodesData, err := ParseNodesResponse(nodesResp.Data)
	if err != nil {
		return types.StringNull(), fmt.Errorf("failed to parse nodes: %w", err)
	}

	lookup := BuildNodeLookupMaps(nodesData)

	// Normalize API value: API may return name or path (e.g. "Parent\\Node Name")
	apiNodeName := ExtractNodeName(apiNodeValue)
	var apiNodeId string
	if id, isId := lookup.IdToIdentifier[apiNodeValue]; isId {
		apiNodeId = apiNodeValue
		apiNodeName = id
	} else if id, isName := lookup.IdentifierToId[apiNodeName]; isName {
		apiNodeId = id
	} else if id, isName := lookup.IdentifierToId[apiNodeValue]; isName {
		apiNodeId = id
	}
	// If apiNodeId still empty, node may have been deleted; we'll return API value or state below

	// No current state (e.g. first read after import) - return API format (name)
	if currentStateNode.IsNull() || currentStateNode.IsUnknown() || currentStateNode.ValueString() == "" {
		if apiNodeId != "" {
			return types.StringValue(apiNodeName), nil
		}
		return types.StringValue(apiNodeValue), nil
	}

	stateVal := strings.TrimSpace(currentStateNode.ValueString())

	// State has node ID (numeric)
	if _, err := strconv.Atoi(stateVal); err == nil {
		if stateVal == apiNodeId {
			return currentStateNode, nil
		}
		// State ID doesn't match API node (external change) - return API name
		if apiNodeId != "" {
			return types.StringValue(apiNodeName), nil
		}
		return types.StringValue(apiNodeValue), nil
	}

	// State has node name - check if it matches API node
	stateName := ExtractNodeName(stateVal)
	if stateName == apiNodeName || (apiNodeId != "" && lookup.IdentifierToId[stateName] == apiNodeId) {
		return currentStateNode, nil
	}

	// State doesn't match (external change) - return API name
	if apiNodeId != "" {
		return types.StringValue(apiNodeName), nil
	}
	return types.StringValue(apiNodeValue), nil
}
