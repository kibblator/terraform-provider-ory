package helpers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetNested tests the GetNested function
func TestGetNested(t *testing.T) {
	ctx := context.Background()

	// Test data
	data := map[string]interface{}{
		"level1": map[string]interface{}{
			"level2": map[string]interface{}{
				"level3": "value",
			},
		},
	}

	// Test cases
	tests := []struct {
		name     string
		keys     []string
		expected interface{}
	}{
		{"Valid keys", []string{"level1", "level2", "level3"}, "value"},
		{"Invalid key", []string{"level1", "invalid"}, nil},
		{"Non-map value", []string{"level1", "level2", "level3", "extra"}, nil},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			result := GetNested(ctx, data, testCase.keys...)
			assert.Equal(t, testCase.expected, result)
		})
	}
}
