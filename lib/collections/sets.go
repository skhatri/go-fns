// Package collections provides utility functions for working with Go collections like maps and sets.
package collections

// Set represents a generic set data structure that can store unique values.
// It supports operations like adding, checking for membership, and converting to a list.
type Set[K comparable, T any] struct {
	values []T
	idx    map[K]bool
}

// NewSet creates a new empty Set with the specified type parameter.
// The type parameter T must be comparable to ensure proper set operations.
func NewSet[T comparable](items []T) Set[T, T] {
	idx := make(map[T]bool)
	dedup := make([]T, 0)
	for _, item := range items {
		if _, ok := idx[item]; !ok {
			idx[item] = true
			dedup = append(dedup, item)
		}
	}
	return Set[T, T]{
		values: dedup,
		idx:    idx,
	}
}

// NewSetWithComparator creates a new Set with a custom comparison function.
// This is useful when working with types that don't implement the comparable interface
// or when you need custom comparison logic.
func NewSetWithComparator[K comparable, T any](items []T, keyFun func(T) K) Set[K, T] {
	idx := make(map[K]bool)
	dedup := make([]T, 0)
	for _, item := range items {
		key := keyFun(item)
		if _, ok := idx[key]; !ok {
			idx[key] = true
			dedup = append(dedup, item)
		}
	}
	return Set[K, T]{
		values: dedup,
		idx:    idx,
	}
}

// Contains checks if the set contains the specified value.
// Returns true if the value is present in the set, false otherwise.
func (set *Set[K, T]) Contains(k K) bool {
	if _, ok := set.idx[k]; ok {
		return true
	}
	return false
}

// AddWithKeyFunc adds a value to the set using a custom key function.
// The key function is used to determine uniqueness of values in the set.
// This is useful when you want to consider certain fields of a struct for uniqueness.
func (set *Set[K, T]) AddWithKeyFunc(item T, keyFunc func(T) K) {
	key := keyFunc(item)
	if _, ok := set.idx[key]; !ok {
		set.idx[key] = true
		set.values = append(set.values, item)
	}
}

var StringKey = func(s string) string {
	return s
}

// AddAll adds all values from the provided slice to the set.
// Duplicate values will be automatically handled by the set's uniqueness constraint.
func (set *Set[K, T]) AddAll(items []T, keyFun func(t T) K) {
	for _, item := range items {
		set.AddWithKeyFunc(item, keyFun)
	}
}

// ToList converts the set to a slice containing all unique values.
// The order of elements in the resulting slice is not guaranteed.
func (set *Set[K, T]) ToList() []T {
	return set.values
}

// Values returns a slice containing all values in the set.
// This is an alias for ToList() for better readability.
func Values[T any](m map[string]T) []T {
	list := make([]T, 0)
	for _, v := range m {
		list = append(list, v)
	}
	return list
}
