package util

// Ptr returns a pointer to valid
func Ptr[T any](v T) *T {
	return &v
}

// Val returns the value of the pointer, or the zero value if nil
func Val[T any](v *T) T {
	if v == nil {
		var zero T
		return zero
	}
	return *v
}

// ValOr returns the value of the pointer, or the default value if nil
func ValOr[T any](v *T, def T) T {
	if v == nil {
		return def
	}
	return *v
}
