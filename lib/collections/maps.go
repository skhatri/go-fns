package collections

func FilteredByKey[K comparable, V any](m map[K]V, pred func(key K) bool) map[K]V {
	filteredMap := make(map[K]V, 0)
	for k, v := range m {
		if pred(k) {
			filteredMap[k] = v
		}
	}
	return filteredMap
}

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

func CopyMap(dest map[string]interface{}, src map[string]interface{}) {
	for k, v := range src {
		dest[k] = v
	}
}
