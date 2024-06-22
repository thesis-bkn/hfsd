package utils

import "github.com/ztrue/tracerr"

func MapToSlices[M ~map[K]V, K comparable, V any](m M) []V {
	r := make([]V, 0, len(m))
	for _, v := range m {
		r = append(r, v)
	}
	return r
}

func SlicesToMap[M map[K]V, K comparable, V any](s []V, fn func(V) K) map[K]V {
	m := make(map[K]V)
	for _, v := range s {
		m[fn(v)] = v
	}
	return m
}

func Map[T, V any](ts []T, fn func(T) V) []V {
	result := make([]V, len(ts))
	for i, t := range ts {
		result[i] = fn(t)
	}
	return result
}

func MapErr[T, V any](ts []T, fn func(T) (V, error)) ([]V, error) {
	result := make([]V, len(ts))
	var err error
	for i, t := range ts {
		result[i], err = fn(t)
		if err != nil {
			return nil, tracerr.Wrap(err)
		}
	}
	return result, nil
}

func Sum[T any, V int | float32 | float64](s []T, fn func(T) V) V {
	var sum V
	for _, v := range s {
		sum += fn(v)
	}
	return sum
}

func Group[K comparable, V any](s []V, keyfn func(a V) K, group func(a, b V) V) []V {
	m := make(map[K]V)
	for _, x := range s {
		k := keyfn(x)
		if y, ok := m[k]; !ok {
			m[k] = x
		} else {
			m[k] = group(x, y)
		}
	}

	return MapToSlices(m)
}
