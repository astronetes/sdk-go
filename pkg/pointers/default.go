package pointers

func GetValueOrDefault[T any](p *T, defaultValue T) T {
	if p == nil {
		return defaultValue
	}
	return *p
}
