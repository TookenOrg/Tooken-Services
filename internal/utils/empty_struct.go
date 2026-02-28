package utils

func IsEmptyStruct[T comparable](v T) bool {
	var zero T
	return v == zero
}

func IsNilOrEmptyStruct[T comparable](v *T) bool {
	if v == nil {
		return true
	}
	var zero T
	return *v == zero
}
