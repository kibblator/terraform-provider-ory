package helpers

import (
	"encoding/json"

	orytypes "terraform-provider-ory/internal/provider/types"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func ConvertToHooks(hooksRaw interface{}) ([]orytypes.Hook, diag.Diagnostics) {
	var diags diag.Diagnostics

	hooksSlice, ok := hooksRaw.([]interface{})
	if !ok {
		diags.AddError("Type Assertion Failed", "Expected hooks to be a slice of interfaces")
		return nil, diags
	}

	hooksJSON, err := json.Marshal(hooksSlice)
	if err != nil {
		diags.AddError("JSON Marshal Failed", "Failed to marshal hooks: "+err.Error())
		return nil, diags
	}

	var hooks []orytypes.Hook
	if err := json.Unmarshal(hooksJSON, &hooks); err != nil {
		diags.AddError("JSON Unmarshal Failed", "Failed to unmarshal hooks into Hook struct: "+err.Error())
		return nil, diags
	}

	return hooks, diags
}
