/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

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
