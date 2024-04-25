package util

func MapSlice[T any, U any](slice []T, fn func(T) U) []U {
	var ret []U
	if slice != nil {
		ret = make([]U, len(slice))
		for i, t := range slice {
			ret[i] = fn(t)
		}
	}
	return ret
}
