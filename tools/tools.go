// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

//go:build generate

package tools

import (
	_ "github.com/hashicorp/copywrite"
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"
)

// Generate copyright headers
//go:generate go run github.com/hashicorp/copywrite headers -d .. --config ../.copywrite.hcl

// Format Terraform code for use in documentation.
// If you do not have Terraform installed, you can remove the formatting command, but it is suggested
// to ensure the documentation is formatted properly.
//go:generate terraform fmt -recursive ../examples/

// Generate documentation templates with Registry subcategories, then render docs.
//go:generate sh -c "cd .. && go run ./docs/generate_doc_templates.go"
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-dir .. -provider-name commander
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs validate --provider-dir .. -provider-name commander --allowed-resource-subcategories-file ../docs/subcategories.txt
