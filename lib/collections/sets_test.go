package collections

import (
	"reflect"
	"testing"
)

func run[A any](input map[string]A, expected int) func(t *testing.T) {
	return func(t *testing.T) {
		result := Values(input)
		if len(result) != expected {
			t.Errorf("empty map expected %d but found %d", expected, len(result))
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

func TestNewSetWithDuplicates(t *testing.T) {
	set := NewSet[string]([]string{"a", "b", "a", "c", "b"})
	expected := []string{"a", "b", "c"}
	result := set.ToList()

	if len(result) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(result))
	}

	for _, v := range expected {
		if !set.Contains(v) {
			t.Errorf("expected value %s not found in set", v)
		}
	}
}

func TestNewSetWithComparator(t *testing.T) {
	type Person struct {
		ID   int
		Name string
	}

	people := []Person{
		{1, "Alice"},
		{2, "Bob"},
		{1, "Alice Clone"}, // Same ID as Alice
		{3, "Charlie"},
	}

	set := NewSetWithComparator[int, Person](people, func(p Person) int {
		return p.ID
	})

	// Should only contain 3 items (duplicate ID=1 should be removed)
	result := set.ToList()
	if len(result) != 3 {
		t.Errorf("expected 3 items, got %d", len(result))
	}

	// Check if IDs 1, 2, 3 are in the set
	for _, id := range []int{1, 2, 3} {
		if !set.Contains(id) {
			t.Errorf("expected ID %d to be in set", id)
		}
	}
}

func TestAddWithKeyFunc(t *testing.T) {
	type Person struct {
		ID   int
		Name string
	}

	set := NewSetWithComparator[int, Person](nil, func(p Person) int {
		return p.ID
	})

	// Add new items
	set.AddWithKeyFunc(Person{1, "Alice"}, func(p Person) int { return p.ID })
	set.AddWithKeyFunc(Person{2, "Bob"}, func(p Person) int { return p.ID })

	// Try to add duplicate (same ID)
	set.AddWithKeyFunc(Person{1, "Alice Clone"}, func(p Person) int { return p.ID })

	result := set.ToList()
	if len(result) != 2 {
		t.Errorf("expected 2 items, got %d", len(result))
	}

	if !set.Contains(1) || !set.Contains(2) {
		t.Error("set is missing expected IDs")
	}
}

func TestAddAll(t *testing.T) {
	set := NewSetWithComparator[string, string](nil, StringKey)

	// Add initial items
	initialItems := []string{"a", "b"}
	set.AddAll(initialItems, StringKey)

	// Add more items including duplicates
	moreItems := []string{"b", "c", "d"}
	set.AddAll(moreItems, StringKey)

	result := set.ToList()
	expected := []string{"a", "b", "c", "d"}

	if len(result) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(result))
	}

	for _, v := range expected {
		if !set.Contains(v) {
			t.Errorf("expected value %s not found in set", v)
		}
	}
}

func TestToList(t *testing.T) {
	// Test empty set
	emptySet := NewSet[string](nil)
	emptyResult := emptySet.ToList()
	if len(emptyResult) != 0 {
		t.Errorf("expected empty list, got length %d", len(emptyResult))
	}

	// Test non-empty set
	set := NewSet[string]([]string{"a", "b", "c"})
	result := set.ToList()
	expected := []string{"a", "b", "c"}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestStringKey(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"test", "test"},
		{"", ""},
		{"hello world", "hello world"},
	}

	for _, tc := range testCases {
		result := StringKey(tc.input)
		if result != tc.expected {
			t.Errorf("for input %q, expected %q, got %q", tc.input, tc.expected, result)
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
