package collections

// Keys {1:a, 2:b} -> [1, 2]
func Keys[K comparable, V any](m map[K]V) []K {
	if len(m) == 0 {
		return nil
	}

	i := 0
	keys := make([]K, len(m))
	for key := range m {
		keys[i] = key
		i++
	}

	return keys
}

// Values {1:a, 2:b} -> [a, b]
func Values[K comparable, V any](m map[K]V) []V {
	if len(m) == 0 {
		return nil
	}

	i := 0
	values := make([]V, len(m))
	for _, v := range m {
		values[i] = v
		i++
	}

	return values
}
