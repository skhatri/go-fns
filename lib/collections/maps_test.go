package collections

import (
	"testing"
)

func TestCopyAttribute(t *testing.T) {
	tests := []struct {
		name           string
		attributeName  string
		source         map[string]string
		target         map[string]string
		expectedValue  *string
		expectedTarget map[string]string
	}{
		{
			name:          "attribute exists",
			attributeName: "key1",
			source:        map[string]string{"key1": "value1"},
			target:        make(map[string]string),
			expectedValue: stringPtr("value1"),
			expectedTarget: map[string]string{
				"key1": "value1",
			},
		},
		{
			name:           "attribute doesn't exist",
			attributeName:  "nonexistent",
			source:         map[string]string{"key1": "value1"},
			target:         make(map[string]string),
			expectedValue:  nil,
			expectedTarget: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CopyAttribute(tt.attributeName, tt.source, tt.target)
			if tt.expectedValue == nil {
				if result != nil {
					t.Errorf("expected nil result, got %v", *result)
				}
			} else if result == nil {
				t.Errorf("expected value %v, got nil", *tt.expectedValue)
			} else if *result != *tt.expectedValue {
				t.Errorf("expected value %v, got %v", *tt.expectedValue, *result)
			}
			if !mapsEqual(tt.target, tt.expectedTarget) {
				t.Errorf("expected target %v, got %v", tt.expectedTarget, tt.target)
			}
		})
	}
}

func TestMapFilteredByKey(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]int
		pred     func(string) bool
		expected map[string]int
	}{
		{
			name:  "filter even keys",
			input: map[string]int{"a": 1, "bb": 2, "ccc": 3},
			pred: func(key string) bool {
				result := len(key)%2 == 0
				t.Logf("Key: %s, Length: %d, Result: %v", key, len(key), result)
				return result
			},
			expected: map[string]int{"bb": 2},
		},
		{
			name:  "filter keys starting with 'a'",
			input: map[string]int{"a": 1, "ab": 2, "b": 3},
			pred: func(key string) bool {
				return key[0] == 'a'
			},
			expected: map[string]int{"a": 1, "ab": 2},
		},
		{
			name:  "empty map",
			input: map[string]int{},
			pred: func(key string) bool {
				return true
			},
			expected: map[string]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FilteredByKey(tt.input, tt.pred)
			t.Logf("Input map: %v", tt.input)
			t.Logf("Result map: %v", result)
			t.Logf("Expected map: %v", tt.expected)
			if !mapsEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestMapByStringKey(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name: "map[interface{}]interface{} to map[string]interface{}",
			input: map[interface{}]interface{}{
				"key1": "value1",
				"key2": map[interface{}]interface{}{
					"nested": "value2",
				},
			},
			expected: map[string]interface{}{
				"key1": "value1",
				"key2": map[string]interface{}{
					"nested": "value2",
				},
			},
		},
		{
			name: "map[string]interface{} to map[string]interface{}",
			input: map[string]interface{}{
				"key1": "value1",
				"key2": map[string]interface{}{
					"nested": "value2",
				},
			},
			expected: map[string]interface{}{
				"key1": "value1",
				"key2": map[string]interface{}{
					"nested": "value2",
				},
			},
		},
		{
			name: "slice with maps",
			input: []interface{}{
				map[interface{}]interface{}{
					"key1": "value1",
				},
				map[interface{}]interface{}{
					"key2": "value2",
				},
			},
			expected: []interface{}{
				map[string]interface{}{
					"key1": "value1",
				},
				map[string]interface{}{
					"key2": "value2",
				},
			},
		},
		{
			name: "slice with mixed types",
			input: []interface{}{
				"string",
				123,
				map[interface{}]interface{}{"key": "value"},
				map[string]interface{}{"key": "value"},
			},
			expected: []interface{}{
				"string",
				123,
				map[string]interface{}{"key": "value"},
				map[string]interface{}{"key": "value"},
			},
		},
		{
			name:     "non-map value",
			input:    "string value",
			expected: "string value",
		},
		{
			name:     "integer value",
			input:    42,
			expected: 42,
		},
		{
			name:     "nil value",
			input:    nil,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MapByStringKey(tt.input)
			if !deepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestCopyMap(t *testing.T) {
	tests := []struct {
		name     string
		src      map[string]interface{}
		dest     map[string]interface{}
		expected map[string]interface{}
	}{
		{
			name: "copy non-empty map",
			src: map[string]interface{}{
				"key1": "value1",
				"key2": 42,
			},
			dest: make(map[string]interface{}),
			expected: map[string]interface{}{
				"key1": "value1",
				"key2": 42,
			},
		},
		{
			name: "copy to existing map",
			src: map[string]interface{}{
				"key1": "new_value",
			},
			dest: map[string]interface{}{
				"existing": "value",
			},
			expected: map[string]interface{}{
				"existing": "value",
				"key1":     "new_value",
			},
		},
		{
			name:     "copy empty map",
			src:      map[string]interface{}{},
			dest:     make(map[string]interface{}),
			expected: map[string]interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CopyMap(tt.dest, tt.src)
			if !deepEqual(tt.dest, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, tt.dest)
			}
		})
	}
}

// Helper functions for testing
func stringPtr(s string) *string {
	return &s
}

func mapsEqual[K comparable, V comparable](m1, m2 map[K]V) bool {
	if len(m1) != len(m2) {
		return false
	}
	for k, v1 := range m1 {
		if v2, ok := m2[k]; !ok || v1 != v2 {
			return false
		}
	}
	return true
}

func deepEqual(v1, v2 interface{}) bool {
	switch val1 := v1.(type) {
	case map[string]interface{}:
		val2, ok := v2.(map[string]interface{})
		if !ok {
			return false
		}
		if len(val1) != len(val2) {
			return false
		}
		for k, v := range val1 {
			if !deepEqual(v, val2[k]) {
				return false
			}
		}
		return true
	case []interface{}:
		val2, ok := v2.([]interface{})
		if !ok {
			return false
		}
		if len(val1) != len(val2) {
			return false
		}
		for i := range val1 {
			if !deepEqual(val1[i], val2[i]) {
				return false
			}
		}
		return true
	default:
		return v1 == v2
	}
}
