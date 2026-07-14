// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamuser

const (
	SchemaDescription         = "Creates and manages a Keeper PAM User record in the vault."
	SchemaMarkdownDescription = "Creates and manages a Keeper **PAM User** record (`pamUser`) in the vault.\n\nA PAM User record stores privileged credentials (login/password) that can be associated with PAM Machines for rotation, connections, and tunneling."

	IDDescription                           = "The unique identifier (UID) of the PAM User record."
	IDMarkdownDescription                   = "The unique identifier (UID) of the PAM User record."
	TitleDescription                        = "Title of the PAM User record."
	TitleMarkdownDescription                = "Title of the PAM User record."
	LoginDescription                        = "Login (username) for the PAM User."
	LoginMarkdownDescription                = "Login (username) for the PAM User."
	PasswordDescription                     = "Password for the PAM User."
	PasswordMarkdownDescription             = "Password for the PAM User."
	FolderDescription                       = "Folder path or UID where the record will be stored."
	FolderMarkdownDescription               = "Folder path or UID where the record will be stored."
	NotesDescription                        = "Optional notes for the PAM User record."
	NotesMarkdownDescription                = "Optional notes for the PAM User record."
	DistinguishedNameDescription            = "LDAP distinguished name of the PAM User (e.g. CN=svc_myapp,OU=Service Accounts,DC=corp,DC=local)."
	DistinguishedNameMarkdownDescription    = "LDAP distinguished name of the PAM User (e.g. `CN=svc_myapp,OU=Service Accounts,DC=corp,DC=local`)."
	PrivatePEMKeyDescription                = "Private PEM key associated with the PAM User."
	PrivatePEMKeyMarkdownDescription        = "Private PEM key associated with the PAM User."
	PublicKeyDescription                    = "Public key associated with the PAM User."
	PublicKeyMarkdownDescription            = "Public key associated with the PAM User."
	PrivateKeyPassphraseDescription         = "Passphrase for the private key associated with the PAM User."
	PrivateKeyPassphraseMarkdownDescription = "Passphrase for the private key associated with the PAM User."
	ConnectDatabaseDescription              = "Database name the PAM User connects to."
	ConnectDatabaseMarkdownDescription      = "Database name the PAM User connects to."
	ManagedDescription                      = "Whether this PAM User account is managed by Keeper."
	ManagedMarkdownDescription              = "Whether this PAM User account is managed by Keeper."

	RotationSettingsDescription         = "Rotation settings for the PAM User record."
	RotationSettingsMarkdownDescription = "Rotation settings for the PAM User record.\n\n**Required:** `rotation_profile`. **Profile-specific:** `general` requires `configuration` and `resource` (do not set `saas_config`); `iam_user` and `scripts_only` require `configuration` (do not set `resource` or `saas_config`); `saas` requires `configuration` and `saas_config` (do not set `resource`).\n\n**Schedule:** when `enabled` is not `false`, set **exactly one** of `on_demand`, `use_default_rotation_schedule`, `schedule_cron`, or `schedule_json` (required and mutually exclusive). When `enabled` is `false`, do not set schedule fields or `complexity`."

	RotStatusNoRotation = "RRS_NO_ROTATION"

	RotProfileGeneral     = "general"
	RotProfileIAMUser     = "iam_user"
	RotProfileScriptsOnly = "scripts_only"
	RotProfileSaaS        = "saas"

	RotProfileScheduleTypeManual    = "manual"
	RotProfileScheduleTypeScheduled = "scheduled"

	RotProfileDescription              = "Rotation profile type: general (resource-based), iam_user (IAM/Azure user), scripts_only (run PAM scripts only), or saas (SaaS Account). Required when rotation_settings is set."
	RotProfileMarkdownDescription      = "Rotation profile type: `general` (resource-based), `iam_user` (IAM/Azure user), `scripts_only` (run PAM scripts only), or `saas` (SaaS Account). **Required** when `rotation_settings` is set."
	RotConfigDescription               = "PAM Configuration UID to use for rotation. Required when rotation_profile is general, iam_user, saas or scripts_only."
	RotConfigMarkdownDescription       = "PAM Configuration UID to use for rotation. **Required** when `rotation_profile` is `general`, `iam_user`, `saas` or `scripts_only`."
	RotIamAadConfigDescription         = "PAM Configuration UID for IAM or Azure AD users. Used instead of resource when rotation_profile is iam_user."
	RotIamAadConfigMarkdownDescription = "PAM Configuration UID for IAM or Azure AD users. Used instead of `resource` when `rotation_profile` is `iam_user`."
	RotSaaSConfigDescription           = "SaaS Configuration UID associated with the PAM Configuration to use for rotation. Required when rotation_profile is saas (along with configuration)."
	RotSaaSConfigMarkdownDescription   = "SaaS Configuration UID associated with the PAM Configuration to use for rotation. **Required** when `rotation_profile` is `saas` (along with `configuration`)."
	RotResourceDescription             = "UID of the PAM resource record (machine or database) this user rotates on. Required when rotation_profile is general (along with configuration)."
	RotResourceMarkdownDescription     = "UID of the PAM resource record (machine or database) this user rotates on. **Required** when `rotation_profile` is `general` (along with `configuration`)."
	RotAdminUserDescription            = "UID of the PAM User record to use as admin credential when rotating."
	RotAdminUserMarkdownDescription    = "UID of the PAM User record to use as admin credential when rotating."
	RotEnabledDescription              = "Whether rotation is enabled for this PAM User. When false, schedule and complexity fields must not be set."
	RotEnabledMarkdownDescription      = "Whether rotation is enabled for this PAM User. When `false`, do not set `on_demand`, `use_default_rotation_schedule`, `schedule_cron`, `schedule_json`, or `complexity`."
	RotScheduleCronDescription         = "Cron schedule for rotation using Keeper Quartz format (6 or 7 fields), e.g. \"0 28 17 ? * *\". Schedules must have at least a 1-hour interval. Invalid expressions are rejected at plan time."
	RotScheduleCronMarkdownDescription = "Cron schedule for rotation using the [Keeper Quartz cron spec](https://docs.keeper.io/keeperpam/privileged-access-manager/references/cron-spec) (**6 or 7 fields**, seconds first, e.g. `0 28 17 ? * *`). Schedules must have at least a **1-hour interval** between executions. Invalid expressions fail at **plan** time so the vault record is not created first."
	RotScheduleJSONDescription         = "Schedule JSON for rotation. Supported types: DAILY, WEEKLY, MONTHLY_BY_WEEKDAY, YEARLY. Use either time (HH:MM:SS) or utcTime (HH:MM), not both. For cron schedules use schedule_cron instead. Examples: {\"type\":\"DAILY\",\"intervalCount\":1,\"time\":\"17:00:00\",\"tz\":\"Asia/Calcutta\"}; {\"type\":\"WEEKLY\",\"time\":\"17:00:00\",\"tz\":\"Asia/Calcutta\",\"weekday\":\"WEDNESDAY\"}."
	RotScheduleJSONMarkdownDescription = "Schedule JSON for rotation. Provide a **single JSON object** with a required `type` field.\n\n" +
		"**Supported types:** `DAILY`, `WEEKLY`, `MONTHLY_BY_WEEKDAY`, `YEARLY`.\n\n" +
		"**Common fields:**\n" +
		"- `time` — `HH:MM:SS` (24-hour), **or** `utcTime` — `HH:MM` (use one, not both)\n" +
		"- `tz` — IANA timezone (recommended), e.g. `Asia/Calcutta`, `Etc/UTC`\n" +
		"- `intervalCount` — optional positive integer (default `1`)\n\n" +
		"**Type-specific fields:**\n" +
		"- `WEEKLY` / `MONTHLY_BY_WEEKDAY` — `weekday` (`SUNDAY`..`SATURDAY`)\n" +
		"- `YEARLY` — `monthDay` (`1`–`28`)\n" +
		"- `MONTHLY_BY_WEEKDAY` — `occurrence` (`FIRST`, `SECOND`, `THIRD`, `FOURTH`, `LAST`)\n" +
		"- `YEARLY` — `month` (`JANUARY`..`DECEMBER`)\n\n" +
		"**Examples by type:**\n\n" +
		"`DAILY` — every day at 5:00 PM IST:\n" +
		"```json\n" +
		"{\"type\":\"DAILY\",\"intervalCount\":1,\"time\":\"17:00:00\",\"tz\":\"Asia/Calcutta\"}\n" +
		"```\n\n" +
		"`WEEKLY` — every Wednesday at 5:00 PM IST:\n" +
		"```json\n" +
		"{\"type\":\"WEEKLY\",\"intervalCount\":1,\"time\":\"17:00:00\",\"tz\":\"Asia/Calcutta\",\"weekday\":\"WEDNESDAY\"}\n" +
		"```\n\n" +
		"`WEEKLY` — every Saturday at midnight UTC using `utcTime`:\n" +
		"```json\n" +
		"{\"type\":\"WEEKLY\",\"utcTime\":\"00:00\",\"weekday\":\"SATURDAY\",\"intervalCount\":1,\"tz\":\"Etc/UTC\"}\n" +
		"```\n\n" +
		"`MONTHLY_BY_WEEKDAY` — on the second Tuesday of each month at 9:30 AM Eastern:\n" +
		"```json\n" +
		"{\"type\":\"MONTHLY_BY_WEEKDAY\",\"intervalCount\":1,\"time\":\"09:30:00\",\"tz\":\"America/New_York\",\"weekday\":\"TUESDAY\",\"occurrence\":\"SECOND\"}\n" +
		"```\n\n" +
		"`YEARLY` — every May 20 at midnight UTC:\n" +
		"```json\n" +
		"{\"type\":\"YEARLY\",\"intervalCount\":1,\"time\":\"00:00:00\",\"tz\":\"Etc/UTC\",\"month\":\"MAY\",\"monthDay\":20}\n" +
		"```\n\n" +
		"**Terraform usage** (recommended — avoids escaping issues):\n" +
		"```hcl\n" +
		"rotation_settings {\n" +
		"  schedule_json = jsonencode({\n" +
		"    type          = \"WEEKLY\"\n" +
		"    intervalCount = 1\n" +
		"    time          = \"17:00:00\"\n" +
		"    tz            = \"Asia/Calcutta\"\n" +
		"    weekday       = \"WEDNESDAY\"\n" +
		"  })\n" +
		"}\n" +
		"```\n\n" +
		"Mutually exclusive with `on_demand`, `schedule_cron`, and `use_default_rotation_schedule`. For cron schedules use `schedule_cron`."
	RotOnDemandDescription                           = "If true, rotation is on-demand (manual) only."
	RotOnDemandMarkdownDescription                   = "If `true`, rotation is on-demand (manual) only."
	RotUseDefaultRotationScheduleDescription         = "If true, uses the schedule from the PAM Configuration."
	RotUseDefaultRotationScheduleMarkdownDescription = "If `true`, uses the schedule from the PAM Configuration instead of a per-record schedule."
	RotComplexityDescription                         = "Password complexity for rotation: five integers length,upper,lower,digits,symbols (length 20\u201399; each count 0\u201399 per Keeper UI). Must not be set when enabled is false."
	RotComplexityMarkdownDescription                 = "Password complexity for rotation: `length,upper,lower,digits,symbols` as **five integers**. Password **length** must be **20\u201399**; upper, lower, digits, and symbols minimums must each be **0\u201399** (Keeper UI limits). Must not be set when `enabled` is `false`. Invalid values fail at plan time."

	ErrSummaryCreateFailed       = "PAM User Record Create Failed"
	ErrSummaryReadFailed         = "PAM User Record Read Failed"
	ErrSummaryUpdateFailed       = "PAM User Record Update Failed"
	ErrSummaryRotationEditFailed = "PAM User Rotation Edit Failed"

	ErrDetailCreateFailed       = "Something went wrong when creating the PAM User record."
	ErrDetailReadFailed         = "Something went wrong when reading the PAM User record."
	ErrDetailUpdateFailed       = "Something went wrong when updating the PAM User record."
	ErrDetailRotationEditFailed = "Something went wrong when configuring rotation for the PAM User record."
	ErrDetailRotationInfoFailed = "Something went wrong when reading rotation info for the PAM User record."

	// Commander field dot-notation keys for pamUser record-add / record-update.
	FieldLogin                = "f.login"
	FieldPassword             = "f.password"
	FieldDistinguishedName    = "f.text.distinguishedName"
	FieldPrivatePEMKey        = "f.secret.privatePEMKey"
	FieldPublicKey            = "f.secret.publicKey"
	FieldPrivateKeyPassphrase = "f.secret.privateKeyPassphrase"
	FieldConnectDatabase      = "f.text.connectDatabase"
	FieldManaged              = "f.checkbox.managed"

	// Commander pam rotation edit command and flags.
	CmdPamRotationEdit  = "pam rotation edit"
	CmdPamRotationInfo  = "pam rotation info"
	FlagConfig          = "--config"
	FlagIamAadConfig    = "--iam-aad-config"
	FlagSaaSConfig      = "--saas-config-uid"
	FlagRotationProfile = "--rotation-profile"
	FlagResource        = "--resource"
	FlagAdminUser       = "--admin-user"
	FlagEnable          = "--enable"
	FlagDisable         = "--disable"
	FlagForce           = "--force"
	FlagOnDemand        = "--on-demand"
	FlagScheduleCron    = "--schedulecron"
	FlagScheduleJSON    = "--schedulejson"
	FlagScheduleConfig  = "--schedule-config"
	FlagComplexity      = "--complexity"
	FlagRecordShort     = "-r"
)
