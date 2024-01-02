package collections

import (
	"testing"
)

func run[A any](input map[string]A, expected int) func(t *testing.T) {
	return func(t *testing.T) {
		result := Values(input)
		if len(result) != expected {
			t.Errorf("empty map expected 0 but found %d", len(result))
		}
	}
}
func TestGetMapValues(t *testing.T) {
	t.Run("Get Values with int", run(map[string]int{"a": 1, "b": 2, "c": 3}, 3))
	t.Run("Get Value of nil", run[string](nil, 0))
	t.Run("Get Value of empty map", run(map[string]string{}, 0))
}

func TestNewSetCompare(t *testing.T) {
	set1 := NewSet[string]([]string{"a", "b", "c"})
	contains := []string{"c", "b", "a"}
	for _, c := range contains {
		if !set1.Contains(c) {
			t.Errorf("expected %s but not found", c)
		}
	}
	doesNotContain := []string{"d", "1"}
	for _, d := range doesNotContain {
		if set1.Contains(d) {
			t.Errorf("did not expect %s in the set", d)
		}
	}
}

func TestFilteredByKey(t *testing.T) {
	m := map[string]string{
		"a": "apple",
		"b": "boy",
		"c": "cat",
	}
	out := FilteredByKey(m, func(key string) bool {
		return key == "a" || key == "b"
	})
	if len(out) != 2 {
		t.Errorf("expected 2 entries")
	}
	if _, ok := out["c"]; ok {
		t.Errorf("deleted key %s was still found", "c")
	}
	if out["a"] != "apple" || out["b"] != "boy" {
		t.Errorf("map content is not as expected.")
	}
}
