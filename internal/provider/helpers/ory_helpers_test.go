package helpers

import (
	"testing"

	orytypes "github.com/kibblator/terraform-provider-ory/internal/provider/types"

	"github.com/stretchr/testify/assert"
)

func TestConvertToHooks(t *testing.T) {
	expectedHooks := []orytypes.Hook{
		{Hook: "session"},
	}

	t.Run("Valid hooks", func(t *testing.T) {
		hooksRaw := []interface{}{
			map[string]interface{}{
				"hook": "session",
			},
		}

		hooks, diags := ConvertToHooks(hooksRaw)
		assert.Equal(t, expectedHooks, hooks)
		assert.False(t, diags.HasError())
	})

	t.Run("Invalid hooks", func(t *testing.T) {
		invalidHooksRaw := "invalid"
		_, diags := ConvertToHooks(invalidHooksRaw)
		assert.True(t, diags.HasError())
	})

	t.Run("Invalid hooks shape", func(t *testing.T) {
		incorrectShapeHooksRaw := []interface{}{
			map[string]interface{}{
				"invalid_key": "value",
			},
		}
		_, diags := ConvertToHooks(incorrectShapeHooksRaw)
		assert.True(t, diags.HasError())
	})
}
