package collections

type Set[K comparable, T any] struct {
	values []T
	idx    map[K]bool
}

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

func (set *Set[K, T]) Contains(k K) bool {
	if _, ok := set.idx[k]; ok {
		return true
	}
	return false
}

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

func (set *Set[K, T]) AddAll(items []T, keyFun func(t T) K) {
	for _, item := range items {
		set.AddWithKeyFunc(item, keyFun)
	}
}

func (set *Set[K, T]) ToList() []T {
	return set.values
}

func Values[T any](m map[string]T) []T {
	list := make([]T, 0)
	for _, v := range m {
		list = append(list, v)
	}
	return list
}
