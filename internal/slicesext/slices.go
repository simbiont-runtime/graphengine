// ---

package slicesext

import (
	"golang.org/x/exp/slices"
)

func ContainsFunc[S ~[]E, E any](s S, f func(E) bool) bool {
	return slices.IndexFunc(s, f) >= 0
}

func FilterFunc[S ~[]E, E any](s S, f func(E) bool) S {
	n := 0
	for i, v := range s {
		if f(v) {
			s[n] = s[i]
			n++
		}
	}
	return s[:n]
}

func FindFunc[S ~[]E, E any](s S, f func(E) bool) (e E, ok bool) {
	for _, v := range s {
		if f(v) {
			return v, true
		}
	}
	return
}
