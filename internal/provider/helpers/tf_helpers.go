package helpers

import "github.com/hashicorp/terraform-plugin-framework/types"

func StringOrNil(value string) types.String {
	if value == "" {
		return types.StringPointerValue(nil)
	}
	return types.StringValue(value)
}
