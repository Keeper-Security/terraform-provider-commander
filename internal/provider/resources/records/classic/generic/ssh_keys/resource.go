// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package sshkeys

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.Resource = &SshKeysResource{}
var _ resource.ResourceWithConfigure = &SshKeysResource{}
var _ resource.ResourceWithImportState = &SshKeysResource{}

type SshKeysResource struct {
	utils.BaseResource
}

func (r *SshKeysResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_classic_ssh_keys"
}

func (r *SshKeysResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.ConfigureResource(ctx, req, resp)
}

func NewSshKeysResource() resource.Resource {
	return &SshKeysResource{}
}
