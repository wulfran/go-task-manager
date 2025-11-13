package helpers

func SliceContains[T comparable](slice []T, value T) bool {
	for _, val := range slice {
		if val == value {
			return true
		}
	}
	return false
}
