// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package paymentcard

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.Resource = &PaymentCardResource{}
var _ resource.ResourceWithConfigure = &PaymentCardResource{}
var _ resource.ResourceWithImportState = &PaymentCardResource{}

type PaymentCardResource struct {
	utils.BaseResource
}

func (r *PaymentCardResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_new_payment_card"
}

func (r *PaymentCardResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.ConfigureResource(ctx, req, resp)
}

func NewPaymentCardResource() resource.Resource {
	return &PaymentCardResource{}
}
