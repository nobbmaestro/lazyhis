package utils

func CycleNext[T comparable](current T, collection []T) T {
	if len(collection) == 0 {
		var zero T
		return zero
	}
	for i, item := range collection {
		if item == current {
			return collection[(i+1)%len(collection)]
		}
	}
	return collection[0]
}
