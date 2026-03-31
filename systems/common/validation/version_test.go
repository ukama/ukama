/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const invalidVersionErrMsg = "Invalid version format. Refer to https://semver.org/ for more information"

func TestParseVersion(t *testing.T) {
	tests := []struct {
		name        string
		version     string
		wantVersion string
		wantErr     bool
		checkErr    func(t *testing.T, err error)
	}{
		{
			name:        "valid semver major minor patch",
			version:     "1.0.0",
			wantVersion: "1.0.0",
			wantErr:     false,
		},
		{
			name:        "valid semver with prerelease",
			version:     "2.3.4-beta",
			wantVersion: "2.3.4-beta",
			wantErr:     false,
		},
		{
			name:        "valid semver with build metadata",
			version:     "1.0.0+20230301",
			wantVersion: "1.0.0+20230301",
			wantErr:     false,
		},
		{
			name:    "empty string",
			version: "",
			wantErr: true,
			checkErr: func(t *testing.T, err error) {
				assert.Equal(t, invalidVersionErrMsg, err.Error())
			},
		},
		{
			name:    "invalid format missing patch",
			version: "1.0",
			wantErr: true,
			checkErr: func(t *testing.T, err error) {
				assert.Equal(t, invalidVersionErrMsg, err.Error())
			},
		},
		{
			name:    "invalid format non-numeric",
			version: "invalid",
			wantErr: true,
			checkErr: func(t *testing.T, err error) {
				assert.Equal(t, invalidVersionErrMsg, err.Error())
			},
		},
		{
			name:    "invalid format with spaces",
			version: "1.0.0 ",
			wantErr: true,
			checkErr: func(t *testing.T, err error) {
				assert.Equal(t, invalidVersionErrMsg, err.Error())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseVersion(tt.version)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
				if tt.checkErr != nil {
					tt.checkErr(t, err)
				}
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, tt.wantVersion, got.String())
		})
	}
}

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		name     string
		version1 string
		version2 string
		want     int
		wantErr  bool
		checkErr func(t *testing.T, err error)
	}{
		{
			name:     "version1 greater than version2",
			version1: "2.0.0",
			version2: "1.0.0",
			want:     1,
		},
		{
			name:     "version1 less than version2",
			version1: "1.0.0",
			version2: "2.0.0",
			want:     -1,
		},
		{
			name:     "version1 equal to version2",
			version1: "1.0.0",
			version2: "1.0.0",
			want:     0,
		},
		{
			name:     "minor version greater",
			version1: "1.2.0",
			version2: "1.1.0",
			want:     1,
		},
		{
			name:     "patch version greater",
			version1: "1.0.1",
			version2: "1.0.0",
			want:     1,
		},
		{
			name:     "prerelease less than release",
			version1: "1.0.0-beta",
			version2: "1.0.0",
			want:     -1,
		},
		{
			name:     "release greater than prerelease",
			version1: "1.0.0",
			version2: "1.0.0-beta",
			want:     1,
		},
		{
			name:     "invalid version1 returns error",
			version1: "invalid",
			version2: "1.0.0",
			wantErr:  true,
			checkErr: func(t *testing.T, err error) {
				assert.Equal(t, invalidVersionErrMsg, err.Error())
			},
		},
		{
			name:     "invalid version2 returns error",
			version1: "1.0.0",
			version2: "invalid",
			wantErr:  true,
			checkErr: func(t *testing.T, err error) {
				assert.Equal(t, invalidVersionErrMsg, err.Error())
			},
		},
		{
			name:     "empty version1 returns error",
			version1: "",
			version2: "1.0.0",
			wantErr:  true,
			checkErr: func(t *testing.T, err error) {
				assert.Equal(t, invalidVersionErrMsg, err.Error())
			},
		},
		{
			name:     "empty version2 returns error",
			version1: "1.0.0",
			version2: "",
			wantErr:  true,
			checkErr: func(t *testing.T, err error) {
				assert.Equal(t, invalidVersionErrMsg, err.Error())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CompareVersions(tt.version1, tt.version2)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, 0, got)
				if tt.checkErr != nil {
					tt.checkErr(t, err)
				}
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
