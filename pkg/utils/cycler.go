package utils

func Cycle[T comparable](current T, collection []T, isForward bool) T {
	length := len(collection)

	if length == 0 {
		var zero T
		return zero
	}

	for i, item := range collection {
		if item == current {
			if !isForward {
				return collection[(i-1+length)%length]
			}
			return collection[(i+1)%length]
		}
	}
	return collection[0]
}

func SafeIndex[T comparable](collection []T, idx int) T {
	length := len(collection)

	if length == 0 {
		var zero T
		return zero
	}

	return collection[idx]
}
