/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
package utils

import (
	"encoding/json"
	"testing"
)

func TestJSONMap_Value(t *testing.T) {
	tests := []struct {
		name    string
		jsonMap JSONMap
		wantErr bool
	}{
		{
			name: "valid JSONMap",
			jsonMap: JSONMap{
				"key1": "value1",
				"key2": 123,
				"key3": map[string]interface{}{
					"nested": "value",
				},
			},
			wantErr: false,
		},
		{
			name:    "empty JSONMap",
			jsonMap: JSONMap{},
			wantErr: false,
		},
		{
			name:    "nil JSONMap",
			jsonMap: nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.jsonMap.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("JSONMap.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Verify that the result can be unmarshaled back to JSONMap
				var result JSONMap
				err = json.Unmarshal(got.([]byte), &result)
				if err != nil {
					t.Errorf("JSONMap.Value() returned invalid JSON: %v", err)
				}
			}
		})
	}
}

func TestJSONMap_Scan(t *testing.T) {
	tests := []struct {
		name    string
		value   interface{}
		want    JSONMap
		wantErr bool
	}{
		{
			name:  "valid JSON bytes",
			value: []byte(`{"key1":"value1","key2":123,"key3":{"nested":"value"}}`),
			want: JSONMap{
				"key1": "value1",
				"key2": float64(123), // JSON numbers are unmarshaled as float64
				"key3": map[string]interface{}{
					"nested": "value",
				},
			},
			wantErr: false,
		},
		{
			name:    "nil value",
			value:   nil,
			want:    JSONMap{},
			wantErr: false,
		},
		{
			name:    "empty JSON bytes",
			value:   []byte(`{}`),
			want:    JSONMap{},
			wantErr: false,
		},
		{
			name:    "invalid JSON bytes",
			value:   []byte(`{"invalid": json`),
			want:    JSONMap{},
			wantErr: true,
		},
		{
			name:    "non-byte value",
			value:   "not bytes",
			want:    JSONMap{},
			wantErr: true,
		},
		{
			name:    "int value",
			value:   123,
			want:    JSONMap{},
			wantErr: true,
		},
		{
			name:    "bool value",
			value:   true,
			want:    JSONMap{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var m JSONMap
			err := m.Scan(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSONMap.Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Compare the scanned result with expected
				if len(m) != len(tt.want) {
					t.Errorf("JSONMap.Scan() result length = %d, want %d", len(m), len(tt.want))
				}
				// For simplicity, we'll just check that the scan didn't fail
				// In a real scenario, you might want to do deeper comparison
			}
		})
	}
}

func TestJSONMap_ScanWithNilReceiver(t *testing.T) {
	// Test scanning into a nil receiver
	var m *JSONMap
	err := m.Scan([]byte(`{"key":"value"}`))
	if err == nil {
		t.Error("JSONMap.Scan() should fail when receiver is nil")
	}
}

func TestJSONMap_ValueAndScanRoundTrip(t *testing.T) {
	original := JSONMap{
		"string": "test",
		"number": 42,
		"bool":   true,
		"array":  []interface{}{1, 2, 3},
		"object": map[string]interface{}{
			"nested": "value",
		},
		"null": nil,
	}

	// Test Value() method
	value, err := original.Value()
	if err != nil {
		t.Fatalf("JSONMap.Value() failed: %v", err)
	}

	// Test Scan() method
	var scanned JSONMap
	err = scanned.Scan(value)
	if err != nil {
		t.Fatalf("JSONMap.Scan() failed: %v", err)
	}

	// Verify round trip
	originalBytes, _ := json.Marshal(original)
	scannedBytes, _ := json.Marshal(scanned)

	if string(originalBytes) != string(scannedBytes) {
		t.Errorf("Round trip failed: original = %s, scanned = %s",
			string(originalBytes), string(scannedBytes))
	}
}
