// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package utils_test

import (
	"context"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestParseNodesResponse(t *testing.T) {
	data := []interface{}{
		map[string]interface{}{"node_id": float64(1), "name": "Node1", "parent_node": "Root", "parent_id": float64(0)},
	}
	nodes, err := utils.ParseNodesResponse(data)
	if err != nil {
		t.Fatalf("ParseNodesResponse: %v", err)
	}
	if len(nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(nodes))
	}
	if nodes[0].NodeId != 1 || nodes[0].Name != "Node1" {
		t.Errorf("got %+v", nodes[0])
	}
}

func TestParseNodesResponse_Empty(t *testing.T) {
	nodes, err := utils.ParseNodesResponse([]interface{}{})
	if err != nil {
		t.Fatalf("ParseNodesResponse: %v", err)
	}
	if len(nodes) != 0 {
		t.Errorf("expected 0 nodes, got %d", len(nodes))
	}
}

func TestBuildNodeLookupMaps(t *testing.T) {
	nodes := []utils.EnterpriseNodeResponse{
		{NodeId: 1, Name: "Root"},
		{NodeId: 2, Name: "Parent\\Child"},
	}
	lookup := utils.BuildNodeLookupMaps(nodes)
	if lookup.IdentifierToId["Root"] != "1" {
		t.Errorf("IdentifierToId[Root] = %s; want 1", lookup.IdentifierToId["Root"])
	}
	if lookup.IdentifierToId["Child"] != "2" {
		t.Errorf("IdentifierToId[Child] = %s; want 2", lookup.IdentifierToId["Child"])
	}
	if lookup.IdToIdentifier["1"] != "Root" {
		t.Errorf("IdToIdentifier[1] = %s; want Root", lookup.IdToIdentifier["1"])
	}
	if lookup.IdToIdentifier["2"] != "Child" {
		t.Errorf("IdToIdentifier[2] = %s; want Child", lookup.IdToIdentifier["2"])
	}
}

func TestBuildNodeLookupMaps_SkipsZeroId(t *testing.T) {
	nodes := []utils.EnterpriseNodeResponse{
		{NodeId: 0, Name: "Bad"},
	}
	lookup := utils.BuildNodeLookupMaps(nodes)
	if len(lookup.IdentifierToId) != 0 {
		t.Errorf("expected no entries for zero node_id, got %v", lookup.IdentifierToId)
	}
}

func TestParseRolesResponse(t *testing.T) {
	data := []interface{}{
		map[string]interface{}{"role_id": float64(1), "name": "Admin"},
	}
	roles, err := utils.ParseRolesResponse(data)
	if err != nil {
		t.Fatalf("ParseRolesResponse: %v", err)
	}
	if len(roles) != 1 {
		t.Fatalf("expected 1 role, got %d", len(roles))
	}
	if roles[0].RoleId != 1 || roles[0].Name != "Admin" {
		t.Errorf("got %+v", roles[0])
	}
}

func TestBuildRoleLookupMaps(t *testing.T) {
	roles := []utils.EnterpriseRoleResponse{
		{RoleId: 1, Name: "Admin"},
		{RoleId: 2, Name: "User"},
	}
	lookup := utils.BuildRoleLookupMaps(roles)
	if lookup.IdentifierToId["Admin"] != "1" {
		t.Errorf("IdentifierToId[Admin] = %s; want 1", lookup.IdentifierToId["Admin"])
	}
	if lookup.IdToIdentifier["1"] != "Admin" {
		t.Errorf("IdToIdentifier[1] = %s; want Admin", lookup.IdToIdentifier["1"])
	}
}

func TestBuildRoleLookupMaps_SkipsZeroId(t *testing.T) {
	roles := []utils.EnterpriseRoleResponse{
		{RoleId: 0, Name: "Bad"},
	}
	lookup := utils.BuildRoleLookupMaps(roles)
	if len(lookup.IdentifierToId) != 0 {
		t.Errorf("expected no entries for zero role_id, got %v", lookup.IdentifierToId)
	}
}

func TestParseTeamsResponse(t *testing.T) {
	data := []interface{}{
		map[string]interface{}{"team_uid": "uid1", "name": "Team1"},
	}
	teams, err := utils.ParseTeamsResponse(data)
	if err != nil {
		t.Fatalf("ParseTeamsResponse: %v", err)
	}
	if len(teams) != 1 {
		t.Fatalf("expected 1 team, got %d", len(teams))
	}
	if teams[0].TeamUid != "uid1" || teams[0].Name != "Team1" {
		t.Errorf("got %+v", teams[0])
	}
}

func TestBuildTeamLookupMaps(t *testing.T) {
	teams := []utils.EnterpriseTeamResponse{
		{TeamUid: "uid1", Name: "Team1"},
		{TeamUid: "uid2", Name: "Team2"},
	}
	lookup := utils.BuildTeamLookupMaps(teams)
	if lookup.IdentifierToId["Team1"] != "uid1" {
		t.Errorf("IdentifierToId[Team1] = %s; want uid1", lookup.IdentifierToId["Team1"])
	}
	if lookup.IdToIdentifier["uid1"] != "Team1" {
		t.Errorf("IdToIdentifier[uid1] = %s; want Team1", lookup.IdToIdentifier["uid1"])
	}
}

func TestParseUsersResponse(t *testing.T) {
	data := []interface{}{
		map[string]interface{}{"user_id": float64(1), "email": "user@example.com"},
	}
	users, err := utils.ParseUsersResponse(data)
	if err != nil {
		t.Fatalf("ParseUsersResponse: %v", err)
	}
	if len(users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(users))
	}
	if users[0].UserId != 1 || users[0].Email != "user@example.com" {
		t.Errorf("got %+v", users[0])
	}
}

func TestBuildUserLookupMaps(t *testing.T) {
	users := []utils.EnterpriseUserResponse{
		{UserId: 1, Email: "a@example.com"},
		{UserId: 2, Email: "b@example.com"},
	}
	lookup := utils.BuildUserLookupMaps(users)
	if lookup.IdentifierToId["a@example.com"] != "1" {
		t.Errorf("IdentifierToId[a@example.com] = %s; want 1", lookup.IdentifierToId["a@example.com"])
	}
	if lookup.IdToIdentifier["1"] != "a@example.com" {
		t.Errorf("IdToIdentifier[1] = %s; want a@example.com", lookup.IdToIdentifier["1"])
	}
}

func TestConvertRolesToIdMap(t *testing.T) {
	ctx := context.Background()
	setVal, _ := types.SetValueFrom(ctx, types.StringType, []types.String{types.StringValue("Admin")})
	lookup := utils.LookupMaps{
		IdentifierToId: map[string]string{"Admin": "1"},
		IdToIdentifier: map[string]string{"1": "Admin"},
	}
	roles := []utils.EnterpriseRoleResponse{{RoleId: 1, Name: "Admin"}}
	result, err := utils.ConvertRolesToIdMap(setVal, lookup, roles)
	if err != nil {
		t.Fatalf("ConvertRolesToIdMap: %v", err)
	}
	if len(result) != 1 || result["1"] != "Admin" {
		t.Errorf("got %v", result)
	}
}

func TestConvertTeamsToIdMap(t *testing.T) {
	ctx := context.Background()
	setVal, _ := types.SetValueFrom(ctx, types.StringType, []types.String{types.StringValue("Team1")})
	lookup := utils.LookupMaps{
		IdentifierToId: map[string]string{"Team1": "uid1"},
		IdToIdentifier: map[string]string{"uid1": "Team1"},
	}
	teams := []utils.EnterpriseTeamResponse{{TeamUid: "uid1", Name: "Team1"}}
	result, err := utils.ConvertTeamsToIdMap(setVal, lookup, teams)
	if err != nil {
		t.Fatalf("ConvertTeamsToIdMap: %v", err)
	}
	if len(result) != 1 || result["uid1"] != "Team1" {
		t.Errorf("got %v", result)
	}
}

func TestConvertUsersToIdMap(t *testing.T) {
	ctx := context.Background()
	setVal, _ := types.SetValueFrom(ctx, types.StringType, []types.String{types.StringValue("user@example.com")})
	lookup := utils.LookupMaps{
		IdentifierToId: map[string]string{"user@example.com": "1"},
		IdToIdentifier: map[string]string{"1": "user@example.com"},
	}
	users := []utils.EnterpriseUserResponse{{UserId: 1, Email: "user@example.com"}}
	result, err := utils.ConvertUsersToIdMap(setVal, lookup, users)
	if err != nil {
		t.Fatalf("ConvertUsersToIdMap: %v", err)
	}
	if len(result) != 1 || result["1"] != "user@example.com" {
		t.Errorf("got %v", result)
	}
}
