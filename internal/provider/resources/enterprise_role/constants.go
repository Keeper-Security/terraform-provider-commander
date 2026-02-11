// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriserole

// Enforcement policy key constants.
const (
	MasterPasswordMinimumLength              = "MASTER_PASSWORD_MINIMUM_LENGTH"
	MasterPasswordMinimumSpecial             = "MASTER_PASSWORD_MINIMUM_SPECIAL"
	MasterPasswordMinimumUpper               = "MASTER_PASSWORD_MINIMUM_UPPER"
	MasterPasswordMinimumLower               = "MASTER_PASSWORD_MINIMUM_LOWER"
	MasterPasswordMinimumDigits              = "MASTER_PASSWORD_MINIMUM_DIGITS"
	MasterPasswordRestrictDaysBeforeReuse    = "MASTER_PASSWORD_RESTRICT_DAYS_BEFORE_REUSE"
	RequireTwoFactor                         = "REQUIRE_TWO_FACTOR"
	MasterPasswordMaximumDaysBeforeChange    = "MASTER_PASSWORD_MAXIMUM_DAYS_BEFORE_CHANGE"
	MasterPasswordExpiredAsOf                = "MASTER_PASSWORD_EXPIRED_AS_OF"
	MasterPasswordReentry                    = "MASTER_PASSWORD_REENTRY"
	MinimumPbkdf2Iterations                  = "MINIMUM_PBKDF2_ITERATIONS"
	MaxSessionLoginTime                      = "MAX_SESSION_LOGIN_TIME"
	RestrictPersistentLogin                  = "RESTRICT_PERSISTENT_LOGIN"
	StayLoggedInDefault                      = "STAY_LOGGED_IN_DEFAULT"
	RestrictOfflineAccess                    = "RESTRICT_OFFLINE_ACCESS"
	RestrictSharingAll                       = "RESTRICT_SHARING_ALL"
	RestrictSharingEnterprise                = "RESTRICT_SHARING_ENTERPRISE"
	RestrictSharingAllOutgoing               = "RESTRICT_SHARING_ALL_OUTGOING"
	RestrictSharingEnterpriseOutgoing        = "RESTRICT_SHARING_ENTERPRISE_OUTGOING"
	RestrictSharingAllIncoming               = "RESTRICT_SHARING_ALL_INCOMING"
	RestrictSharingEnterpriseIncoming        = "RESTRICT_SHARING_ENTERPRISE_INCOMING"
	RestrictSharingRecordWithAttachments     = "RESTRICT_SHARING_RECORD_WITH_ATTACHMENTS"
	RestrictSharingOutsideOfIsolatedNodes    = "RESTRICT_SHARING_OUTSIDE_OF_ISOLATED_NODES"
	RestrictSharingRecordToSharedFolders     = "RESTRICT_SHARING_RECORD_TO_SHARED_FOLDERS"
	RequireAccountShare                      = "REQUIRE_ACCOUNT_SHARE"
	RestrictExport                           = "RESTRICT_EXPORT"
	RestrictImport                           = "RESTRICT_IMPORT"
	RestrictImportSharedFolders              = "RESTRICT_IMPORT_SHARED_FOLDERS"
	RestrictFileUpload                       = "RESTRICT_FILE_UPLOAD"
	RestrictIpAddresses                      = "RESTRICT_IP_ADDRESSES"
	RestrictVaultIpAddresses                 = "RESTRICT_VAULT_IP_ADDRESSES"
	TipZoneRestrictAllowedIpRanges           = "TIP_ZONE_RESTRICT_ALLOWED_IP_RANGES"
	RestrictIpAutoapproval                   = "RESTRICT_IP_AUTOAPPROVAL"
	RequireDeviceApproval                    = "REQUIRE_DEVICE_APPROVAL"
	RequireAccountRecoveryApproval           = "REQUIRE_ACCOUNT_RECOVERY_APPROVAL"
	RestrictIosFingerprint                   = "RESTRICT_IOS_FINGERPRINT"
	RestrictMacFingerprint                   = "RESTRICT_MAC_FINGERPRINT"
	RestrictAndroidFingerprint               = "RESTRICT_ANDROID_FINGERPRINT"
	RestrictWindowsFingerprint               = "RESTRICT_WINDOWS_FINGERPRINT"
	LogoutTimerWeb                           = "LOGOUT_TIMER_WEB"
	LogoutTimerMobile                        = "LOGOUT_TIMER_MOBILE"
	LogoutTimerDesktop                       = "LOGOUT_TIMER_DESKTOP"
	RestrictWebVaultAccess                   = "RESTRICT_WEB_VAULT_ACCESS"
	RestrictExtensionsAccess                 = "RESTRICT_EXTENSIONS_ACCESS"
	RestrictMobileAccess                     = "RESTRICT_MOBILE_ACCESS"
	RestrictDesktopAccess                    = "RESTRICT_DESKTOP_ACCESS"
	RestrictMobileIosAccess                  = "RESTRICT_MOBILE_IOS_ACCESS"
	RestrictMobileAndroidAccess              = "RESTRICT_MOBILE_ANDROID_ACCESS"
	RestrictMobileWindowsPhoneAccess         = "RESTRICT_MOBILE_WINDOWS_PHONE_ACCESS"
	RestrictDesktopWinAccess                 = "RESTRICT_DESKTOP_WIN_ACCESS"
	RestrictDesktopMacAccess                 = "RESTRICT_DESKTOP_MAC_ACCESS"
	RestrictChatDesktopAccess                = "RESTRICT_CHAT_DESKTOP_ACCESS"
	RestrictChatMobileAccess                 = "RESTRICT_CHAT_MOBILE_ACCESS"
	RestrictCommanderAccess                  = "RESTRICT_COMMANDER_ACCESS"
	RestrictTwoFactorChannelText             = "RESTRICT_TWO_FACTOR_CHANNEL_TEXT"
	RestrictTwoFactorChannelGoogle           = "RESTRICT_TWO_FACTOR_CHANNEL_GOOGLE"
	RestrictTwoFactorChannelDna              = "RESTRICT_TWO_FACTOR_CHANNEL_DNA"
	RestrictTwoFactorChannelDuo              = "RESTRICT_TWO_FACTOR_CHANNEL_DUO"
	RestrictTwoFactorChannelRsa              = "RESTRICT_TWO_FACTOR_CHANNEL_RSA"
	RestrictTwoFactorChannelSecurityKeys     = "RESTRICT_TWO_FACTOR_CHANNEL_SECURITY_KEYS"
	TwoFactorDurationWeb                     = "TWO_FACTOR_DURATION_WEB"
	TwoFactorDurationMobile                  = "TWO_FACTOR_DURATION_MOBILE"
	TwoFactorDurationDesktop                 = "TWO_FACTOR_DURATION_DESKTOP"
	TwoFactorByIp                            = "TWO_FACTOR_BY_IP"
	RestrictDomainAccess                     = "RESTRICT_DOMAIN_ACCESS"
	RestrictDomainCreate                     = "RESTRICT_DOMAIN_CREATE"
	RestrictHoverLocks                       = "RESTRICT_HOVER_LOCKS"
	RestrictPromptToLogin                    = "RESTRICT_PROMPT_TO_LOGIN"
	RestrictPromptToFill                     = "RESTRICT_PROMPT_TO_FILL"
	RestrictAutoSubmit                       = "RESTRICT_AUTO_SUBMIT"
	RestrictPromptToSave                     = "RESTRICT_PROMPT_TO_SAVE"
	RestrictPromptToChange                   = "RESTRICT_PROMPT_TO_CHANGE"
	RestrictAutoFill                         = "RESTRICT_AUTO_FILL"
	RestrictPromptToDisable                  = "RESTRICT_PROMPT_TO_DISABLE"
	RestrictHttpFillWarning                  = "RESTRICT_HTTP_FILL_WARNING"
	KeeperFillHoverLocks                     = "KEEPER_FILL_HOVER_LOCKS"
	KeeperFillAutoFill                       = "KEEPER_FILL_AUTO_FILL"
	KeeperFillAutoSubmit                     = "KEEPER_FILL_AUTO_SUBMIT"
	KeeperFillMatchOnSubdomain               = "KEEPER_FILL_MATCH_ON_SUBDOMAIN"
	KeeperFillAutoSuggest                    = "KEEPER_FILL_AUTO_SUGGEST"
	RestrictCreateFolder                     = "RESTRICT_CREATE_FOLDER"
	RestrictCreateFolderToOnlySharedFolders  = "RESTRICT_CREATE_FOLDER_TO_ONLY_SHARED_FOLDERS"
	RestrictCreateRecord                     = "RESTRICT_CREATE_RECORD"
	RestrictCreateRecordToSharedFolders      = "RESTRICT_CREATE_RECORD_TO_SHARED_FOLDERS"
	RestrictCreateSharedFolder               = "RESTRICT_CREATE_SHARED_FOLDER"
	RestrictCreateIdentityPaymentRecords     = "RESTRICT_CREATE_IDENTITY_PAYMENT_RECORDS"
	DisableCreateDuplicate                   = "DISABLE_CREATE_DUPLICATE"
	MaskCustomFields                         = "MASK_CUSTOM_FIELDS"
	MaskNotes                                = "MASK_NOTES"
	MaskPasswordsWhileEditing                = "MASK_PASSWORDS_WHILE_EDITING"
	GeneratedPasswordComplexity              = "GENERATED_PASSWORD_COMPLEXITY"
	GeneratedSecurityQuestionComplexity      = "GENERATED_SECURITY_QUESTION_COMPLEXITY"
	DaysBeforeDeletedRecordsClearedPerm      = "DAYS_BEFORE_DELETED_RECORDS_CLEARED_PERM"
	DaysBeforeDeletedRecordsAutoCleared      = "DAYS_BEFORE_DELETED_RECORDS_AUTO_CLEARED"
	RestrictRecordTypes                      = "RESTRICT_RECORD_TYPES"
	MaximumRecordSize                        = "MAXIMUM_RECORD_SIZE"
	RestrictSfRecordRemoval                  = "RESTRICT_SF_RECORD_REMOVAL"
	RestrictSfFolderDeletion                 = "RESTRICT_SF_FOLDER_DELETION"
	AllowAlternatePasswords                  = "ALLOW_ALTERNATE_PASSWORDS"
	RestrictPersonalLicense                  = "RESTRICT_PERSONAL_LICENSE"
	RestrictAccountRecovery                  = "RESTRICT_ACCOUNT_RECOVERY"
	RestrictAccountSwitching                 = "RESTRICT_ACCOUNT_SWITCHING"
	RequireSecurityKeyPin                    = "REQUIRE_SECURITY_KEY_PIN"
	RestrictPasskeyLogin                     = "RESTRICT_PASSKEY_LOGIN"
	RestrictLinkSharing                      = "RESTRICT_LINK_SHARING"
	DisableSetupTour                         = "DISABLE_SETUP_TOUR"
	DisableOnboarding                        = "DISABLE_ONBOARDING"
	DisallowV2Clients                        = "DISALLOW_V2_CLIENTS"
	RestrictEmailChange                      = "RESTRICT_EMAIL_CHANGE"
	SendInviteAtRegistration                 = "SEND_INVITE_AT_REGISTRATION"
	ResendEnterpriseInviteInXDays            = "RESEND_ENTERPRISE_INVITE_IN_X_DAYS"
	AutomaticBackupEveryXDays                = "AUTOMATIC_BACKUP_EVERY_X_DAYS"
	SendBreachWatchEvents                    = "SEND_BREACH_WATCH_EVENTS"
	RestrictBreachWatch                      = "RESTRICT_BREACH_WATCH"
	AllowPamRotation                         = "ALLOW_PAM_ROTATION"
	AllowPamDiscovery                        = "ALLOW_PAM_DISCOVERY"
	AllowPamGateway                          = "ALLOW_PAM_GATEWAY"
	AllowConfigureRotationSettings           = "ALLOW_CONFIGURE_ROTATION_SETTINGS"
	AllowRotateCredentials                   = "ALLOW_ROTATE_CREDENTIALS"
	AllowConfigurePamCloudConnectionSettings = "ALLOW_CONFIGURE_PAM_CLOUD_CONNECTION_SETTINGS"
	AllowLaunchPamOnCloudConnection          = "ALLOW_LAUNCH_PAM_ON_CLOUD_CONNECTION"
	AllowConfigurePamTunnelingSettings       = "ALLOW_CONFIGURE_PAM_TUNNELING_SETTINGS"
	AllowLaunchPamTunnels                    = "ALLOW_LAUNCH_PAM_TUNNELS"
	AllowLaunchRbi                           = "ALLOW_LAUNCH_RBI"
	AllowConfigureRbi                        = "ALLOW_CONFIGURE_RBI"
	AllowViewKcmRecordings                   = "ALLOW_VIEW_KCM_RECORDINGS"
	AllowViewRbiRecordings                   = "ALLOW_VIEW_RBI_RECORDINGS"
	RequireSelfDestruct                      = "REQUIRE_SELF_DESTRUCT"
	RestrictSelfDestructRecords              = "RESTRICT_SELF_DESTRUCT_RECORDS"
	AllowSecretsManager                      = "ALLOW_SECRETS_MANAGER"
	AllowCanEditExternalShares               = "ALLOW_CAN_EDIT_EXTERNAL_SHARES"
	RestrictSnapshotTool                     = "RESTRICT_SNAPSHOT_TOOL"
	RestrictForcefield                       = "RESTRICT_FORCEFIELD"
	RestrictClipboardExpireInXSecs           = "RESTRICT_CLIPBOARD_EXPIRE_IN_X_SECS"
	RestrictManageTla                        = "RESTRICT_MANAGE_TLA"
)

// Managing node privilege constants.
const (
	PrivilegeManageNodes          = "manage_nodes"
	PrivilegeManageUser           = "manage_user"
	PrivilegeManageRoles          = "manage_roles"
	PrivilegeManageTeams          = "manage_teams"
	PrivilegeRunReports           = "run_reports"
	PrivilegeManageBridge         = "manage_bridge"
	PrivilegeApproveDevice        = "approve_device"
	PrivilegeManageRecordTypes    = "manage_record_types"
	PrivilegeRunComplianceReports = "run_compliance_reports"
	PrivilegeTransferAccount      = "transfer_account"
	PrivilegeSharingAdministrator = "sharing_administrator"
	PrivilegeManageCompanies      = "manage_companies"
)

// ValidPrivileges contains all valid privilege values as a slice.
var ValidPrivileges = []string{
	PrivilegeManageNodes,
	PrivilegeManageUser,
	PrivilegeManageRoles,
	PrivilegeManageTeams,
	PrivilegeRunReports,
	PrivilegeManageBridge,
	PrivilegeApproveDevice,
	PrivilegeManageRecordTypes,
	PrivilegeRunComplianceReports,
	PrivilegeTransferAccount,
	PrivilegeSharingAdministrator,
	PrivilegeManageCompanies,
}

// ValidEnforcementPolicyKeys contains all valid enforcement policy keys as a slice.
var ValidEnforcementPolicyKeys = []string{
	MasterPasswordMinimumLength,
	MasterPasswordMinimumSpecial,
	MasterPasswordMinimumUpper,
	MasterPasswordMinimumLower,
	MasterPasswordMinimumDigits,
	MasterPasswordRestrictDaysBeforeReuse,
	RequireTwoFactor,
	MasterPasswordMaximumDaysBeforeChange,
	MasterPasswordExpiredAsOf,
	MinimumPbkdf2Iterations,
	MaxSessionLoginTime,
	RestrictPersistentLogin,
	StayLoggedInDefault,
	RestrictSharingAll,
	RestrictSharingEnterprise,
	RestrictSharingAllOutgoing,
	RestrictSharingEnterpriseOutgoing,
	RestrictExport,
	RestrictFileUpload,
	RequireAccountShare,
	RestrictSharingAllIncoming,
	RestrictSharingEnterpriseIncoming,
	RestrictSharingRecordWithAttachments,
	RestrictIpAddresses,
	RequireDeviceApproval,
	RequireAccountRecoveryApproval,
	RestrictVaultIpAddresses,
	TipZoneRestrictAllowedIpRanges,
	AutomaticBackupEveryXDays,
	RestrictOfflineAccess,
	SendInviteAtRegistration,
	RestrictEmailChange,
	RestrictIosFingerprint,
	RestrictMacFingerprint,
	RestrictAndroidFingerprint,
	RestrictWindowsFingerprint,
	LogoutTimerWeb,
	LogoutTimerMobile,
	LogoutTimerDesktop,
	RestrictWebVaultAccess,
	RestrictExtensionsAccess,
	RestrictMobileAccess,
	RestrictDesktopAccess,
	RestrictMobileIosAccess,
	RestrictMobileAndroidAccess,
	RestrictMobileWindowsPhoneAccess,
	RestrictDesktopWinAccess,
	RestrictDesktopMacAccess,
	RestrictChatDesktopAccess,
	RestrictChatMobileAccess,
	RestrictCommanderAccess,
	RestrictTwoFactorChannelText,
	RestrictTwoFactorChannelGoogle,
	RestrictTwoFactorChannelDna,
	RestrictTwoFactorChannelDuo,
	RestrictTwoFactorChannelRsa,
	TwoFactorDurationWeb,
	TwoFactorDurationMobile,
	TwoFactorDurationDesktop,
	RestrictTwoFactorChannelSecurityKeys,
	TwoFactorByIp,
	RestrictDomainAccess,
	RestrictDomainCreate,
	RestrictHoverLocks,
	RestrictPromptToLogin,
	RestrictPromptToFill,
	RestrictAutoSubmit,
	RestrictPromptToSave,
	RestrictPromptToChange,
	RestrictAutoFill,
	RestrictCreateFolder,
	RestrictCreateFolderToOnlySharedFolders,
	RestrictCreateIdentityPaymentRecords,
	MaskCustomFields,
	MaskNotes,
	MaskPasswordsWhileEditing,
	GeneratedPasswordComplexity,
	GeneratedSecurityQuestionComplexity,
	RestrictImport,
	DaysBeforeDeletedRecordsClearedPerm,
	DaysBeforeDeletedRecordsAutoCleared,
	AllowAlternatePasswords,
	RestrictCreateRecord,
	RestrictCreateRecordToSharedFolders,
	RestrictCreateSharedFolder,
	RestrictLinkSharing,
	RestrictSharingOutsideOfIsolatedNodes,
	RestrictSharingRecordToSharedFolders,
	DisableSetupTour,
	RestrictPersonalLicense,
	DisableOnboarding,
	DisallowV2Clients,
	RestrictIpAutoapproval,
	SendBreachWatchEvents,
	RestrictBreachWatch,
	ResendEnterpriseInviteInXDays,
	MasterPasswordReentry,
	RestrictAccountRecovery,
	KeeperFillHoverLocks,
	KeeperFillAutoFill,
	KeeperFillAutoSubmit,
	KeeperFillMatchOnSubdomain,
	KeeperFillAutoSuggest,
	RestrictPromptToDisable,
	RestrictHttpFillWarning,
	RestrictRecordTypes,
	AllowSecretsManager,
	RequireSelfDestruct,
	MaximumRecordSize,
	AllowPamRotation,
	AllowPamDiscovery,
	RestrictImportSharedFolders,
	RequireSecurityKeyPin,
	DisableCreateDuplicate,
	AllowPamGateway,
	AllowConfigureRotationSettings,
	AllowRotateCredentials,
	AllowConfigurePamCloudConnectionSettings,
	AllowLaunchPamOnCloudConnection,
	AllowConfigurePamTunnelingSettings,
	AllowLaunchPamTunnels,
	AllowLaunchRbi,
	AllowConfigureRbi,
	AllowViewKcmRecordings,
	AllowViewRbiRecordings,
	RestrictManageTla,
	RestrictSelfDestructRecords,
	RestrictAccountSwitching,
	RestrictPasskeyLogin,
	AllowCanEditExternalShares,
	RestrictSnapshotTool,
	RestrictForcefield,
	RestrictClipboardExpireInXSecs,
	RestrictSfRecordRemoval,
	RestrictSfFolderDeletion,
}
