package enterpriseteam

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

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

type RoleInfo struct {
	RoleId int    `json:"role_id"`
	Name   string `json:"name"`
}

type UserInfo struct {
	UserId int    `json:"user_id"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

type NodeInfo struct {
	NodeId int    `json:"node_id"`
	Name   string `json:"name"`
}

// lookupMaps holds the mappings between identifiers (name/email) and IDs
type lookupMaps struct {
	identifierToId map[string]string // identifier (name/email) -> id
	idToIdentifier map[string]string // id -> identifier (name/email)
}

// parseTeamsResponse parses the JSON response from enterprise-info -t command
func parseTeamsResponse(data interface{}) ([]TeamInfo, error) {
	var teams []TeamInfo

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("unable to process the response from Keeper Commander Service Mode API: %w", err)
	}

	if err := json.Unmarshal(dataBytes, &teams); err != nil {
		return nil, fmt.Errorf("unable to parse enterprise teams list from Service Mode API response: %w", err)
	}

	return teams, nil
}

// findTeamUidByName finds a team UID by name from the teams list
func findTeamUidByName(teams []TeamInfo, teamName string) (string, error) {
	for _, team := range teams {
		if team.Name == teamName {
			return team.TeamUid, nil
		}
	}
	return "", fmt.Errorf("enterprise team with name '%s' not found in the response", teamName)
}

// fetchTeamUidByName fetches the team UID by name using the API
func fetchTeamUidByName(ctx context.Context, apiManager *api.ApiManager, teamName string) (string, error) {
	teamsResp, err := apiManager.ExecuteCommand(ctx, fmt.Sprintf("enterprise-info -t --format json -q '%s'", teamName), "Unable to fetch enterprise team ID")
	if err != nil {
		return "", err
	}

	teams, err := parseTeamsResponse(teamsResp.Data)
	if err != nil {
		return "", err
	}

	return findTeamUidByName(teams, teamName)
}

func parseRolesResponse(data interface{}) ([]RoleInfo, error) {
	var roles []RoleInfo

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("unable to process the response from Keeper Commander Service Mode API: %w", err)
	}

	if err := json.Unmarshal(dataBytes, &roles); err != nil {
		return nil, fmt.Errorf("unable to parse enterprise roles list from Service Mode API response: %w", err)
	}

	return roles, nil
}

func parseUsersResponse(data interface{}) ([]UserInfo, error) {
	var users []UserInfo

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("unable to process the response from Keeper Commander Service Mode API: %w", err)
	}

	if err := json.Unmarshal(dataBytes, &users); err != nil {
		return nil, fmt.Errorf("unable to parse enterprise users list from Service Mode API response: %w", err)
	}

	return users, nil
}

func parseNodesResponse(data interface{}) ([]NodeInfo, error) {
	var nodes []NodeInfo

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("unable to process the response from Keeper Commander Service Mode API: %w", err)
	}

	if err := json.Unmarshal(dataBytes, &nodes); err != nil {
		return nil, fmt.Errorf("unable to parse enterprise nodes list from Service Mode API response: %w", err)
	}

	return nodes, nil
}

// buildRoleLookupMaps creates lookup maps from API response
func buildRoleLookupMaps(rolesRespData []RoleInfo) lookupMaps {
	identifierToId := make(map[string]string)
	idToIdentifier := make(map[string]string)

	for _, role := range rolesRespData {
		if role.RoleId > 0 && role.Name != "" {
			roleIdStr := strconv.Itoa(role.RoleId)
			identifierToId[role.Name] = roleIdStr
			idToIdentifier[roleIdStr] = role.Name
		}
	}

	return lookupMaps{
		identifierToId: identifierToId,
		idToIdentifier: idToIdentifier,
	}
}

// buildUserLookupMaps creates lookup maps from API response
func buildUserLookupMaps(usersRespData []UserInfo) lookupMaps {
	identifierToId := make(map[string]string)
	idToIdentifier := make(map[string]string)

	for _, user := range usersRespData {
		if user.UserId > 0 && user.Email != "" {
			userIdStr := strconv.Itoa(user.UserId)
			identifierToId[user.Email] = userIdStr
			idToIdentifier[userIdStr] = user.Email
		}
	}

	return lookupMaps{
		identifierToId: identifierToId,
		idToIdentifier: idToIdentifier,
	}
}

// convertItemsToIdMap is a generic function that converts a types.Set to a map of id -> original input
// It works for both roles and users by accepting lookup maps and validation functions
func convertItemsToIdMap(
	items types.Set,
	lookup lookupMaps,
	itemType string, // "role" or "user"
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
		itemStr := itemElem.(types.String)
		userInput := itemStr.ValueString()

		if userInput == "" {
			continue
		}

		var itemId string
		var itemIdentifier string

		// Check if input is an id
		if existingIdentifier, isId := lookup.idToIdentifier[userInput]; isId {
			itemId = userInput
			itemIdentifier = existingIdentifier
		} else if id, isIdentifier := lookup.identifierToId[userInput]; isIdentifier {
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

// convertRolesToIdMap converts a types.Set of roles to a map of role_id -> original input
func convertRolesToIdMap(roles types.Set, lookup lookupMaps, rolesRespData []RoleInfo) (map[string]string, error) {
	validateRole := func(userInput string) (bool, string) {
		for _, role := range rolesRespData {
			if role.Name == userInput && role.RoleId <= 0 {
				return false, "role '" + userInput + "' exists but has no valid role_id. This role cannot be used"
			}
		}
		return true, ""
	}

	return convertItemsToIdMap(roles, lookup, "role", validateRole)
}

// convertUsersToIdMap converts a types.Set of users to a map of user_id -> original input
func convertUsersToIdMap(users types.Set, lookup lookupMaps, usersRespData []UserInfo) (map[string]string, error) {
	validateUser := func(userInput string) (bool, string) {
		for _, user := range usersRespData {
			if user.Email == userInput && user.UserId <= 0 {
				return false, "user '" + userInput + "' exists but has no valid user_id. This user cannot be used"
			}
		}
		return true, ""
	}

	return convertItemsToIdMap(users, lookup, "user", validateUser)
}

// fetchAndProcessRoles processes roles for both create and update operations
// For create: stateRoles should be null/empty, planRoles contains roles to add
// For update: compares stateRoles (old) with planRoles (new) to determine additions and removals
// Returns a string with -ar "role_id" for additions and -rr "role_id" for removals
func fetchAndProcessRoles(ctx context.Context, apiManager *api.ApiManager, stateRoles types.Set, planRoles types.Set) (string, error) {
	// Early return if both are empty/null
	if (stateRoles.IsNull() || len(stateRoles.Elements()) == 0) &&
		(planRoles.IsNull() || len(planRoles.Elements()) == 0) {
		return "", nil
	}

	// Fetch roles from API
	rolesResp, err := apiManager.ExecuteCommand(ctx, "enterprise-info -r --format json", "Unable to fetch enterprise roles")
	if err != nil {
		return "", err
	}

	// Parse the roles response
	rolesRespData, err := parseRolesResponse(rolesResp.Data)
	if err != nil {
		return "", err
	}

	// Build lookup maps
	lookup := buildRoleLookupMaps(rolesRespData)

	// Convert state roles to role_id map (old roles)
	stateRoleIdMap, err := convertRolesToIdMap(stateRoles, lookup, rolesRespData)
	if err != nil {
		return "", err
	}

	// Convert plan roles to role_id map (new roles)
	planRoleIdMap, err := convertRolesToIdMap(planRoles, lookup, rolesRespData)
	if err != nil {
		return "", err
	}

	// Early return if no changes
	if len(stateRoleIdMap) == 0 && len(planRoleIdMap) == 0 {
		return "", nil
	}

	// Find roles to add and remove
	var parts []string

	// Add roles that are in plan but not in state
	for roleId := range planRoleIdMap {
		if _, exists := stateRoleIdMap[roleId]; !exists {
			parts = append(parts, fmt.Sprintf("-ar '%s'", roleId))
		}
	}

	// Remove roles that are in state but not in plan
	for roleId := range stateRoleIdMap {
		if _, exists := planRoleIdMap[roleId]; !exists {
			parts = append(parts, fmt.Sprintf("-rr '%s'", roleId))
		}
	}

	if len(parts) == 0 {
		return "", nil
	}

	return strings.Join(parts, " "), nil
}

// fetchAndProcessUsers processes users for both create and update operations
// For create: stateUsers should be null/empty, planUsers contains users to add
// For update: compares stateUsers (old) with planUsers (new) to determine additions and removals
// Returns a string with -au "user_id" for additions and -ru "user_id" for removals
func fetchAndProcessUsers(ctx context.Context, apiManager *api.ApiManager, stateUsers types.Set, planUsers types.Set) (string, error) {
	// Early return if both are empty/null
	if (stateUsers.IsNull() || len(stateUsers.Elements()) == 0) &&
		(planUsers.IsNull() || len(planUsers.Elements()) == 0) {
		return "", nil
	}

	// Fetch users from API
	usersResp, err := apiManager.ExecuteCommand(ctx, "enterprise-info -u --format json", "Unable to fetch enterprise users")
	if err != nil {
		return "", err
	}

	// Parse the users response
	usersRespData, err := parseUsersResponse(usersResp.Data)
	if err != nil {
		return "", err
	}

	// Build lookup maps
	lookup := buildUserLookupMaps(usersRespData)

	// Create a map of user_id -> UserInfo for status checking
	userIdToUserInfo := make(map[string]UserInfo)
	for _, user := range usersRespData {
		if user.UserId > 0 {
			userIdStr := strconv.Itoa(user.UserId)
			userIdToUserInfo[userIdStr] = user
		}
	}

	// Convert state users to user_id map (old users)
	stateUserIdMap, err := convertUsersToIdMap(stateUsers, lookup, usersRespData)
	if err != nil {
		return "", err
	}

	// Convert plan users to user_id map (new users)
	planUserIdMap, err := convertUsersToIdMap(planUsers, lookup, usersRespData)
	if err != nil {
		return "", err
	}

	// Early return if no changes
	if len(stateUserIdMap) == 0 && len(planUserIdMap) == 0 {
		return "", nil
	}

	// Find users to add and remove
	var parts []string

	// Add users that are in plan but not in state
	for userId := range planUserIdMap {
		if _, exists := stateUserIdMap[userId]; !exists {
			// Check if user has "Invited" status
			if userInfo, exists := userIdToUserInfo[userId]; exists {
				if userInfo.Status == "Invited" {
					userIdentifier := planUserIdMap[userId] // Get original input (email or user_id)
					return "", fmt.Errorf("user '%s' has status 'Invited'. Users must accept invitation before being added to a team", userIdentifier)
				}
			}
			parts = append(parts, fmt.Sprintf("-au '%s'", userId))
		}
	}

	// Remove users that are in state but not in plan
	for userId := range stateUserIdMap {
		if _, exists := planUserIdMap[userId]; !exists {
			parts = append(parts, fmt.Sprintf("-ru '%s'", userId))
		}
	}

	if len(parts) == 0 {
		return "", nil
	}

	return strings.Join(parts, " "), nil
}

func buildEnterpriseTeamAddCommand(data EnterpriseTeamResourceModel) string {
	var parts []string

	parts = append(parts, "enterprise-team")

	// Required parameters
	parts = append(parts, fmt.Sprintf("--add '%s'", data.Name.ValueString()))

	// Optional parameters
	if !data.RestrictEdit.IsNull() {
		parts = append(parts, "--restrict-edit on")
	}

	if !data.RestrictShare.IsNull() {
		parts = append(parts, "--restrict-share on")
	}

	if !data.RestrictView.IsNull() {
		parts = append(parts, "--restrict-view on")
	}

	if !data.Node.IsNull() {
		parts = append(parts, fmt.Sprintf("--node '%s'", data.Node.ValueString()))
	}

	return strings.Join(parts, " ")
}

func buildEnterpriseTeamUpdateCommand(ctx context.Context, apiManager *api.ApiManager, plan *EnterpriseTeamResourceModel, state *EnterpriseTeamResourceModel) (string, error) {
	var parts []string

	parts = append(parts, "enterprise-team")

	// Required parameters
	if !state.Name.Equal(plan.Name) {
		parts = append(parts, fmt.Sprintf("--name '%s'", plan.Name.ValueString()))
	}

	// Optional parameters
	if !state.RestrictEdit.Equal(plan.RestrictEdit) {
		if !plan.RestrictEdit.IsNull() && plan.RestrictEdit.ValueBool() {
			parts = append(parts, "--restrict-edit on")
		} else {
			parts = append(parts, "--restrict-edit off")
		}
	}

	if !state.RestrictShare.Equal(plan.RestrictShare) {
		if !plan.RestrictShare.IsNull() && plan.RestrictShare.ValueBool() {
			parts = append(parts, "--restrict-share on")
		} else {
			parts = append(parts, "--restrict-share off")
		}
	}

	if !state.RestrictView.Equal(plan.RestrictView) {
		if !plan.RestrictView.IsNull() && plan.RestrictView.ValueBool() {
			parts = append(parts, "--restrict-view on")
		} else {
			parts = append(parts, "--restrict-view off")
		}
	}

	// TODO: we will node with its id all time
	if !state.Node.Equal(plan.Node) {
		parts = append(parts, fmt.Sprintf("--node '%s'", plan.Node.ValueString()))
	}

	// Process users and roles changes
	if !state.Users.Equal(plan.Users) {
		users, err := fetchAndProcessUsers(ctx, apiManager, state.Users, plan.Users)
		if err != nil {
			return "", err
		}
		if users != "" {
			parts = append(parts, users)
		}
	}

	if !state.Roles.Equal(plan.Roles) {
		roles, err := fetchAndProcessRoles(ctx, apiManager, state.Roles, plan.Roles)
		if err != nil {
			return "", err
		}
		if roles != "" {
			parts = append(parts, roles)
		}
	}

	parts = append(parts, fmt.Sprintf("'%s'", state.Id.ValueString()))

	return strings.Join(parts, " "), nil
}

// parseEnterpriseTeamReadResponse parses the JSON response from enterprise-info -t command
func parseEnterpriseTeamReadResponse(data interface{}) ([]EnterpriseTeamReadResponse, error) {
	var teams []EnterpriseTeamReadResponse

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("unable to process the response from Keeper Commander Service Mode API: %w", err)
	}

	if err := json.Unmarshal(dataBytes, &teams); err != nil {
		return nil, fmt.Errorf("unable to parse enterprise team information from Service Mode API response: %w", err)
	}

	return teams, nil
}

// parseRestrictsString parses the restricts string (e.g., "R W S") and returns boolean values
// R -> enable_privacy_screen (RestrictView)
// S -> restrict_record_re_share (RestrictShare)
// W -> restrict_record_edit (RestrictEdit)
func parseRestrictsString(restricts string) (restrictEdit, restrictShare, restrictView bool) {
	restricts = strings.TrimSpace(restricts)
	if restricts == "" {
		return false, false, false
	}

	// Split by space and check for each flag
	parts := strings.Fields(restricts)
	for _, part := range parts {
		switch part {
		case "W":
			restrictEdit = true
		case "S":
			restrictShare = true
		case "R":
			restrictView = true
		}
	}

	return restrictEdit, restrictShare, restrictView
}

// convertApiIdentifiersToOriginalFormat converts API response identifiers (names/emails)
// back to the original format from state (preserves IDs or names/emails as user provided)
func convertApiIdentifiersToOriginalFormat(
	ctx context.Context,
	apiManager *api.ApiManager,
	apiIdentifiers []string, // Names/emails from API response
	currentState types.Set, // Current state (what user originally provided)
	itemType string, // "role" or "user"
	fetchCommand string, // "enterprise-info -r --format json" or "enterprise-info -u --format json"
	parseFunc func(interface{}) (interface{}, error), // parseRolesResponse or parseUsersResponse
	buildLookupFunc func(interface{}) lookupMaps, // buildRoleLookupMaps or buildUserLookupMaps
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
		if identifier, isId := lookup.idToIdentifier[stateVal]; isId {
			// State has ID, map ID -> original ID
			stateIdToOriginal[stateVal] = stateVal
			// Also map the identifier (name/email) -> original ID
			stateIdentifierToOriginal[identifier] = stateVal
		} else if id, isIdentifier := lookup.identifierToId[stateVal]; isIdentifier {
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
		if _, isId := lookup.idToIdentifier[apiIdentifier]; isId {
			// API returned an ID (unlikely but handle it)
			id = apiIdentifier
		} else if foundId, isIdentifier := lookup.identifierToId[apiIdentifier]; isIdentifier {
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

// convertApiRolesToOriginalFormat converts role names from API back to original format from state
func convertApiRolesToOriginalFormat(ctx context.Context, apiManager *api.ApiManager, roleNames []string, currentState types.Set) (types.Set, error) {
	return convertApiIdentifiersToOriginalFormat(
		ctx,
		apiManager,
		roleNames,
		currentState,
		"role",
		"enterprise-info -r --format json",
		func(data interface{}) (interface{}, error) { return parseRolesResponse(data) },
		func(data interface{}) lookupMaps { return buildRoleLookupMaps(data.([]RoleInfo)) },
	)
}

// convertApiUsersToOriginalFormat converts user emails from API back to original format from state
func convertApiUsersToOriginalFormat(ctx context.Context, apiManager *api.ApiManager, userEmails []string, currentState types.Set) (types.Set, error) {
	return convertApiIdentifiersToOriginalFormat(
		ctx,
		apiManager,
		userEmails,
		currentState,
		"user",
		"enterprise-info -u --format json",
		func(data interface{}) (interface{}, error) { return parseUsersResponse(data) },
		func(data interface{}) lookupMaps { return buildUserLookupMaps(data.([]UserInfo)) },
	)
}

// convertNodeToName converts node from API response to node name
// API may return node as name or ID. If it's an ID, we fetch nodes and convert to name.
// Final state always stores node name (not ID).
func convertNodeToName(ctx context.Context, apiManager *api.ApiManager, nodeFromApi string) (types.String, error) {
	// Handle empty/null cases
	nodeFromApi = strings.TrimSpace(nodeFromApi)
	if nodeFromApi == "" {
		return types.StringNull(), nil
	}

	// Check if nodeFromApi is numeric (likely an ID)
	// If it's not numeric, assume it's already a name and use it directly
	nodeIdInt, err := strconv.Atoi(nodeFromApi)
	if err != nil {
		// Not numeric - assume it's a name, use it directly
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
				return types.StringValue(node.Name), nil
			}
		}
	}

	// Node ID not found - return null (node may have been deleted)
	// This handles edge case where node was deleted but team still references it
	return types.StringNull(), nil
}

// mapTeamReadResponseToModel maps the API response to the resource model
func mapTeamReadResponseToModel(ctx context.Context, apiManager *api.ApiManager, teamResp EnterpriseTeamReadResponse, state *EnterpriseTeamResourceModel) error {
	// Map ID
	state.Id = types.StringValue(teamResp.TeamUid)

	// Map Name
	state.Name = types.StringValue(teamResp.Name)

	// Parse and map restricts flags
	restrictEdit, restrictShare, restrictView := parseRestrictsString(teamResp.Restricts)
	state.RestrictEdit = types.BoolValue(restrictEdit)
	state.RestrictShare = types.BoolValue(restrictShare)
	state.RestrictView = types.BoolValue(restrictView)

	// Convert API response identifiers back to original format from state
	// Roles: preserve original format (name or ID) as user provided
	if len(teamResp.Roles) > 0 {
		rolesSet, err := convertApiRolesToOriginalFormat(ctx, apiManager, teamResp.Roles, state.Roles)
		if err != nil {
			return fmt.Errorf("failed to convert roles to original format: %w", err)
		}
		state.Roles = rolesSet
	} else {
		state.Roles = types.SetNull(types.StringType)
	}

	// Users: preserve original format (email or ID) as user provided
	if len(teamResp.Users) > 0 {
		usersSet, err := convertApiUsersToOriginalFormat(ctx, apiManager, teamResp.Users, state.Users)
		if err != nil {
			return fmt.Errorf("failed to convert users to original format: %w", err)
		}
		state.Users = usersSet
	} else {
		state.Users = types.SetNull(types.StringType)
	}

	// Node: convert to name (API always returns node name, but user may provide name or ID)
	// Always update from API response to detect external changes
	if teamResp.Node != "" {
		nodeName, err := convertNodeToName(ctx, apiManager, teamResp.Node)
		if err != nil {
			return fmt.Errorf("failed to convert node to name: %w", err)
		}
		// Update state with node name from API (always a name, not ID)
		state.Node = nodeName
	} else {
		// API returned empty node - set to null (team has no node assigned)
		state.Node = types.StringNull()
	}

	return nil
}
