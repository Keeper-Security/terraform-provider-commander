// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package bankaccount

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.Resource = &BankAccountResource{}
var _ resource.ResourceWithConfigure = &BankAccountResource{}
var _ resource.ResourceWithImportState = &BankAccountResource{}

type BankAccountResource struct {
	utils.BaseResource
}

func (r *BankAccountResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_classic_bank_account"
}

func (r *BankAccountResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.ConfigureResource(ctx, req, resp)
}

func NewBankAccountResource() resource.Resource {
	return &BankAccountResource{}
}
