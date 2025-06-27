/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	ukama "github.com/ukama/ukama/systems/common/ukama"
	pb "github.com/ukama/ukama/systems/subscriber/sim-pool/pb/gen"
	"github.com/ukama/ukama/systems/subscriber/sim-pool/pkg/db"
)

func TestPoolStats(t *testing.T) {
	tests := []struct {
		name     string
		sims     []db.Sim
		expected *pb.GetStatsResponse
	}{
		{
			name: "Empty slice",
			sims: []db.Sim{},
			expected: &pb.GetStatsResponse{
				Total:     0,
				Failed:    0,
				Available: 0,
				Consumed:  0,
				Physical:  0,
				Esim:      0,
			},
		},
		{
			name: "All available physical sims",
			sims: []db.Sim{
				{Iccid: "1", IsAllocated: false, IsFailed: false, IsPhysical: true},
				{Iccid: "2", IsAllocated: false, IsFailed: false, IsPhysical: true},
				{Iccid: "3", IsAllocated: false, IsFailed: false, IsPhysical: true},
			},
			expected: &pb.GetStatsResponse{
				Total:     3,
				Failed:    0,
				Available: 3,
				Consumed:  0,
				Physical:  3,
				Esim:      0,
			},
		},
		{
			name: "All available esims",
			sims: []db.Sim{
				{Iccid: "1", IsAllocated: false, IsFailed: false, IsPhysical: false},
				{Iccid: "2", IsAllocated: false, IsFailed: false, IsPhysical: false},
			},
			expected: &pb.GetStatsResponse{
				Total:     2,
				Failed:    0,
				Available: 2,
				Consumed:  0,
				Physical:  0,
				Esim:      2,
			},
		},
		{
			name: "Mixed status sims",
			sims: []db.Sim{
				{Iccid: "1", IsAllocated: true, IsFailed: false, IsPhysical: true},   // consumed physical
				{Iccid: "2", IsAllocated: false, IsFailed: true, IsPhysical: true},   // failed physical
				{Iccid: "3", IsAllocated: false, IsFailed: false, IsPhysical: true},  // available physical
				{Iccid: "4", IsAllocated: true, IsFailed: false, IsPhysical: false},  // consumed esim
				{Iccid: "5", IsAllocated: false, IsFailed: true, IsPhysical: false},  // failed esim
				{Iccid: "6", IsAllocated: false, IsFailed: false, IsPhysical: false}, // available esim
			},
			expected: &pb.GetStatsResponse{
				Total:     6,
				Failed:    2,
				Available: 2,
				Consumed:  2,
				Physical:  1,
				Esim:      1,
			},
		},
		{
			name: "All consumed sims",
			sims: []db.Sim{
				{Iccid: "1", IsAllocated: true, IsFailed: false, IsPhysical: true},
				{Iccid: "2", IsAllocated: true, IsFailed: false, IsPhysical: false},
			},
			expected: &pb.GetStatsResponse{
				Total:     2,
				Failed:    0,
				Available: 0,
				Consumed:  2,
				Physical:  0,
				Esim:      0,
			},
		},
		{
			name: "All failed sims",
			sims: []db.Sim{
				{Iccid: "1", IsAllocated: false, IsFailed: true, IsPhysical: true},
				{Iccid: "2", IsAllocated: false, IsFailed: true, IsPhysical: false},
			},
			expected: &pb.GetStatsResponse{
				Total:     2,
				Failed:    2,
				Available: 0,
				Consumed:  0,
				Physical:  0,
				Esim:      0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PoolStats(tt.sims)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPbParseToModel(t *testing.T) {
	tests := []struct {
		name     string
		input    []*pb.AddSim
		expected []db.Sim
	}{
		{
			name:     "Empty slice",
			input:    []*pb.AddSim{},
			expected: []db.Sim(nil),
		},
		{
			name: "Single sim",
			input: []*pb.AddSim{
				{
					Iccid:          "12345678901234567890",
					Msisdn:         "+1234567890",
					SmDpAddress:    "smdp.example.com",
					ActivationCode: "ACT123",
					QrCode:         "QR123",
					SimType:        "test",
					IsPhysical:     true,
				},
			},
			expected: []db.Sim{
				{
					Iccid:          "12345678901234567890",
					Msisdn:         "+1234567890",
					SmDpAddress:    "smdp.example.com",
					ActivationCode: "ACT123",
					QrCode:         "QR123",
					SimType:        ukama.SimTypeTest,
					IsPhysical:     true,
				},
			},
		},
		{
			name: "Multiple sims with different types",
			input: []*pb.AddSim{
				{
					Iccid:          "11111111111111111111",
					Msisdn:         "+1111111111",
					SmDpAddress:    "smdp1.example.com",
					ActivationCode: "ACT111",
					QrCode:         "QR111",
					SimType:        "operator_data",
					IsPhysical:     true,
				},
				{
					Iccid:          "22222222222222222222",
					Msisdn:         "+2222222222",
					SmDpAddress:    "smdp2.example.com",
					ActivationCode: "ACT222",
					QrCode:         "QR222",
					SimType:        "ukama_data",
					IsPhysical:     false,
				},
			},
			expected: []db.Sim{
				{
					Iccid:          "11111111111111111111",
					Msisdn:         "+1111111111",
					SmDpAddress:    "smdp1.example.com",
					ActivationCode: "ACT111",
					QrCode:         "QR111",
					SimType:        ukama.SimTypeOperatorData,
					IsPhysical:     true,
				},
				{
					Iccid:          "22222222222222222222",
					Msisdn:         "+2222222222",
					SmDpAddress:    "smdp2.example.com",
					ActivationCode: "ACT222",
					QrCode:         "QR222",
					SimType:        ukama.SimTypeUkamaData,
					IsPhysical:     false,
				},
			},
		},
		{
			name: "Unknown sim type",
			input: []*pb.AddSim{
				{
					Iccid:          "33333333333333333333",
					Msisdn:         "+3333333333",
					SmDpAddress:    "smdp3.example.com",
					ActivationCode: "ACT333",
					QrCode:         "QR333",
					SimType:        "unknown_type",
					IsPhysical:     true,
				},
			},
			expected: []db.Sim{
				{
					Iccid:          "33333333333333333333",
					Msisdn:         "+3333333333",
					SmDpAddress:    "smdp3.example.com",
					ActivationCode: "ACT333",
					QrCode:         "QR333",
					SimType:        ukama.SimTypeUnknown,
					IsPhysical:     true,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PbParseToModel(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseBytesToRawSim(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		expected    []RawSim
		expectError bool
	}{
		{
			name:     "Empty CSV",
			input:    []byte(""),
			expected: []RawSim(nil),
		},
		{
			name: "Valid CSV with headers",
			input: []byte(`ICCID,MSISDN,SmDpAddress,ActivationCode,QrCode,IsPhysical
12345678901234567890,+1234567890,smdp.example.com,ACT123,QR123,TRUE
98765432109876543210,+9876543210,smdp2.example.com,ACT456,QR456,FALSE`),
			expected: []RawSim{
				{
					Iccid:          "12345678901234567890",
					Msisdn:         "+1234567890",
					SmDpAddress:    "smdp.example.com",
					ActivationCode: "ACT123",
					QrCode:         "QR123",
					IsPhysical:     "TRUE",
				},
				{
					Iccid:          "98765432109876543210",
					Msisdn:         "+9876543210",
					SmDpAddress:    "smdp2.example.com",
					ActivationCode: "ACT456",
					QrCode:         "QR456",
					IsPhysical:     "FALSE",
				},
			},
		},
		{
			name: "Single row CSV",
			input: []byte(`ICCID,MSISDN,SmDpAddress,ActivationCode,QrCode,IsPhysical
11111111111111111111,+1111111111,smdp1.example.com,ACT111,QR111,TRUE`),
			expected: []RawSim{
				{
					Iccid:          "11111111111111111111",
					Msisdn:         "+1111111111",
					SmDpAddress:    "smdp1.example.com",
					ActivationCode: "ACT111",
					QrCode:         "QR111",
					IsPhysical:     "TRUE",
				},
			},
		},
		{
			name: "CSV with empty fields",
			input: []byte(`ICCID,MSISDN,SmDpAddress,ActivationCode,QrCode,IsPhysical
12345678901234567890,,smdp.example.com,,QR123,TRUE`),
			expected: []RawSim{
				{
					Iccid:          "12345678901234567890",
					Msisdn:         "",
					SmDpAddress:    "smdp.example.com",
					ActivationCode: "",
					QrCode:         "QR123",
					IsPhysical:     "TRUE",
				},
			},
		},
		{
			name: "Invalid CSV format - wrong number of columns",
			input: []byte(`ICCID,MSISDN,SmDpAddress,ActivationCode,QrCode,IsPhysical
12345678901234567890,+1234567890,smdp.example.com,ACT123,QR123,TRUE,EXTRA_COLUMN`),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseBytesToRawSim(tt.input)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestRawSimToPb(t *testing.T) {
	tests := []struct {
		name     string
		input    []RawSim
		simType  string
		expected []db.Sim
	}{
		{
			name:     "Empty slice",
			input:    []RawSim{},
			simType:  "test",
			expected: []db.Sim(nil),
		},
		{
			name: "Single physical sim",
			input: []RawSim{
				{
					Iccid:          "12345678901234567890",
					Msisdn:         "+1234567890",
					SmDpAddress:    "smdp.example.com",
					ActivationCode: "ACT123",
					QrCode:         "QR123",
					IsPhysical:     "TRUE",
				},
			},
			simType: "test",
			expected: []db.Sim{
				{
					Iccid:          "12345678901234567890",
					Msisdn:         "+1234567890",
					SmDpAddress:    "smdp.example.com",
					ActivationCode: "ACT123",
					QrCode:         "QR123",
					SimType:        ukama.SimTypeTest,
					IsPhysical:     true,
				},
			},
		},
		{
			name: "Single esim",
			input: []RawSim{
				{
					Iccid:          "98765432109876543210",
					Msisdn:         "+9876543210",
					SmDpAddress:    "smdp2.example.com",
					ActivationCode: "ACT456",
					QrCode:         "QR456",
					IsPhysical:     "FALSE",
				},
			},
			simType: "operator_data",
			expected: []db.Sim{
				{
					Iccid:          "98765432109876543210",
					Msisdn:         "+9876543210",
					SmDpAddress:    "smdp2.example.com",
					ActivationCode: "ACT456",
					QrCode:         "QR456",
					SimType:        ukama.SimTypeOperatorData,
					IsPhysical:     false,
				},
			},
		},
		{
			name: "Multiple sims with different physical types",
			input: []RawSim{
				{
					Iccid:          "11111111111111111111",
					Msisdn:         "+1111111111",
					SmDpAddress:    "smdp1.example.com",
					ActivationCode: "ACT111",
					QrCode:         "QR111",
					IsPhysical:     "TRUE",
				},
				{
					Iccid:          "22222222222222222222",
					Msisdn:         "+2222222222",
					SmDpAddress:    "smdp2.example.com",
					ActivationCode: "ACT222",
					QrCode:         "QR222",
					IsPhysical:     "FALSE",
				},
			},
			simType: "ukama_data",
			expected: []db.Sim{
				{
					Iccid:          "11111111111111111111",
					Msisdn:         "+1111111111",
					SmDpAddress:    "smdp1.example.com",
					ActivationCode: "ACT111",
					QrCode:         "QR111",
					SimType:        ukama.SimTypeUkamaData,
					IsPhysical:     true,
				},
				{
					Iccid:          "22222222222222222222",
					Msisdn:         "+2222222222",
					SmDpAddress:    "smdp2.example.com",
					ActivationCode: "ACT222",
					QrCode:         "QR222",
					SimType:        ukama.SimTypeUkamaData,
					IsPhysical:     false,
				},
			},
		},
		{
			name: "Case insensitive physical flag",
			input: []RawSim{
				{
					Iccid:          "33333333333333333333",
					Msisdn:         "+3333333333",
					SmDpAddress:    "smdp3.example.com",
					ActivationCode: "ACT333",
					QrCode:         "QR333",
					IsPhysical:     "true",
				},
				{
					Iccid:          "44444444444444444444",
					Msisdn:         "+4444444444",
					SmDpAddress:    "smdp4.example.com",
					ActivationCode: "ACT444",
					QrCode:         "QR444",
					IsPhysical:     "false",
				},
			},
			simType: "test",
			expected: []db.Sim{
				{
					Iccid:          "33333333333333333333",
					Msisdn:         "+3333333333",
					SmDpAddress:    "smdp3.example.com",
					ActivationCode: "ACT333",
					QrCode:         "QR333",
					SimType:        ukama.SimTypeTest,
					IsPhysical:     false, // Only "TRUE" (uppercase) is considered true
				},
				{
					Iccid:          "44444444444444444444",
					Msisdn:         "+4444444444",
					SmDpAddress:    "smdp4.example.com",
					ActivationCode: "ACT444",
					QrCode:         "QR444",
					SimType:        ukama.SimTypeTest,
					IsPhysical:     false,
				},
			},
		},
		{
			name: "Unknown sim type",
			input: []RawSim{
				{
					Iccid:          "55555555555555555555",
					Msisdn:         "+5555555555",
					SmDpAddress:    "smdp5.example.com",
					ActivationCode: "ACT555",
					QrCode:         "QR555",
					IsPhysical:     "TRUE",
				},
			},
			simType: "unknown_type",
			expected: []db.Sim{
				{
					Iccid:          "55555555555555555555",
					Msisdn:         "+5555555555",
					SmDpAddress:    "smdp5.example.com",
					ActivationCode: "ACT555",
					QrCode:         "QR555",
					SimType:        ukama.SimTypeUnknown,
					IsPhysical:     true,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RawSimToPb(tt.input, tt.simType)
			assert.Equal(t, tt.expected, result)
		})
	}
}
