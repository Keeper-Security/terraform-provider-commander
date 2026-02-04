package enterpriseteam

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ImportState supports import ID formats:
//   - Team only: team name or team UID (e.g. "Engineering" or 1234567890123456)
//   - With managed company: "managed_company_name_or_id,team_name_or_uid" (comma-separated, e.g. "Test Company,Engineering" or "1169425105420462,1234567890123456")
//
// After import, Terraform runs Read to refresh state from the API.
func (r *EnterpriseTeamResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if err := r.ensureApiManager(); err != nil {
		resp.Diagnostics.AddError(
			"Provider Configuration Error",
			err.Error(),
		)
		return
	}

	importID := strings.TrimSpace(req.ID)
	if importID == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			"Import ID cannot be empty. Use: (1) team name or team UID alone, e.g. Engineering or 1234567890123456; or (2) for a team in a managed company, use \"managed_company_name_or_id,team_name_or_uid\" (comma-separated), e.g. \"Test Company,Engineering\" or \"1169425105420462,1234567890123456\".",
		)
		return
	}

	var resourceIdentifier, managedCompany string
	if parts := strings.SplitN(importID, ",", 2); len(parts) == 2 {
		managedCompany = strings.TrimSpace(parts[0])
		resourceIdentifier = strings.TrimSpace(parts[1])
	} else {
		resourceIdentifier = importID
	}

	if resourceIdentifier == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			"When using managed company format \"managed_company_name_or_id,team\", the team part cannot be empty. Examples: \"Test Company,Engineering\" or \"1169425105420462,1234567890123456\".",
		)
		return
	}

	state := EnterpriseTeamResourceModel{
		Id:             types.StringValue(resourceIdentifier),
		Name:           types.StringNull(),
		RestrictEdit:   types.BoolNull(),
		RestrictShare:  types.BoolNull(),
		RestrictView:   types.BoolNull(),
		Users:          types.SetNull(types.StringType),
		Roles:          types.SetNull(types.StringType),
		Node:           types.StringNull(),
		ManagedCompany: types.StringNull(),
	}
	if managedCompany != "" {
		state.ManagedCompany = types.StringValue(managedCompany)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
