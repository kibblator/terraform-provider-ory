package custom_validators

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type Base64Validator struct{}

func (b Base64Validator) Description(_ context.Context) string {
	return "Ensures the string is a valid base64-encoded value"
}

func (b Base64Validator) MarkdownDescription(_ context.Context) string {
	return "Ensures the string is a **valid base64-encoded** value"
}

func (b Base64Validator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()
	_, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Base64 Encoding",
			fmt.Sprintf("The provided string is not valid base64: %s", err),
		)
	}
}
