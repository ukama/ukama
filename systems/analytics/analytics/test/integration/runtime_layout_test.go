/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package integration

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRuntimeLayoutMatchesArchitecture(t *testing.T) {
	root := filepath.Clean("../..")

	mustExist := []string{
		"analytics/cmd/server/main.go",
		"api-gateway/cmd/main.go",
		"collector/cmd/server/main.go",
		"business/pkg",
		"customer/pkg",
		"network/pkg",
		"business/pb",
		"customer/pb",
		"network/pb",
	}

	mustNotExist := []string{
		"business/cmd",
		"customer/cmd",
		"network/cmd",
		"business/Dockerfile",
		"customer/Dockerfile",
		"network/Dockerfile",
		"business/Int.Dockerfile",
		"customer/Int.Dockerfile",
		"network/Int.Dockerfile",
	}

	for _, path := range mustExist {
		if _, err := os.Stat(filepath.Join(root, path)); err != nil {
			t.Fatalf("expected %s to exist: %v", path, err)
		}
	}

	for _, path := range mustNotExist {
		if _, err := os.Stat(filepath.Join(root, path)); err == nil {
			t.Fatalf("expected old runtime artifact %s to be removed", path)
		} else if !os.IsNotExist(err) {
			t.Fatalf("failed to stat %s: %v", path, err)
		}
	}
}
