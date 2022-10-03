// Helper functions to work with slices
// Eventually, this will be replaced by https://pkg.go.dev/golang.org/x/exp/slices but now it's not stable

package slices

func Find[T any](slice []T, predicate func(*T) bool) *T {
	for i, v := range slice {
		if predicate(&v) {
			return &slice[i]
		}
	}
	return nil
}

func FindPointer[T any](slice []*T, predicate func(*T) bool) *T {
	for i, v := range slice {
		if predicate(v) {
			return slice[i]
		}
	}
	return nil
}
