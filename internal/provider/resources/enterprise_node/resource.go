// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisenode

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &EnterpriseNodeResource{}
var _ resource.ResourceWithConfigure = &EnterpriseNodeResource{}

type EnterpriseNodeResource struct {
	apiManager *api.ApiManager
}

type EnterpriseNodeResourceModel struct {
	Id             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Parent         types.String `tfsdk:"parent"`
	WipeOut        types.Bool   `tfsdk:"wipe_out"`
	ToggleIsolated types.Bool   `tfsdk:"toggle_isolated"`
	// LogoFile       types.String `tfsdk:"logo_file"` // NOTE: In commander cli not working and in admin console there no feature like this
	ManagedCompany types.String `tfsdk:"managed_company"`
}

func (r *EnterpriseNodeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_enterprise_node"
}

func (r *EnterpriseNodeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					nameValidator{},
				},
			},
			"parent": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					parentValidator{},
				},
			},
			"wipe_out": schema.BoolAttribute{
				Optional: true,
			},
			"toggle_isolated": schema.BoolAttribute{
				Optional: true,
			},
			// "logo_file": schema.StringAttribute{
			// 	Optional: true,
			// },
			"managed_company": schema.StringAttribute{
				Optional:   true,
				Validators: []validator.String{
					// managedCompanyValidator{},
				},
			},
		},
	}
}

func (r *EnterpriseNodeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	apiManager, ok := req.ProviderData.(*api.ApiManager)
	if !ok {
		resp.Diagnostics.AddError(
			"Provider Configuration Error",
			fmt.Sprintf("The provider was not configured correctly. Expected API manager, but got: %T. Please check your provider configuration.", req.ProviderData),
		)
		return
	}

	r.apiManager = apiManager

}

// ensureApiManager validates that apiManager is configured and returns an error if not
func (r *EnterpriseNodeResource) ensureApiManager() error {
	if r.apiManager == nil {
		return fmt.Errorf("the Keeper Commander provider is not properly configured. Please ensure the provider is set up with valid service_mode_url and service_mode_api_key")
	}
	return nil
}

func (r *EnterpriseNodeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data EnterpriseNodeResourceModel

	// Get planned data from Terraform
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate ApiManager is configured
	if err := r.ensureApiManager(); err != nil {
		resp.Diagnostics.AddError(
			"Provider Configuration Error",
			err.Error(),
		)
		return
	}

	r.apiManager.ExecuteCommand(ctx, "msp-down", "Unable to msp down latest changes")

	// Execute with managed company context if provided
	err := utils.ExecuteWithManagedCompanyContext(ctx, r.apiManager, data.ManagedCompany, func() error {
		// Build the Commander command string
		command := buildEnterpriseNodeAddCommand(data)

		createNodeResponse, err := r.apiManager.ExecuteCommand(ctx, command, "Unable to create enterprise node")
		if err != nil {
			return fmt.Errorf("create enterprise node failed: %w", err)
		}

		nodeID, ok := utils.ExtractNodeIDFromCreateNodeResponse(string(createNodeResponse.Message))

		idValue := data.Name.ValueString()

		if ok {
			idValue = nodeID
		}

		data.Id = types.StringValue(idValue)

		// Set the ID in the state
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		return nil
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Create Enterprise Node Failed",
			err.Error(),
		)
		return
	}
}
func (r *EnterpriseNodeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state EnterpriseNodeResourceModel

	// Get current state (old values)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate ApiManager is configured
	if err := r.ensureApiManager(); err != nil {
		resp.Diagnostics.AddError(
			"Provider Configuration Error",
			err.Error(),
		)
		return
	}

	// Execute msp-down command (setup/initialization)
	r.apiManager.ExecuteCommand(ctx, "msp-down", "Unable to msp down latest changes")

	// Execute with managed company context if provided
	err := utils.ExecuteWithManagedCompanyContext(ctx, r.apiManager, state.ManagedCompany, func() error {
		// Build the Commander command string
		command := fmt.Sprintf("enterprise-info -n -v --format json --node '%s'", state.Id.ValueString())

		apiResp, err := r.apiManager.ExecuteCommand(ctx, command, "Unable to read enterprise node")
		if err != nil {
			return fmt.Errorf("Read Enterprise Node Failed: %w", err)
		}

		// Parse the JSON response - it's an array of node objects
		var nodes []struct {
			NodeId         int    `json:"node_id"`
			Name           string `json:"name"`
			ParentNodeName string `json:"parent_node"`
			ParentNodeId   int    `json:"parent_id"`
		}

		// Convert apiResp.Data to JSON bytes and unmarshal
		dataBytes, err := json.Marshal(apiResp.Data)
		if err != nil {
			return fmt.Errorf("unable to process the response from Keeper Commander API: %w", err)
		}

		if err := json.Unmarshal(dataBytes, &nodes); err != nil {
			return fmt.Errorf("unable to parse enterprise nodes list from API response: %w", err)
		}

		// Find the node matching state.Id (which is the node name)
		var nodeInfo *struct {
			NodeId         int    `json:"node_id"`
			Name           string `json:"name"`
			ParentNodeName string `json:"parent_node"`
			ParentNodeId   int    `json:"parent_id"`
		}

		stateId := state.Id.ValueString()
		for i := range nodes {
			if strconv.Itoa(nodes[i].NodeId) == stateId {
				nodeInfo = &nodes[i]
				break
			}
		}

		if nodeInfo == nil {
			// Resource not found - remove from state
			resp.State.RemoveResource(ctx)
			return nil
		}

		// Map the response to the state
		state.Id = types.StringValue(strconv.Itoa(nodeInfo.NodeId))
		state.Name = types.StringValue(nodeInfo.Name)

		// TODO: here we need to check with client if we need to set the parent node name or id bec previous state value is cant be accessed so we will not know if parent id is selecetd or name
		// Set parent if it exists in the response
		state.Parent = types.StringValue(nodeInfo.ParentNodeName)

		// Set the updated state
		resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
		return nil
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Read Enterprise Node Failed",
			err.Error(),
		)
		return
	}
}
func (r *EnterpriseNodeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan EnterpriseNodeResourceModel
	var state EnterpriseNodeResourceModel

	// Get planned data (new values)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state (old values)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate ApiManager is configured
	if err := r.ensureApiManager(); err != nil {
		resp.Diagnostics.AddError(
			"Provider Configuration Error",
			err.Error(),
		)
		return
	}

	// Use managed company from plan (or state if plan doesn't have it)
	managedCompany := plan.ManagedCompany
	if managedCompany.IsNull() || managedCompany.IsUnknown() {
		managedCompany = state.ManagedCompany
	}

	// Execute with managed company context if provided
	err := utils.ExecuteWithManagedCompanyContext(ctx, r.apiManager, managedCompany, func() error {
		command := buildEnterpriseNodeUpdateCommand(&plan, &state)

		_, err := r.apiManager.ExecuteCommand(ctx, command, "Unable to update enterprise node")
		if err != nil {
			return fmt.Errorf("update enterprise node failed: %w", err)
		}

		// Keep the same ID
		plan.Id = state.Id
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
		return nil
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Update Enterprise Node Failed",
			err.Error(),
		)
		return
	}
}
func (r *EnterpriseNodeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state EnterpriseNodeResourceModel

	// Get state from Terraform
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate ApiManager is configured
	if err := r.ensureApiManager(); err != nil {
		resp.Diagnostics.AddError(
			"Provider Configuration Error",
			err.Error(),
		)
		return
	}

	// Execute with managed company context if provided
	err := utils.ExecuteWithManagedCompanyContext(ctx, r.apiManager, state.ManagedCompany, func() error {
		// Build delete command
		command := fmt.Sprintf("enterprise-node --delete '%s'", state.Id.ValueString())

		_, err := r.apiManager.ExecuteCommand(ctx, command, "Unable to delete enterprise node")
		if err != nil {
			return fmt.Errorf("Delete Enterprise Node Failed: %w", err)
		}
		return nil
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Delete Enterprise Node Failed",
			err.Error(),
		)
		return
	}
}

func NewEnterpriseNodeResource() resource.Resource {
	return &EnterpriseNodeResource{}
}
