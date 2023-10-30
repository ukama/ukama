/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package archiver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/* Testing success case for unarchiving */
func Test_Unarchive(t *testing.T) {
	fname := "testdata/test.tar.gz"
	dest := "testdata"

	var want error = nil

	err := Unarchive(fname, dest)
	if want != err {
		t.Errorf("Expected Error '%s', but got '%s'", want, err.Error())
	}
}

/* Testing if source file msiising */
func Test_UnarchiveMissingFile(t *testing.T) {
	fname := "testdata/missingfile.tar.gz"
	dest := "testdata"

	err := Unarchive(fname, dest)

	assert.Contains(t, err.Error(), "invalid source")
}

/* Testing if source file msiising */
func Test_UnarchiveMissingDest(t *testing.T) {
	fname := "testdata/test.tar.gz"
	dest := "testdata/xyz"

	err := Unarchive(fname, dest)

	assert.Contains(t, err.Error(), "invalid destination")
}

/* Archive format failure */
func Test_UnarchiveInvalid(t *testing.T) {
	fname := "testdata/archiveformatfailure.txt"
	dest := "testdata"

	err := Unarchive(fname, dest)

	assert.Contains(t, err.Error(), "format unrecognized by filename")
}
