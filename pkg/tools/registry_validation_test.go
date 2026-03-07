package tools

import (
	"context"
	"testing"
)

// MockToolWithParams is a mock tool with parameter validation
type MockToolWithParams struct {
	name        string
	description string
	params      map[string]any
	executed    bool
	lastArgs    map[string]any
}

func (m *MockToolWithParams) Name() string        { return m.name }
func (m *MockToolWithParams) Description() string { return m.description }
func (m *MockToolWithParams) Parameters() map[string]any {
	return m.params
}

func (m *MockToolWithParams) Execute(ctx context.Context, args map[string]any) *ToolResult {
	m.executed = true
	m.lastArgs = args
	return NewToolResult("executed successfully")
}

func TestToolArgumentValidation(t *testing.T) {
	t.Run("missing required parameter", func(t *testing.T) {
		registry := NewToolRegistry()
		tool := &MockToolWithParams{
			name:        "test_tool",
			description: "Test tool",
			params: map[string]any{
				"required": []string{"param1", "param2"},
				"properties": map[string]any{
					"param1": map[string]any{"type": "string"},
					"param2": map[string]any{"type": "number"},
				},
			},
		}
		registry.Register(tool)

		// Missing param2
		result := registry.Execute(context.Background(), "test_tool", map[string]any{
			"param1": "value",
		})

		if !result.IsError {
			t.Error("Expected error for missing required parameter")
		}
		if result.ForLLM == "" || result.Err == nil {
			t.Error("Expected error message and error object")
		}
		t.Logf("Error message: %s", result.ForLLM)
	})

	t.Run("wrong parameter type - string expected", func(t *testing.T) {
		registry := NewToolRegistry()
		tool := &MockToolWithParams{
			name:        "test_tool",
			description: "Test tool",
			params: map[string]any{
				"required": []string{"param1"},
				"properties": map[string]any{
					"param1": map[string]any{"type": "string"},
				},
			},
		}
		registry.Register(tool)

		result := registry.Execute(context.Background(), "test_tool", map[string]any{
			"param1": 123, // Should be string
		})

		if !result.IsError {
			t.Error("Expected error for wrong parameter type")
		}
		t.Logf("Error message: %s", result.ForLLM)
	})

	t.Run("wrong parameter type - number expected", func(t *testing.T) {
		registry := NewToolRegistry()
		tool := &MockToolWithParams{
			name:        "test_tool",
			description: "Test tool",
			params: map[string]any{
				"required": []string{"count"},
				"properties": map[string]any{
					"count": map[string]any{"type": "number"},
				},
			},
		}
		registry.Register(tool)

		result := registry.Execute(context.Background(), "test_tool", map[string]any{
			"count": "not a number",
		})

		if !result.IsError {
			t.Error("Expected error for wrong parameter type")
		}
		t.Logf("Error message: %s", result.ForLLM)
	})

	t.Run("wrong parameter type - integer expected", func(t *testing.T) {
		registry := NewToolRegistry()
		tool := &MockToolWithParams{
			name:        "test_tool",
			description: "Test tool",
			params: map[string]any{
				"required": []string{"count"},
				"properties": map[string]any{
					"count": map[string]any{"type": "integer"},
				},
			},
		}
		registry.Register(tool)

		// Float that's not a whole number
		result := registry.Execute(context.Background(), "test_tool", map[string]any{
			"count": 3.14,
		})

		if !result.IsError {
			t.Error("Expected error for non-integer float")
		}
		t.Logf("Error message: %s", result.ForLLM)
	})

	t.Run("integer accepts whole number float", func(t *testing.T) {
		registry := NewToolRegistry()
		tool := &MockToolWithParams{
			name:        "test_tool",
			description: "Test tool",
			params: map[string]any{
				"required": []string{"count"},
				"properties": map[string]any{
					"count": map[string]any{"type": "integer"},
				},
			},
		}
		registry.Register(tool)

		// Float that IS a whole number
		result := registry.Execute(context.Background(), "test_tool", map[string]any{
			"count": 42.0,
		})

		if result.IsError {
			t.Errorf("Should accept whole number float: %s", result.ForLLM)
		}
	})

	t.Run("wrong parameter type - boolean expected", func(t *testing.T) {
		registry := NewToolRegistry()
		tool := &MockToolWithParams{
			name:        "test_tool",
			description: "Test tool",
			params: map[string]any{
				"required": []string{"enabled"},
				"properties": map[string]any{
					"enabled": map[string]any{"type": "boolean"},
				},
			},
		}
		registry.Register(tool)

		result := registry.Execute(context.Background(), "test_tool", map[string]any{
			"enabled": "true", // Should be bool, not string
		})

		if !result.IsError {
			t.Error("Expected error for wrong parameter type")
		}
		t.Logf("Error message: %s", result.ForLLM)
	})

	t.Run("wrong parameter type - array expected", func(t *testing.T) {
		registry := NewToolRegistry()
		tool := &MockToolWithParams{
			name:        "test_tool",
			description: "Test tool",
			params: map[string]any{
				"required": []string{"items"},
				"properties": map[string]any{
					"items": map[string]any{"type": "array"},
				},
			},
		}
		registry.Register(tool)

		result := registry.Execute(context.Background(), "test_tool", map[string]any{
			"items": "not an array",
		})

		if !result.IsError {
			t.Error("Expected error for wrong parameter type")
		}
		t.Logf("Error message: %s", result.ForLLM)
	})

	t.Run("wrong parameter type - object expected", func(t *testing.T) {
		registry := NewToolRegistry()
		tool := &MockToolWithParams{
			name:        "test_tool",
			description: "Test tool",
			params: map[string]any{
				"required": []string{"config"},
				"properties": map[string]any{
					"config": map[string]any{"type": "object"},
				},
			},
		}
		registry.Register(tool)

		result := registry.Execute(context.Background(), "test_tool", map[string]any{
			"config": "not an object",
		})

		if !result.IsError {
			t.Error("Expected error for wrong parameter type")
		}
		t.Logf("Error message: %s", result.ForLLM)
	})

	t.Run("valid arguments - all types", func(t *testing.T) {
		registry := NewToolRegistry()
		tool := &MockToolWithParams{
			name:        "test_tool",
			description: "Test tool",
			params: map[string]any{
				"required": []string{"str", "num", "int", "bool", "arr", "obj"},
				"properties": map[string]any{
					"str":  map[string]any{"type": "string"},
					"num":  map[string]any{"type": "number"},
					"int":  map[string]any{"type": "integer"},
					"bool": map[string]any{"type": "boolean"},
					"arr":  map[string]any{"type": "array"},
					"obj":  map[string]any{"type": "object"},
				},
			},
		}
		registry.Register(tool)

		result := registry.Execute(context.Background(), "test_tool", map[string]any{
			"str":  "hello",
			"num":  3.14,
			"int":  42,
			"bool": true,
			"arr":  []string{"a", "b"},
			"obj":  map[string]any{"key": "value"},
		})

		if result.IsError {
			t.Errorf("Should accept valid arguments: %s", result.ForLLM)
		}
		if !tool.executed {
			t.Error("Tool should have been executed")
		}
	})

	t.Run("null values allowed", func(t *testing.T) {
		registry := NewToolRegistry()
		tool := &MockToolWithParams{
			name:        "test_tool",
			description: "Test tool",
			params: map[string]any{
				"required": []string{}, // No required params
				"properties": map[string]any{
					"optional": map[string]any{"type": "string"},
				},
			},
		}
		registry.Register(tool)

		result := registry.Execute(context.Background(), "test_tool", map[string]any{
			"optional": nil,
		})

		if result.IsError {
			t.Errorf("Should accept null values: %s", result.ForLLM)
		}
	})

	t.Run("unknown parameter warning", func(t *testing.T) {
		registry := NewToolRegistry()
		tool := &MockToolWithParams{
			name:        "test_tool",
			description: "Test tool",
			params: map[string]any{
				"required": []string{"param1"},
				"properties": map[string]any{
					"param1": map[string]any{"type": "string"},
				},
			},
		}
		registry.Register(tool)

		// Include unknown parameter (should log warning but not fail)
		result := registry.Execute(context.Background(), "test_tool", map[string]any{
			"param1":  "value",
			"unknown": "extra",
		})

		if result.IsError {
			t.Errorf("Should not fail on unknown parameter: %s", result.ForLLM)
		}
	})

	t.Run("no validation when no parameters defined", func(t *testing.T) {
		registry := NewToolRegistry()
		tool := &MockToolWithParams{
			name:        "test_tool",
			description: "Test tool",
			params:      nil, // No parameters
		}
		registry.Register(tool)

		result := registry.Execute(context.Background(), "test_tool", map[string]any{
			"anything": "goes",
		})

		if result.IsError {
			t.Errorf("Should not validate when no parameters defined: %s", result.ForLLM)
		}
	})

	t.Run("required as interface slice", func(t *testing.T) {
		registry := NewToolRegistry()
		tool := &MockToolWithParams{
			name:        "test_tool",
			description: "Test tool",
			params: map[string]any{
				"required": []interface{}{"param1", "param2"}, // interface{} slice
				"properties": map[string]any{
					"param1": map[string]any{"type": "string"},
					"param2": map[string]any{"type": "string"},
				},
			},
		}
		registry.Register(tool)

		// Missing param2
		result := registry.Execute(context.Background(), "test_tool", map[string]any{
			"param1": "value",
		})

		if !result.IsError {
			t.Error("Expected error for missing required parameter")
		}
	})

	t.Run("number types accepted", func(t *testing.T) {
		registry := NewToolRegistry()
		tool := &MockToolWithParams{
			name:        "test_tool",
			description: "Test tool",
			params: map[string]any{
				"required": []string{"num"},
				"properties": map[string]any{
					"num": map[string]any{"type": "number"},
				},
			},
		}
		registry.Register(tool)

		testCases := []struct {
			name  string
			value any
		}{
			{"float64", float64(3.14)},
			{"float32", float32(3.14)},
			{"int", int(42)},
			{"int32", int32(42)},
			{"int64", int64(42)},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				result := registry.Execute(context.Background(), "test_tool", map[string]any{
					"num": tc.value,
				})

				if result.IsError {
					t.Errorf("Should accept %s: %s", tc.name, result.ForLLM)
				}
			})
		}
	})

	t.Run("array types accepted", func(t *testing.T) {
		registry := NewToolRegistry()
		tool := &MockToolWithParams{
			name:        "test_tool",
			description: "Test tool",
			params: map[string]any{
				"required": []string{"arr"},
				"properties": map[string]any{
					"arr": map[string]any{"type": "array"},
				},
			},
		}
		registry.Register(tool)

		testCases := []struct {
			name  string
			value any
		}{
			{"[]any", []any{"a", "b"}},
			{"[]string", []string{"a", "b"}},
			{"[]int", []int{1, 2}},
			{"[]float64", []float64{1.1, 2.2}},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				result := registry.Execute(context.Background(), "test_tool", map[string]any{
					"arr": tc.value,
				})

				if result.IsError {
					t.Errorf("Should accept %s: %s", tc.name, result.ForLLM)
				}
			})
		}
	})
}
