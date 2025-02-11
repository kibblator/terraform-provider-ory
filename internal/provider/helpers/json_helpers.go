package helpers

import (
	"context"
)

func GetNested(ctx context.Context, m map[string]interface{}, keys ...string) interface{} {
	var current interface{} = m

	for _, key := range keys {
		if nestedMap, ok := current.(map[string]interface{}); ok {
			if val, exists := nestedMap[key]; exists {
				// tflog.Debug(ctx, "Found key", map[string]interface{}{"key": key, "value": val})
				current = val
			} else {
				// tflog.Debug(ctx, "Key not found", map[string]interface{}{"key": key})
				return nil
			}
		} else {
			// tflog.Debug(ctx, "Current value is not a map", map[string]interface{}{"current": current})
			return nil
		}
	}

	return current
}
