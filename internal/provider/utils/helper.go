// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ErrResourceRemoved is returned by Read callbacks when the resource was not found
// and has been removed from state. The caller should return without calling State.Set.
var ErrResourceRemoved = errors.New("resource removed from state")

// SwitchToManagedCompany switches to the specified managed company.
func SwitchToManagedCompany(ctx context.Context, apiManager *api.ApiManager, manageCompany string) error {
	command := fmt.Sprintf("switch-to-mc '%s'", manageCompany)
	_, err := apiManager.ExecuteCommand(ctx, command, "Unable to switch to manage company")
	return err
}

// SwitchToMsp switches back to MSP context.
func SwitchToMsp(ctx context.Context, apiManager *api.ApiManager) error {
	command := "switch-to-msp"
	_, err := apiManager.ExecuteCommand(ctx, command, "Unable to switch to msp")

	// NOTE: For now we are commenting bec commander cli sending 500 status code
	if err != nil && strings.Contains(err.Error(), "Already MSP") {
		return nil
	}
	return err
}

// Perform msp down.
func MspDown(ctx context.Context, apiManager *api.ApiManager) error {
	command := "msp-down"
	_, err := apiManager.ExecuteCommand(ctx, command, "Unable to perform msp down")
	return err
}

// Perform enterprise down.
func EnterpriseDown(ctx context.Context, apiManager *api.ApiManager) error {
	command := "enterprise-down"
	_, err := apiManager.ExecuteCommand(ctx, command, "Unable to perform enterprise down")
	return err
}

// Perform sync down.
func SyncDown(ctx context.Context, apiManager *api.ApiManager) error {
	command := "sync-down -f"
	_, err := apiManager.ExecuteCommand(ctx, command, "Unable to perform sync down")
	return err
}

// Perform epm sync down.
func EpmSyncDown(ctx context.Context, apiManager *api.ApiManager) error {
	command := "epm sync-down"
	_, err := apiManager.ExecuteCommand(ctx, command, "Unable to perform epm sync down")
	return err
}

// ExecuteWithManagedCompanyContext executes a function with managed company context switching
// If managedCompany is provided and not null, it switches to MC before execution and back to MSP after
// If managedCompany is not provided, it ensures we're in the correct base context (MSP for MSP accounts, Enterprise for Enterprise)
// This prevents race conditions when Terraform runs resources in parallel or different orders.
func ExecuteWithManagedCompanyContext(
	ctx context.Context,
	apiManager *api.ApiManager,
	managedCompany types.String,
	operation func() error,
) (err error) {
	apiManager.LockContext()
	defer apiManager.UnlockContext()

	switchedToMC := false

	hasManagedCompany := !managedCompany.IsNull() && !managedCompany.IsUnknown() && managedCompany.ValueString() != ""

	if hasManagedCompany {
		// Has managed_company (non-empty) - switch to it
		if err := EnterpriseDown(ctx, apiManager); err != nil {
			return fmt.Errorf("failed to sync enterprise data: %w", err)
		}
		if err := SwitchToManagedCompany(ctx, apiManager, managedCompany.ValueString()); err != nil {
			return fmt.Errorf("failed to switch to managed company: %w", err)
		}
		apiManager.SetCurrentContext(managedCompany.ValueString())
		switchedToMC = true
	} else {
		// No managed_company - ensure we're in MSP context only when needed
		if apiManager.IsMspAccount {
			if apiManager.GetCurrentContext() != "" {
				// We think we're in MC (e.g. previous op failed to switch back); switch to MSP
				if err := SwitchToMsp(ctx, apiManager); err != nil {
					return fmt.Errorf("failed to switch to MSP context: %w", err)
				}
				apiManager.SetCurrentContext("")
			}
			// Else already in MSP - skip redundant switch-to-msp API call
		}
		// Always run EnterpriseDown to sync with backend (e.g. external changes); we only skip switch-to-msp above.
		if err := EnterpriseDown(ctx, apiManager); err != nil {
			return fmt.Errorf("failed to sync enterprise data: %w", err)
		}
	}

	defer func() {
		if switchedToMC && apiManager.IsMspAccount {
			if switchErr := SwitchToMsp(ctx, apiManager); switchErr != nil {
				if err != nil {
					err = fmt.Errorf("operation failed: %w; also failed to switch back to MSP: %w", err, switchErr)
				} else {
					err = fmt.Errorf("failed to switch back to MSP: %w", switchErr)
				}
			} else {
				apiManager.SetCurrentContext("")
			}
		}
	}()

	err = operation()
	return err
}

// Note: After creating a node, service mode api returns message like: "Node is created with Node ID: 1169425105420462".
// This function extracts the node id from the response.
func ExtractNodeIDFromCreateNodeResponse(s string) (string, bool) {
	re := regexp.MustCompile(`Node ID:\s*(\d+)`)
	match := re.FindStringSubmatch(s)

	if len(match) < 2 {
		return "", false
	}
	return match[1], true
}

// Note: After creating a user, service mode api returns message like: "user@example.com user invited with Enterprise User ID : 116942510420593".
// This function extracts the user id from the response.
func ExtractUserIDFromCreateUserResponse(s string) (string, bool) {
	re := regexp.MustCompile(`User ID :\s*(\d+)`)
	match := re.FindStringSubmatch(s)

	if len(match) < 2 {
		return "", false
	}
	return match[1], true
}

// Function to extract the node name from the input string like "Metronlabs\\Aditya Dev Inc" -> "Aditya Dev Inc"
// msp-info returns node_name as "Metronlabs\\Aditya Dev Inc" if present in child node or node_name as "Metronlabs" if present in root node.
func ExtractNodeName(input string) string {
	if idx := strings.LastIndex(input, `\`); idx != -1 {
		return input[idx+1:]
	}
	return input
}

// UnmarshalApiResponse unmarshals API response data into a target struct.
// It handles the common pattern of marshaling interface{} to JSON bytes and then unmarshaling into the target type.
// Parameters:
//   - data: The API response data (typically apiResp.Data from ExecuteCommand)
//   - target: A pointer to the struct/slice that should receive the unmarshaled data
//
// Returns an error if marshaling or unmarshaling fails.
//
// Example usage:
//
//	var roles []RoleInfo
//	if err := utils.UnmarshalApiResponse(apiResp.Data, &roles); err != nil {
//	    return fmt.Errorf("failed to parse roles: %w", err)
//	}
func UnmarshalApiResponse(data interface{}, target interface{}) error {
	// Convert apiResp.Data to JSON bytes
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("unable to process the response from Keeper Commander API: %w", err)
	}

	// Unmarshal JSON bytes into target struct
	if err := json.Unmarshal(dataBytes, target); err != nil {
		return fmt.Errorf("unable to parse API response: %w", err)
	}

	return nil
}

// LookupMaps holds the mappings between identifiers (name/email) and IDs.
type LookupMaps struct {
	IdentifierToId map[string]string // identifier (name/email) -> id
	IdToIdentifier map[string]string // id -> identifier (name/email)
}

// ConvertItemsToIdMap is a generic function that converts a types.Set to a map of id -> original input.
// It works for roles, users, and teams by accepting lookup maps and validation functions.
func ConvertItemsToIdMap(
	items types.Set,
	lookup LookupMaps,
	itemType string, // "role", "user", or "team"
	validateItem func(string) (bool, string), // returns (isValid, errorMessage)
) (map[string]string, error) {
	result := make(map[string]string)

	if items.IsNull() || items.IsUnknown() {
		return result, nil
	}

	elements := items.Elements()
	if len(elements) == 0 {
		return result, nil
	}

	seenIds := make(map[string]string) // id -> original input

	for _, itemElem := range elements {
		itemStr, ok := itemElem.(types.String)
		if !ok {
			continue
		}
		userInput := itemStr.ValueString()

		if userInput == "" {
			continue
		}

		var itemId string
		var itemIdentifier string

		// Check if input is an id
		if existingIdentifier, isId := lookup.IdToIdentifier[userInput]; isId {
			itemId = userInput
			itemIdentifier = existingIdentifier
		} else if id, isIdentifier := lookup.IdentifierToId[userInput]; isIdentifier {
			// Input is an identifier (name/email), convert to id
			itemId = id
			itemIdentifier = userInput
		} else {
			// Validate if item exists but has invalid id
			isValid, errMsg := validateItem(userInput)
			if !isValid {
				return nil, fmt.Errorf("%s", errMsg)
			}
			return nil, fmt.Errorf("%s '%s' not found. Please provide a valid %s identifier or %s Id", itemType, userInput, itemType, itemType)
		}

		if itemId == "" {
			return nil, fmt.Errorf("%s '%s' resulted in an empty %s_id. This should not happen - please report this issue", itemType, userInput, itemType)
		}

		// Check for duplicates
		if originalInput, exists := seenIds[itemId]; exists {
			return nil, fmt.Errorf("duplicate %s detected: '%s' and '%s' both map to the same %s Id '%s' (%s identifier: '%s')",
				itemType, originalInput, userInput, itemType, itemId, itemType, itemIdentifier)
		}

		seenIds[itemId] = userInput
		result[itemId] = userInput
	}

	return result, nil
}

// RestoreUserInputFormatFromApiResponse converts API response identifiers back to the format
// that the user originally provided in their Terraform configuration.
//
// This function is critical for preventing false diffs in Terraform plans. When users provide
// identifiers in different formats (e.g., IDs vs names/emails), the API typically returns
// a standardized format (usually names/emails). This function ensures the Terraform state
// preserves the original user input format, preventing unnecessary plan changes.
//
// How it works:
//  1. Takes API response identifiers (typically names/emails) and current Terraform state
//  2. Builds lookup maps to convert between IDs and identifiers (names/emails)
//  3. For each API identifier, finds the corresponding ID
//  4. Looks up the original format from state (what user originally provided)
//  5. Returns a Set preserving the original format (IDs or names/emails)
//
// Examples:
//   - User provided user_id "123" → API returns "user@example.com" → Function returns "123"
//   - User provided email "user@example.com" → API returns "user@example.com" → Function returns "user@example.com"
//   - New item added outside Terraform → Function returns the API identifier (name/email)
//
// This is a generic function that works for roles, users, and teams. Type-specific wrappers
// are provided in process_roles.go, process_users.go, and process_teams.go.
//
// Parameters:
//   - apiIdentifiers: Identifiers returned by the API (typically names/emails)
//   - currentState: Current Terraform state (what user originally provided)
//   - itemType: Type of item ("role", "user", or "team") - used for error messages
//   - fetchCommand: Commander CLI command to fetch all items for lookup map building
//   - parseFunc: Function to parse the API response into the appropriate type
//   - buildLookupFunc: Function to build lookup maps from parsed data
//
// Returns:
//   - types.Set: Set of identifiers in the original user input format
//   - error: Error if fetching or parsing fails
func RestoreUserInputFormatFromApiResponse(
	ctx context.Context,
	apiManager *api.ApiManager,
	apiIdentifiers []string, // Names/emails from API response
	currentState types.Set, // Current state (what user originally provided)
	itemType string, // "role", "user", or "team"
	fetchCommand string, // "enterprise-info -r --format json" or "enterprise-info -u --format json" or "enterprise-info -t --format json"
	parseFunc func(interface{}) (interface{}, error), // ParseRolesResponse, ParseUsersResponse, or ParseTeamsResponse
	buildLookupFunc func(interface{}) LookupMaps, // BuildRoleLookupMaps, BuildUserLookupMaps, or BuildTeamLookupMaps
) (types.Set, error) {
	// Handle empty/null cases
	if len(apiIdentifiers) == 0 {
		return types.SetNull(types.StringType), nil
	}

	// Build a map of current state values for quick lookup
	stateMap := make(map[string]bool)
	stateIdToOriginal := make(map[string]string)         // ID -> original value from state
	stateIdentifierToOriginal := make(map[string]string) // name/email -> original value from state

	if !currentState.IsNull() {
		var stateValues []string
		currentState.ElementsAs(ctx, &stateValues, false)
		for _, val := range stateValues {
			val = strings.TrimSpace(val)
			if val != "" {
				stateMap[val] = true
			}
		}
	}

	// Fetch all items to build lookup maps
	itemsResp, err := apiManager.ExecuteCommand(ctx, fetchCommand, fmt.Sprintf("Unable to fetch enterprise %ss", itemType))
	if err != nil {
		return types.SetNull(types.StringType), fmt.Errorf("failed to fetch %ss: %w", itemType, err)
	}

	// Parse the response
	itemsData, err := parseFunc(itemsResp.Data)
	if err != nil {
		return types.SetNull(types.StringType), fmt.Errorf("failed to parse %ss: %w", itemType, err)
	}

	// Build lookup maps
	lookup := buildLookupFunc(itemsData)

	// Build reverse maps: for each value in state, map its ID to the original value
	for stateVal := range stateMap {
		// Check if state value is an ID
		if identifier, isId := lookup.IdToIdentifier[stateVal]; isId {
			// State has ID, map ID -> original ID
			stateIdToOriginal[stateVal] = stateVal
			// Also map the identifier (name/email) -> original ID
			stateIdentifierToOriginal[identifier] = stateVal
		} else if id, isIdentifier := lookup.IdentifierToId[stateVal]; isIdentifier {
			// State has name/email, map ID -> original name/email
			stateIdToOriginal[id] = stateVal
			// Also map identifier -> original identifier
			stateIdentifierToOriginal[stateVal] = stateVal
		}
	}

	// Convert API identifiers to original format
	var resultStrings []string
	seen := make(map[string]bool)

	for _, apiIdentifier := range apiIdentifiers {
		apiIdentifier = strings.TrimSpace(apiIdentifier)
		if apiIdentifier == "" {
			continue
		}

		// Find the ID for this API identifier
		var id string
		if _, isId := lookup.IdToIdentifier[apiIdentifier]; isId {
			// API returned an ID (unlikely but handle it)
			id = apiIdentifier
		} else if foundId, isIdentifier := lookup.IdentifierToId[apiIdentifier]; isIdentifier {
			// API returned name/email, get its ID
			id = foundId
		} else {
			// Not found - item might have been deleted, skip it
			continue
		}

		// Find the original format from state
		var originalValue string
		if original, found := stateIdToOriginal[id]; found {
			// State had this ID, use the original format (could be ID or name/email)
			originalValue = original
		} else if original, found := stateIdentifierToOriginal[apiIdentifier]; found {
			// State had this name/email, use it
			originalValue = original
		} else {
			// New item added outside Terraform - use the identifier from API (name/email)
			originalValue = apiIdentifier
		}

		// Avoid duplicates
		if !seen[originalValue] {
			resultStrings = append(resultStrings, originalValue)
			seen[originalValue] = true
		}
	}

	// Convert to types.Set
	if len(resultStrings) == 0 {
		return types.SetNull(types.StringType), nil
	}

	resultElements := make([]types.String, len(resultStrings))
	for i, val := range resultStrings {
		resultElements[i] = types.StringValue(val)
	}

	resultSet, diags := types.SetValueFrom(ctx, types.StringType, resultElements)
	if diags.HasError() {
		return types.SetNull(types.StringType), fmt.Errorf("failed to create set: %v", diags)
	}

	return resultSet, nil
}

func FetchEnterpriseNodeByNameOrId(ctx context.Context, apiManager *api.ApiManager, nodeNameOrId string) (*EnterpriseNodeResponse, error) {
	// node can be id or name

	command := fmt.Sprintf("enterprise-info -n -v --format json --node '%s' --columns='isolated,parent_node,parent_id'", nodeNameOrId)

	apiResp, err := apiManager.ExecuteCommand(ctx, command, "Failed to retrieve enterprise node information")
	if err != nil {
		if errors.Is(err, api.ErrResourceNotFound) {
			return nil, nil
		}
		return nil, err
	}

	// Parse the JSON response - it's an array of node objects
	var nodes []EnterpriseNodeResponse

	if err := UnmarshalApiResponse(apiResp.Data, &nodes); err != nil {
		return nil, fmt.Errorf("unable to parse enterprise nodes list from API response: %w", err)
	}

	// Find the node matching state.Id (which is the node name)
	var nodeInfo *EnterpriseNodeResponse

	for i := range nodes {
		// convert node id to string to compare with id
		if strconv.Itoa(nodes[i].NodeId) == nodeNameOrId || nodes[i].Name == nodeNameOrId {
			nodeInfo = &nodes[i]
			break
		}
	}

	// Node not in list - resource was likely deleted outside Terraform.
	// Return (nil, nil) so Read can remove it from state instead of erroring.
	if nodeInfo == nil {
		return nil, nil
	}

	return nodeInfo, nil
}

func FetchEnterpriseRoleByNameOrId(ctx context.Context, apiManager *api.ApiManager, roleNameOrId string) (*EnterpriseRoleResponse, error) {
	// Build the Commander command string
	command := fmt.Sprintf("enterprise-info '%s' -r --format json --columns='visible_below,default_role,admin,node,users,teams,managed_nodes_permissions,enforcements' -q", roleNameOrId)

	apiResp, err := apiManager.ExecuteCommand(ctx, command, "Failed to retrieve enterprise role information")
	if err != nil {
		if errors.Is(err, api.ErrResourceNotFound) {
			return nil, nil
		}
		return nil, err
	}

	// Parse the JSON response - it's an array of role objects
	var roles []EnterpriseRoleResponse

	// Unmarshal API response into roles struct
	if err := UnmarshalApiResponse(apiResp.Data, &roles); err != nil {
		return nil, fmt.Errorf("unable to parse enterprise roles list from API response: %w", err)
	}

	// Find the role matching roleNameOrId (by ID or name)
	var roleInfo *EnterpriseRoleResponse
	for i := range roles {
		if strconv.Itoa(roles[i].RoleId) == roleNameOrId || roles[i].Name == roleNameOrId {
			roleInfo = &roles[i]
			break
		}
	}

	if roleInfo == nil {
		return nil, nil
	}

	return roleInfo, nil
}

func FetchEnterpriseTeamByNameOrId(ctx context.Context, apiManager *api.ApiManager, teamNameOrId string) (*EnterpriseTeamResponse, error) {
	// Build command to get enterprise team info
	command := fmt.Sprintf("enterprise-info '%s' -t --format json --columns='users,roles,restricts,node' -q", teamNameOrId)

	// Execute the command
	apiResp, err := apiManager.ExecuteCommand(ctx, command, "Failed to retrieve enterprise team information")
	if err != nil {
		if errors.Is(err, api.ErrResourceNotFound) {
			return nil, nil
		}
		return nil, err
	}

	// Parse the JSON response - it's an array of team objects
	var teams []EnterpriseTeamResponse

	// Unmarshal API response into teams struct
	if err := UnmarshalApiResponse(apiResp.Data, &teams); err != nil {
		return nil, fmt.Errorf("unable to parse enterprise teams list from API response: %w", err)
	}

	// Find the team matching teamNameOrId (by ID or name)
	var teamInfo *EnterpriseTeamResponse
	for i := range teams {
		if teams[i].TeamUid == teamNameOrId || teams[i].Name == teamNameOrId {
			teamInfo = &teams[i]
			break
		}
	}

	if teamInfo == nil {
		return nil, nil
	}

	return teamInfo, nil
}

func FetchEnterpriseUserByEmailOrId(ctx context.Context, apiManager *api.ApiManager, emailOrId string) (*EnterpriseUserResponse, error) {
	// Build the Commander command string
	command := fmt.Sprintf("enterprise-info '%s' -u --format json --columns='name,status,node,teams,roles,alias,job_title' -q", emailOrId)

	apiResp, err := apiManager.ExecuteCommand(ctx, command, "Failed to retrieve enterprise user information")
	if err != nil {
		if errors.Is(err, api.ErrResourceNotFound) {
			return nil, nil
		}
		return nil, err
	}

	// Parse the JSON response - it's an array of user objects
	var users []EnterpriseUserResponse

	// Unmarshal API response into users struct
	if err := UnmarshalApiResponse(apiResp.Data, &users); err != nil {
		return nil, fmt.Errorf("unable to parse enterprise users list from API response: %w", err)
	}

	// Find the user matching state.Id
	var userInfo *EnterpriseUserResponse
	for i := range users {
		if strconv.Itoa(users[i].UserId) == emailOrId || users[i].Email == emailOrId {
			userInfo = &users[i]
			break
		}
	}

	if userInfo == nil {
		return nil, nil
	}

	return userInfo, nil
}

func FetchManagedCompanyByNameOrId(ctx context.Context, apiManager *api.ApiManager, nameOrId string) (*ManagedCompanyResponse, error) {
	// Build command to get all companies info
	command := fmt.Sprintf("msp-info -m '%s' --format json -v", nameOrId)

	apiResp, err := apiManager.ExecuteCommand(ctx, command, "Failed to retrieve managed company information")
	if err != nil {
		if errors.Is(err, api.ErrResourceNotFound) {
			return nil, nil
		}
		return nil, err
	}

	// Parse the JSON response - it's an array of company objects
	var companies []ManagedCompanyResponse

	if err := UnmarshalApiResponse(apiResp.Data, &companies); err != nil {
		return nil, fmt.Errorf("unable to parse managed companies list from API response: %w", err)
	}

	var companyInfo *ManagedCompanyResponse

	for i := range companies {
		if strconv.Itoa(companies[i].CompanyId) == nameOrId || companies[i].CompanyName == nameOrId {
			companyInfo = &companies[i]
			break
		}
	}

	if companyInfo == nil {
		return nil, nil
	}

	return companyInfo, nil
}

// FetchEnterpriseScimById retrieves a single SCIM configuration by scim_id.
// The API returns a single object. Returns (nil, nil) if not found so Read can remove from state.
func FetchEnterpriseScimById(ctx context.Context, apiManager *api.ApiManager, scimId string) (*EnterpriseScimResponse, error) {
	command := fmt.Sprintf("scim view '%s' --format json", scimId)

	apiResp, err := apiManager.ExecuteCommand(ctx, command, "Failed to retrieve enterprise SCIM information")
	if err != nil {
		return nil, err
	}

	var scimInfo EnterpriseScimResponse
	if err := UnmarshalApiResponse(apiResp.Data, &scimInfo); err != nil {
		return nil, fmt.Errorf("unable to parse enterprise SCIM from API response: %w", err)
	}

	// Treat empty or zero scim_id as not found
	if scimInfo.ScimID == 0 {
		return nil, nil
	}

	return &scimInfo, nil
}

func FetchEpmPolicyById(ctx context.Context, apiManager *api.ApiManager, policyId string) (*EpmPolicyResponse, error) {
	command := fmt.Sprintf("epm policy view '%s' --format json", policyId)
	apiResp, err := apiManager.ExecuteCommand(ctx, command, "Unable to read EPM policy")
	if err != nil {
		if errors.Is(err, api.ErrResourceNotFound) {
			return nil, nil
		}
		return nil, err
	}

	// Parse the JSON response - it's an array of company objects
	var policyInfo *EpmPolicyResponse

	if err := UnmarshalApiResponse(apiResp.Data, &policyInfo); err != nil {
		return nil, fmt.Errorf("unable to parse EPM policy from API response: %w", err)
	}

	if policyInfo == nil {
		return nil, nil
	}

	return policyInfo, nil
}

// ParseManagedCompanyImportID parses an import ID that may be "resource_id" or "managed_company,resource_id".
// resourceName is used in error messages (e.g. "node" or "role").
// Returns resourceID, managedCompany (empty if not in ID), and any diagnostics.
func ParseManagedCompanyImportID(importID string, resourceName string) (resourceID, managedCompany string, diags diag.Diagnostics) {
	importID = strings.TrimSpace(importID)
	if importID == "" {
		diags = append(diags, diag.NewErrorDiagnostic(
			ERR_MSG_INVALID_IMPORT_ID,
			"Import ID cannot be empty. Use: (1) "+resourceName+" name or "+resourceName+" ID alone; or (2) for a "+resourceName+" in a managed company, use \"managed_company_name_or_id,"+resourceName+"_name_or_id\" (comma-separated).",
		))
		return "", "", diags
	}

	if parts := strings.SplitN(importID, ",", 2); len(parts) == 2 {
		managedCompany = strings.TrimSpace(parts[0])
		resourceID = strings.TrimSpace(parts[1])
	} else {
		resourceID = importID
	}

	if resourceID == "" {
		diags = append(diags, diag.NewErrorDiagnostic(
			ERR_MSG_INVALID_IMPORT_ID,
			"When using managed company format \"managed_company_name_or_id,"+resourceName+"\", the "+resourceName+" part cannot be empty.",
		))
		return "", "", diags
	}

	return resourceID, managedCompany, diags
}

// ManagingNodesMapElemType is the object type for each entry in the managing_nodes map (privileges set + cascade bool).
// Used by enterprise role resource and data source.
var ManagingNodesMapElemType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"privileges": types.SetType{ElemType: types.StringType},
		"cascade":    types.BoolType,
	},
}

// MapManagedNodesPermissionsToState converts API ManagedNodesPermissions to a types.Map for state.
// Map key is the node name (or node ID if name empty). Value is object with privileges (Set) and cascade (Bool).
func MapManagedNodesPermissionsToState(_ context.Context, perms []ManagedNodePermission) (types.Map, error) {
	if len(perms) == 0 {
		return types.MapNull(ManagingNodesMapElemType), nil
	}
	elements := make(map[string]attr.Value)
	for _, p := range perms {
		key := p.NodeName
		if key == "" {
			key = strconv.FormatInt(p.NodeId, 10)
		}
		privilegeElems := make([]attr.Value, len(p.Privileges))
		for i, pr := range p.Privileges {
			privilegeElems[i] = types.StringValue(pr)
		}
		privilegesSet := types.SetValueMust(types.StringType, privilegeElems)
		obj := types.ObjectValueMust(
			ManagingNodesMapElemType.AttrTypes,
			map[string]attr.Value{
				"privileges": privilegesSet,
				"cascade":    types.BoolValue(p.Cascade),
			},
		)
		elements[key] = obj
	}
	mapVal, diags := types.MapValue(ManagingNodesMapElemType, elements)
	if diags.HasError() {
		return types.MapNull(ManagingNodesMapElemType), fmt.Errorf("failed to build managing_nodes map: %v", diags)
	}
	return mapVal, nil
}

// CanonicalizeGeneratedPasswordComplexityJSON parses the JSON string and re-encodes it with sorted keys
// so that semantically equal values compare equal (avoids perpetual diff from whitespace/key order).
func CanonicalizeGeneratedPasswordComplexityJSON(s string) string {
	if s == "" {
		return s
	}
	var arr []map[string]interface{}
	if err := json.Unmarshal([]byte(s), &arr); err != nil {
		return s
	}
	out, err := json.Marshal(arr)
	if err != nil {
		return s
	}
	return string(out)
}

const (
	restrictSharingAll                = "RESTRICT_SHARING_ALL"
	restrictSharingAllIncoming        = "RESTRICT_SHARING_ALL_INCOMING"
	restrictSharingAllOutgoing        = "RESTRICT_SHARING_ALL_OUTGOING"
	restrictSharingEnterprise         = "RESTRICT_SHARING_ENTERPRISE"
	restrictSharingEnterpriseIncoming = "RESTRICT_SHARING_ENTERPRISE_INCOMING"
	restrictSharingEnterpriseOutgoing = "RESTRICT_SHARING_ENTERPRISE_OUTGOING"
)

// MapEnforcementsToState converts API enforcements (key -> string value) to a types.Map for state.
// Keys are normalized to UPPER_SNAKE_CASE. GENERATED_PASSWORD_COMPLEXITY value is canonicalized.
// Commander expands RESTRICT_SHARING_ALL into _INCOMING and _OUTGOING; when both are "true" we collapse back to RESTRICT_SHARING_ALL for stable plan.
func MapEnforcementsToState(enforcements map[string]string, GeneratedPasswordComplexityKey string) (types.Map, error) {
	if len(enforcements) == 0 {
		return types.MapNull(types.StringType), nil
	}
	elements := make(map[string]attr.Value)
	for key, val := range enforcements {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		normalizedKey := strings.ToUpper(strings.ReplaceAll(key, "-", "_"))
		if normalizedKey == GeneratedPasswordComplexityKey {
			val = CanonicalizeGeneratedPasswordComplexityJSON(val)
		}
		elements[normalizedKey] = types.StringValue(val)
	}
	// Collapse RESTRICT_SHARING_ALL_INCOMING + RESTRICT_SHARING_ALL_OUTGOING (both "true") -> RESTRICT_SHARING_ALL = "true"
	if incoming, okIn := elements[restrictSharingAllIncoming]; okIn {
		if outgoing, okOut := elements[restrictSharingAllOutgoing]; okOut {
			if sIn, ok := incoming.(types.String); ok && sIn.ValueString() == "true" {
				if sOut, ok := outgoing.(types.String); ok && sOut.ValueString() == "true" {
					elements[restrictSharingAll] = types.StringValue("true")
					delete(elements, restrictSharingAllIncoming)
					delete(elements, restrictSharingAllOutgoing)
				}
			}

		}
	}
	// Collapse RESTRICT_SHARING_ENTERPRISE_INCOMING + RESTRICT_SHARING_ENTERPRISE_OUTGOING (both "true") -> RESTRICT_SHARING_ENTERPRISE = "true"
	if incoming, okIn := elements[restrictSharingEnterpriseIncoming]; okIn {
		if outgoing, okOut := elements[restrictSharingEnterpriseOutgoing]; okOut {
			if sIn, ok := incoming.(types.String); ok && sIn.ValueString() == "true" {
				if sOut, ok := outgoing.(types.String); ok && sOut.ValueString() == "true" {
					elements[restrictSharingEnterprise] = types.StringValue("true")
					delete(elements, restrictSharingEnterpriseIncoming)
					delete(elements, restrictSharingEnterpriseOutgoing)
				}
			}
		}
	}
	if len(elements) == 0 {
		return types.MapNull(types.StringType), nil
	}
	mapVal, diags := types.MapValue(types.StringType, elements)
	if diags.HasError() {
		return types.MapNull(types.StringType), fmt.Errorf("failed to build enforcement_policies map: %v", diags)
	}
	return mapVal, nil
}

// ImmutableAttributeCheck represents one "plan vs state must be equal" check.
type ImmutableAttributeCheck struct {
	PlanValue  attr.Value
	StateValue attr.Value
	Summary    string
	Detail     string
}

// AddErrorIfImmutableAttributesChanged checks each attribute; for any where plan != state
// it adds a diagnostic error to diags. Returns diags so caller can Append() or check HasError().
func RestrictAttributeUpdate(diags *diag.Diagnostics, checks []ImmutableAttributeCheck) {
	for _, c := range checks {
		if c.PlanValue == nil || c.StateValue == nil {
			continue
		}
		if !c.PlanValue.Equal(c.StateValue) {
			diags.AddError(c.Summary, c.Detail)
		}
	}
}

// VaultRecordListEntry is one element returned by `list --format json`.
type VaultRecordListEntry struct {
	RecordUID   string `json:"record_uid"`
	Type        string `json:"type"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Shared      bool   `json:"shared"`
}

// ValidateVaultRecordIdentifiers runs `list --format json` and ensures each non-empty identifier
// matches some record's record_uid or title (exact string match). Duplicate identifiers are validated once.
func ValidateVaultRecordIdentifiers(ctx context.Context, apiManager *api.ApiManager, identifiers []string) error {
	entries, err := FetchVaultRecordList(ctx, apiManager)
	if err != nil {
		return err
	}
	return ValidateRecordIdentifiersAgainstList(identifiers, entries)
}

// FetchVaultRecordList returns vault records from Commander `list --format json`.
func FetchVaultRecordList(ctx context.Context, apiManager *api.ApiManager) ([]VaultRecordListEntry, error) {
	// Commander CLI: list vault records as JSON (array of objects with record_uid, title, etc.).
	const CmdListVaultRecordsJSON = "list --format json"

	apiResp, err := apiManager.ExecuteCommand(ctx, CmdListVaultRecordsJSON, ErrOpListVaultRecords)
	if err != nil {
		return nil, err
	}
	var entries []VaultRecordListEntry
	if apiResp.Data == nil {
		return entries, nil
	}
	if err := UnmarshalApiResponse(apiResp.Data, &entries); err != nil {
		return nil, fmt.Errorf("%s: %w", ErrOpListVaultRecords, err)
	}
	return entries, nil
}

// ValidateRecordIdentifiersAgainstList returns an error if any non-empty identifier does not match
// any entry's RecordUID or Title (exact match).
func ValidateRecordIdentifiersAgainstList(identifiers []string, entries []VaultRecordListEntry) error {
	seen := make(map[string]struct{})
	var invalid []string
	for _, raw := range identifiers {
		id := strings.TrimSpace(raw)
		if id == "" {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		if !recordIdentifierMatchesList(id, entries) {
			invalid = append(invalid, id)
		}
	}
	if len(invalid) == 0 {
		return nil
	}
	return fmt.Errorf(
		"invalid record reference(s): %s — each must match a vault record UID or title",
		strings.Join(invalid, ", "),
	)
}

func recordIdentifierMatchesList(ref string, entries []VaultRecordListEntry) bool {
	for _, e := range entries {
		if ref == e.RecordUID {
			return true
		}
		if ref == e.Title {
			return true
		}
	}
	return false
}

// ExtractFolderValue resolves the folder value from the API response against the
// current Terraform state. If the user originally specified a UID or path that
// matches the API response, the state value is preserved so Terraform does not
// show a spurious diff. Otherwise the folder path from the response is used.
// Path comparison normalizes spaces around "/" so "Test / My Folder" matches "Test/My Folder".
func ExtractFolderValue(folderResponse *FolderLocationResponse, stateFolder types.String) types.String {
	if folderResponse == nil || (folderResponse.UID == "" && strings.TrimSpace(folderResponse.Path) == "/") || strings.TrimSpace(folderResponse.Path) == "/" || strings.TrimSpace(folderResponse.Path) == "" {
		return types.StringNull()
	}
	if !stateFolder.IsNull() && !stateFolder.IsUnknown() {
		sv := strings.TrimSpace(stateFolder.ValueString())
		if sv == strings.TrimSpace(folderResponse.UID) || NormalizeFolderPath(sv) == NormalizeFolderPath(folderResponse.Path) {
			return stateFolder
		}
	}
	return types.StringValue(folderResponse.Path)
}

// NormalizeFolderPath removes spaces around "/" separators so that
// "Test / My Folder" and "Test/My Folder" compare as equal.
func NormalizeFolderPath(p string) string {
	parts := strings.Split(p, "/")
	for i, part := range parts {
		parts[i] = strings.TrimSpace(part)
	}
	return strings.Join(parts, "/")
}

// StringOrNull returns types.StringValue when s is non-empty after trim, otherwise types.StringNull.
func StringOrNull(s string) types.String {
	if strings.TrimSpace(s) == "" {
		return types.StringNull()
	}
	return types.StringValue(s)
}

// MergeResourceAttributes combines resource attribute maps; later maps override earlier keys.
func MergeResourceAttributes(maps ...map[string]schema.Attribute) map[string]schema.Attribute {
	result := map[string]schema.Attribute{}
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

// MergeDataSourceAttributes combines data source attribute maps; later maps override earlier keys.
func MergeDataSourceAttributes(maps ...map[string]dschema.Attribute) map[string]dschema.Attribute {
	result := map[string]dschema.Attribute{}
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

// QuoteShellSingle wraps s for use as a single-quoted shell argument (bash-style escaping of ').
func QuoteShellSingle(s string) string {
	return `'` + strings.ReplaceAll(s, `'`, `'"'"'`) + `'`
}

// MoveNsfFromSourceToDestination move's a nsf-record/nsf-folder when plan and state folder paths differ.
func MoveNsfFromSourceToDestination(ctx context.Context, apiManager *api.ApiManager, uID string, planFolderData string, stateFolderData string) error {
	if planFolderData == stateFolderData {
		return nil
	}

	dest := planFolderData
	if dest == "" {
		dest = "root"
	}

	command := fmt.Sprintf("%s '%s' '%s'", CmdNsfMove, uID, dest)
	_, err := apiManager.ExecuteCommand(ctx, command, ErrSummaryNsfMoveRecordFailed)
	return err
}
