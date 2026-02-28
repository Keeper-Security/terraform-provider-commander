// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriserole

import (
	"context"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *EnterpriseRoleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates and manages an enterprise role in your Keeper MSP or Enterprise account.<br><br>" +
			"Roles provide the organization the ability to define enforcements based on a user's job responsibility as well as provide delegated administrative functions.<br><br>" +
			"For more information, see https://docs.keeper.io/en/enterprise-guide/roles",
		MarkdownDescription: "Creates and manages an **enterprise role** in your Keeper MSP or Enterprise account.<br><br>" +
			"Roles provide the organization the ability to define enforcements based on a user's job responsibility as well as provide delegated administrative functions.<br><br>" +
			"For more information, see [Enterprise Roles documentation](https://docs.keeper.io/en/enterprise-guide/roles).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				Description: "Role ID assigned by Keeper to the role after it is created. " +
					"Use this value to import an existing role into Terraform state or to reference the role from other resources.",
				MarkdownDescription: "**Role ID** assigned by Keeper to the role after it is created. " +
					"Use this value to **import** an existing role into Terraform state or to reference the role from other resources.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "Set the display name for the enterprise role. Must be at least one character.",
				MarkdownDescription: "Set the **display name** for the enterprise role. Must be at least **one character**.",
				Validators: []validator.String{
					utils.StringMinLengthValidator("Enterprise Role Name", 1, false),
				},
			},
			"node": schema.StringAttribute{
				Required:            true,
				Description:         "The node that will manage this enterprise role. Provide the node name or node ID. ",
				MarkdownDescription: "The **node** that will manage this enterprise role. Provide the **node name** or **node ID**. ",
				Validators: []validator.String{
					utils.StringMinLengthValidator("Node", 1, true),
				},
			},
			"users": schema.SetAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				Description:         "Set of users assigned to this enterprise role. Provide user email addresses or user IDs. ",
				MarkdownDescription: "Set of **users** assigned to this enterprise role. Provide **user email addresses** or **user IDs**. ",
				Validators: []validator.Set{
					utils.SetNoEmptyStringsValidator("User"),
				},
			},
			"teams": schema.SetAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				Description:         "Set of teams assigned to this enterprise role. Provide team names or team IDs. ",
				MarkdownDescription: "Set of **teams** assigned to this enterprise role. Provide **team names** or **team IDs**. ",
				Validators: []validator.Set{
					utils.TeamsValidator,
				},
			},
			"managing_nodes": schema.MapNestedAttribute{
				Optional: true,
				Validators: []validator.Map{
					utils.MapKeysMinLengthValidator("managing node name", 1),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"privileges": schema.SetAttribute{
							Optional:    true,
							ElementType: types.StringType,
							Validators: []validator.Set{
								privilegesValidator{},
							},
							Description:         "Manage privileges to grant for this managing node. Valid values: manage_nodes, manage_user, manage_roles, manage_teams, run_reports, manage_bridge, approve_device, manage_record_types, run_compliance_reports, transfer_account, sharing_administrator, manage_companies",
							MarkdownDescription: "Manage **privileges** to grant for this managing node. Valid values: `manage_nodes`, `manage_user`, `manage_roles`, `manage_teams`, `run_reports`, `manage_bridge`, `approve_device`, `manage_record_types`, `run_compliance_reports`, `transfer_account`, `sharing_administrator`, `manage_companies`",
						},
						"cascade": schema.BoolAttribute{
							Optional:            true,
							Description:         "Manage extending admin-privileges for the specified role(s) to child nodes as well.",
							MarkdownDescription: "Manage extending **admin-privileges** for the **specified role(s)** to **child nodes** as well. ",
						},
					},
				},
				Description:         "Manage administrative permissions for the enterprise role. The map key is the node name/ID, and the value is an object with optional `privileges` and `cascade` fields.",
				MarkdownDescription: "Manage **administrative permissions** for the enterprise role. The map **key** is the **node name/ID** and the **value** is an object with optional `privileges` and `cascade` fields.",
			},
			"enforcement_policies": schema.MapAttribute{
				Optional: true,
				Validators: []validator.Map{
					enforcementPoliciesMapKeyValidator{},
				},
				PlanModifiers: []planmodifier.Map{
					enforcementPoliciesGPCPlanModifier{},
				},
				ElementType:         types.StringType,
				Description:         enforcementPoliciesDescription(),
				MarkdownDescription: enforcementPoliciesMarkdownDescription(),
			},
			"managed_company": schema.StringAttribute{
				Optional:            true,
				Description:         utils.EnterpriseManagedCompanySchemaAttributeDescription,
				MarkdownDescription: utils.EnterpriseManagedCompanySchemaAttributeMarkdownDescription,
				Validators: []validator.String{
					utils.ManagedCompanyValidator,
				},
			},
		},
	}
}

// enforcementPoliciesDescription returns the plain-text description for enforcement_policies.
func enforcementPoliciesDescription() string {
	return "Enforcement policies applied to users in this role. Map key = policy name (see valid keys below), value = string (format depends on policy type).\n\n" +
		"General: All values must be strings (e.g. \"true\" for bool, \"12\" for number).\n\n" +
		"Value types:\n" +
		"• Boolean: \"true\" or \"false\". Example: REQUIRE_TWO_FACTOR = \"true\".\n" +
		"• Long (integer): string number for days, minutes, etc. Example: MASTER_PASSWORD_MINIMUM_LENGTH = \"12\", LOGOUT_TIMER_WEB = \"30\".\n" +
		"• String: arbitrary string. Example: RESTRICT_DOMAIN_ACCESS = \"example.com\".\n" +
		"• ACCOUNT_SHARE (REQUIRE_ACCOUNT_SHARE): role name or role ID. Example: \"Compliance_Admins\" or \"12345\".\n" +
		"• IP_WHITELIST (RESTRICT_IP_ADDRESSES, RESTRICT_VAULT_IP_ADDRESSES, TIP_ZONE_RESTRICT_ALLOWED_IP_RANGES): comma-separated IPs; use CIDR (e.g. 192.168.1.0/24) or range (10.0.0.1-10.0.0.50). Example: \"192.168.1.0/24,10.0.0.1-10.0.0.50\".\n" +
		"• TWO_FACTOR_DURATION (TWO_FACTOR_DURATION_WEB, _MOBILE, _DESKTOP): one of \"login\", \"12_hours\", \"24_hours\", \"30_days\", \"forever\".\n" +
		"• RECORD_TYPES: comma-separated record type IDs or \"all\". Example: \"login,password,bank_account\".\n" +
		"• TERNARY (KEEPER_FILL_* policies): \"enforce\", \"disable\", or \"null\".\n" +
		"• GENERATED_PASSWORD_COMPLEXITY: JSON string; use jsonencode([{ domains, length, lower-use, upper-use, ... }]).\n" +
		"• GENERATED_SECURITY_QUESTION_COMPLEXITY: string per API.\n\n" +
		"Valid policy keys: " + strings.Join(ValidEnforcementPolicyKeys, ", ") + ".\n\n" +
		"For more information, see https://docs.keeper.io/en/keeperpam/commander-cli/command-reference/enterprise-management-commands#enforcement-policies."
}

// enforcementPolicyValueType returns the value type label for an enforcement policy key (for docs table).
func enforcementPolicyValueType(key string) string {
	if TwoFactorDurationPolicyKeys[key] {
		return "`\"login\"`, `\"12_hours\"`, `\"24_hours\"`, `\"30_days\"`, `\"forever\"`"
	}
	if KeeperFillPolicyKeys[key] {
		return "`\"enforce\"`, `\"disable\"`, `\"null\"`"
	}
	switch key {
	case RequireAccountShare:
		return "Role name or role ID. That role must be an admin role with TRANSFER_ACCOUNT privilege.<br>eg: `\"Compliance_Admins\"` or `\"12345\"`"
	case RestrictIpAddresses, RestrictVaultIpAddresses, TipZoneRestrictAllowedIpRanges:
		return "IP_WHITELIST (comma-separated CIDR or range).<br>eg: `\"192.168.1.0/24,10.0.0.1-10.0.0.50\"`"
	case RestrictRecordTypes:
		return "RECORD_TYPES (comma-separated or all).<br>eg: `\"login,password,bank_account\"`"
	case GeneratedPasswordComplexity:
		return "JSON string as jsonencode"
	case GeneratedSecurityQuestionComplexity, RestrictDomainAccess, RestrictDomainCreate:
		return "String"
	case MasterPasswordMinimumLength, MasterPasswordMinimumSpecial, MasterPasswordMinimumUpper,
		MasterPasswordMinimumLower, MasterPasswordMinimumDigits, MasterPasswordRestrictDaysBeforeReuse,
		MasterPasswordMaximumDaysBeforeChange, MasterPasswordExpiredAsOf, MinimumPbkdf2Iterations,
		MaxSessionLoginTime, AutomaticBackupEveryXDays, LogoutTimerWeb, LogoutTimerMobile, LogoutTimerDesktop,
		DaysBeforeDeletedRecordsClearedPerm, DaysBeforeDeletedRecordsAutoCleared, ResendEnterpriseInviteInXDays,
		MaximumRecordSize, RestrictClipboardExpireInXSecs:
		return "String(Integer)"
	default:
		return "String(Boolean)"
	}
}

// enforcementPoliciesMarkdownDescription returns the Markdown description for enforcement_policies.
func enforcementPoliciesMarkdownDescription() string {
	var b strings.Builder
	b.WriteString("**Enforcement policies** applied to users in this role. Map **key** = policy name, **value** = string (format depends on policy type).\n")
	b.WriteString("**General:** All values must be **strings** (e.g. `\"true\"` for bool, `\"12\"` for number).\n")
	b.WriteString("| Enforcement Policy Key | Value Type |\n")
	b.WriteString("|------------------------|------------|\n")
	for _, key := range ValidEnforcementPolicyKeys {
		valType := enforcementPolicyValueType(key)
		b.WriteString("| `" + key + "` | " + valType + " |<br>\n")
	}
	b.WriteString("\nFor more information, see [Enforcement Policies](https://docs.keeper.io/en/keeperpam/commander-cli/command-reference/enterprise-management-commands#enforcement-policies).")
	return b.String()
}
