// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisescimpush

const ResourceDescription = "Push SCIM (Google, AD, or record-based) data to a Keeper SCIM endpoint in one step. Changing any value runs the push again."

const ResourceMarkdownDescription = "Use this resource to **push SCIM data** to a Keeper SCIM endpoint in a single step. You choose where the data comes from (Google Workspace, Active Directory, or a record) and whether to auto-approve teams.\n\n" +
	"## What this resource does\n\n" +
	"- When you **apply**, it runs a one-time SCIM push with your settings.\n" +
	"- Terraform does **not** read or delete anything on the server. It only runs the push and tracks that it happened.\n" +
	"- If you change **scim_id**, **source**, **record**, **auto_approve**, or **managed_company**, Terraform will run the push again with the new values.\n\n"
