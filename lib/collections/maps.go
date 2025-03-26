// Package collections provides utility functions for working with Go collections like maps and sets.
package collections

// CopyAttribute copies a value from a source map to a target map using the specified key.
// If the key exists in the source map, it will be copied to the target map with the same key.
// Returns true if the copy was successful, false otherwise.
func CopyAttribute(attributeName string, source map[string]string, target map[string]string) *string {
	var out *string
	if value, ok := source[attributeName]; ok {
		target[attributeName] = value
		out = &value
	}
	return out
}

// FilteredByKey creates a new map containing only the key-value pairs from the input map
// that satisfy the provided filter function.
// The filter function should return true for keys that should be included in the result.
func FilteredByKey[K comparable, V any](m map[K]V, pred func(key K) bool) map[K]V {
	filteredMap := make(map[K]V)
	for k, v := range m {
		if pred(k) {
			filteredMap[k] = v
		}
	}
	return filteredMap
}

// MapByStringKey creates a new map by applying a transformation function to each value in the input map.
// The transformation function receives the key and value of each entry and returns a new value.
// The keys remain unchanged in the resulting map.
func MapByStringKey(source interface{}) interface{} {
	switch source.(type) {
	case map[interface{}]interface{}:
		m := source.(map[interface{}]interface{})
		result := make(map[string]interface{}, len(m))
		for k, v := range m {
			result[k.(string)] = MapByStringKey(v)
		}
		return result

	case map[string]interface{}:
		target := make(map[string]interface{})

		for k, v := range source.(map[string]interface{}) {
			target[k] = MapByStringKey(v)
		}
		return target
	case []interface{}:
		items := make([]interface{}, 0)
		for _, interfaceItem := range source.([]interface{}) {
			items = append(items, MapByStringKey(interfaceItem))
		}
		return items
	default:
		return source
	}
}

// CopyMap creates a deep copy of the input map.
// It recursively copies nested maps and slices, ensuring that modifications to the copy
// do not affect the original map.
func CopyMap(dest map[string]interface{}, src map[string]interface{}) {
	for k, v := range src {
		dest[k] = v
	}
}
