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

const (
	// ICCID values for PoolStats tests
	TestIccidPool1 = "1"
	TestIccidPool2 = "2"
	TestIccidPool3 = "3"
	TestIccidPool4 = "4"
	TestIccidPool5 = "5"
	TestIccidPool6 = "6"

	// ICCID values for PbParseToModel tests
	TestIccidPb1 = "12345678901234567890"
	TestIccidPb2 = "11111111111111111111"
	TestIccidPb3 = "22222222222222222222"
	TestIccidPb4 = "33333333333333333333"

	// ICCID values for ParseBytesToRawSim tests
	TestIccidCsv1 = "12345678901234567890"
	TestIccidCsv2 = "98765432109876543210"
	TestIccidCsv3 = "11111111111111111111"

	// ICCID values for RawSimToPb tests
	TestIccidRaw1 = "12345678901234567890"
	TestIccidRaw2 = "98765432109876543210"
	TestIccidRaw3 = "11111111111111111111"
	TestIccidRaw4 = "22222222222222222222"
	TestIccidRaw5 = "33333333333333333333"
	TestIccidRaw6 = "44444444444444444444"
	TestIccidRaw7 = "55555555555555555555"

	// MSISDN values
	TestMsisdn1 = "+1234567890"
	TestMsisdn2 = "+9876543210"
	TestMsisdn3 = "+1111111111"
	TestMsisdn4 = "+2222222222"
	TestMsisdn5 = "+3333333333"
	TestMsisdn6 = "+4444444444"
	TestMsisdn7 = "+5555555555"

	// SmDpAddress values
	TestSmDpAddress1 = "smdp.example.com"
	TestSmDpAddress2 = "smdp2.example.com"
	TestSmDpAddress3 = "smdp1.example.com"
	TestSmDpAddress4 = "smdp3.example.com"
	TestSmDpAddress5 = "smdp4.example.com"
	TestSmDpAddress6 = "smdp5.example.com"

	// ActivationCode values
	TestActivationCode1 = "ACT123"
	TestActivationCode2 = "ACT456"
	TestActivationCode3 = "ACT111"
	TestActivationCode4 = "ACT222"
	TestActivationCode5 = "ACT333"
	TestActivationCode6 = "ACT444"
	TestActivationCode7 = "ACT555"

	// QR Code values
	TestQrCode1 = "QR123"
	TestQrCode2 = "QR456"
	TestQrCode3 = "QR111"
	TestQrCode4 = "QR222"
	TestQrCode5 = "QR333"
	TestQrCode6 = "QR444"
	TestQrCode7 = "QR555"

	// SimType string values
	TestSimTypeTest         = "test"
	TestSimTypeOperatorData = "operator_data"
	TestSimTypeUkamaData    = "ukama_data"
	TestSimTypeUnknown      = "unknown_type"

	// Physical flag values
	TestPhysicalTrue       = "TRUE"
	TestPhysicalFalse      = "FALSE"
	TestPhysicalTrueLower  = "true"
	TestPhysicalFalseLower = "false"

	// CSV headers
	TestCsvHeader = "ICCID,MSISDN,SmDpAddress,ActivationCode,QrCode,IsPhysical"
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
				{Iccid: TestIccidPool1, IsAllocated: false, IsFailed: false, IsPhysical: true},
				{Iccid: TestIccidPool2, IsAllocated: false, IsFailed: false, IsPhysical: true},
				{Iccid: TestIccidPool3, IsAllocated: false, IsFailed: false, IsPhysical: true},
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
				{Iccid: TestIccidPool1, IsAllocated: false, IsFailed: false, IsPhysical: false},
				{Iccid: TestIccidPool2, IsAllocated: false, IsFailed: false, IsPhysical: false},
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
				{Iccid: TestIccidPool1, IsAllocated: true, IsFailed: false, IsPhysical: true},   // consumed physical
				{Iccid: TestIccidPool2, IsAllocated: false, IsFailed: true, IsPhysical: true},   // failed physical
				{Iccid: TestIccidPool3, IsAllocated: false, IsFailed: false, IsPhysical: true},  // available physical
				{Iccid: TestIccidPool4, IsAllocated: true, IsFailed: false, IsPhysical: false},  // consumed esim
				{Iccid: TestIccidPool5, IsAllocated: false, IsFailed: true, IsPhysical: false},  // failed esim
				{Iccid: TestIccidPool6, IsAllocated: false, IsFailed: false, IsPhysical: false}, // available esim
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
				{Iccid: TestIccidPool1, IsAllocated: true, IsFailed: false, IsPhysical: true},
				{Iccid: TestIccidPool2, IsAllocated: true, IsFailed: false, IsPhysical: false},
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
				{Iccid: TestIccidPool1, IsAllocated: false, IsFailed: true, IsPhysical: true},
				{Iccid: TestIccidPool2, IsAllocated: false, IsFailed: true, IsPhysical: false},
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
					Iccid:          TestIccidPb1,
					Msisdn:         TestMsisdn1,
					SmDpAddress:    TestSmDpAddress1,
					ActivationCode: TestActivationCode1,
					QrCode:         TestQrCode1,
					SimType:        TestSimTypeTest,
					IsPhysical:     true,
				},
			},
			expected: []db.Sim{
				{
					Iccid:          TestIccidPb1,
					Msisdn:         TestMsisdn1,
					SmDpAddress:    TestSmDpAddress1,
					ActivationCode: TestActivationCode1,
					QrCode:         TestQrCode1,
					SimType:        ukama.SimTypeTest,
					IsPhysical:     true,
				},
			},
		},
		{
			name: "Multiple sims with different types",
			input: []*pb.AddSim{
				{
					Iccid:          TestIccidPb2,
					Msisdn:         TestMsisdn3,
					SmDpAddress:    TestSmDpAddress3,
					ActivationCode: TestActivationCode3,
					QrCode:         TestQrCode3,
					SimType:        TestSimTypeOperatorData,
					IsPhysical:     true,
				},
				{
					Iccid:          TestIccidPb3,
					Msisdn:         TestMsisdn4,
					SmDpAddress:    TestSmDpAddress2,
					ActivationCode: TestActivationCode4,
					QrCode:         TestQrCode4,
					SimType:        TestSimTypeUkamaData,
					IsPhysical:     false,
				},
			},
			expected: []db.Sim{
				{
					Iccid:          TestIccidPb2,
					Msisdn:         TestMsisdn3,
					SmDpAddress:    TestSmDpAddress3,
					ActivationCode: TestActivationCode3,
					QrCode:         TestQrCode3,
					SimType:        ukama.SimTypeOperatorData,
					IsPhysical:     true,
				},
				{
					Iccid:          TestIccidPb3,
					Msisdn:         TestMsisdn4,
					SmDpAddress:    TestSmDpAddress2,
					ActivationCode: TestActivationCode4,
					QrCode:         TestQrCode4,
					SimType:        ukama.SimTypeUkamaData,
					IsPhysical:     false,
				},
			},
		},
		{
			name: "Unknown sim type",
			input: []*pb.AddSim{
				{
					Iccid:          TestIccidPb4,
					Msisdn:         TestMsisdn5,
					SmDpAddress:    TestSmDpAddress4,
					ActivationCode: TestActivationCode5,
					QrCode:         TestQrCode5,
					SimType:        TestSimTypeUnknown,
					IsPhysical:     true,
				},
			},
			expected: []db.Sim{
				{
					Iccid:          TestIccidPb4,
					Msisdn:         TestMsisdn5,
					SmDpAddress:    TestSmDpAddress4,
					ActivationCode: TestActivationCode5,
					QrCode:         TestQrCode5,
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
			input: []byte(TestCsvHeader + "\n" +
				TestIccidCsv1 + "," + TestMsisdn1 + "," + TestSmDpAddress1 + "," + TestActivationCode1 + "," + TestQrCode1 + "," + TestPhysicalTrue + "\n" +
				TestIccidCsv2 + "," + TestMsisdn2 + "," + TestSmDpAddress2 + "," + TestActivationCode2 + "," + TestQrCode2 + "," + TestPhysicalFalse),
			expected: []RawSim{
				{
					Iccid:          TestIccidCsv1,
					Msisdn:         TestMsisdn1,
					SmDpAddress:    TestSmDpAddress1,
					ActivationCode: TestActivationCode1,
					QrCode:         TestQrCode1,
					IsPhysical:     TestPhysicalTrue,
				},
				{
					Iccid:          TestIccidCsv2,
					Msisdn:         TestMsisdn2,
					SmDpAddress:    TestSmDpAddress2,
					ActivationCode: TestActivationCode2,
					QrCode:         TestQrCode2,
					IsPhysical:     TestPhysicalFalse,
				},
			},
		},
		{
			name: "Single row CSV",
			input: []byte(TestCsvHeader + "\n" +
				TestIccidCsv3 + "," + TestMsisdn3 + "," + TestSmDpAddress3 + "," + TestActivationCode3 + "," + TestQrCode3 + "," + TestPhysicalTrue),
			expected: []RawSim{
				{
					Iccid:          TestIccidCsv3,
					Msisdn:         TestMsisdn3,
					SmDpAddress:    TestSmDpAddress3,
					ActivationCode: TestActivationCode3,
					QrCode:         TestQrCode3,
					IsPhysical:     TestPhysicalTrue,
				},
			},
		},
		{
			name: "CSV with empty fields",
			input: []byte(TestCsvHeader + "\n" +
				TestIccidCsv1 + ",," + TestSmDpAddress1 + ",," + TestQrCode1 + "," + TestPhysicalTrue),
			expected: []RawSim{
				{
					Iccid:          TestIccidCsv1,
					Msisdn:         "",
					SmDpAddress:    TestSmDpAddress1,
					ActivationCode: "",
					QrCode:         TestQrCode1,
					IsPhysical:     TestPhysicalTrue,
				},
			},
		},
		{
			name: "Invalid CSV format - wrong number of columns",
			input: []byte(TestCsvHeader + "\n" +
				TestIccidCsv1 + "," + TestMsisdn1 + "," + TestSmDpAddress1 + "," + TestActivationCode1 + "," + TestQrCode1 + "," + TestPhysicalTrue + ",EXTRA_COLUMN"),
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
					Iccid:          TestIccidRaw1,
					Msisdn:         TestMsisdn1,
					SmDpAddress:    TestSmDpAddress1,
					ActivationCode: TestActivationCode1,
					QrCode:         TestQrCode1,
					IsPhysical:     TestPhysicalTrue,
				},
			},
			simType: TestSimTypeTest,
			expected: []db.Sim{
				{
					Iccid:          TestIccidRaw1,
					Msisdn:         TestMsisdn1,
					SmDpAddress:    TestSmDpAddress1,
					ActivationCode: TestActivationCode1,
					QrCode:         TestQrCode1,
					SimType:        ukama.SimTypeTest,
					IsPhysical:     true,
				},
			},
		},
		{
			name: "Single esim",
			input: []RawSim{
				{
					Iccid:          TestIccidRaw2,
					Msisdn:         TestMsisdn2,
					SmDpAddress:    TestSmDpAddress2,
					ActivationCode: TestActivationCode2,
					QrCode:         TestQrCode2,
					IsPhysical:     TestPhysicalFalse,
				},
			},
			simType: TestSimTypeOperatorData,
			expected: []db.Sim{
				{
					Iccid:          TestIccidRaw2,
					Msisdn:         TestMsisdn2,
					SmDpAddress:    TestSmDpAddress2,
					ActivationCode: TestActivationCode2,
					QrCode:         TestQrCode2,
					SimType:        ukama.SimTypeOperatorData,
					IsPhysical:     false,
				},
			},
		},
		{
			name: "Multiple sims with different physical types",
			input: []RawSim{
				{
					Iccid:          TestIccidRaw3,
					Msisdn:         TestMsisdn3,
					SmDpAddress:    TestSmDpAddress3,
					ActivationCode: TestActivationCode3,
					QrCode:         TestQrCode3,
					IsPhysical:     TestPhysicalTrue,
				},
				{
					Iccid:          TestIccidRaw4,
					Msisdn:         TestMsisdn4,
					SmDpAddress:    TestSmDpAddress2,
					ActivationCode: TestActivationCode4,
					QrCode:         TestQrCode4,
					IsPhysical:     TestPhysicalFalse,
				},
			},
			simType: TestSimTypeUkamaData,
			expected: []db.Sim{
				{
					Iccid:          TestIccidRaw3,
					Msisdn:         TestMsisdn3,
					SmDpAddress:    TestSmDpAddress3,
					ActivationCode: TestActivationCode3,
					QrCode:         TestQrCode3,
					SimType:        ukama.SimTypeUkamaData,
					IsPhysical:     true,
				},
				{
					Iccid:          TestIccidRaw4,
					Msisdn:         TestMsisdn4,
					SmDpAddress:    TestSmDpAddress2,
					ActivationCode: TestActivationCode4,
					QrCode:         TestQrCode4,
					SimType:        ukama.SimTypeUkamaData,
					IsPhysical:     false,
				},
			},
		},
		{
			name: "Case insensitive physical flag",
			input: []RawSim{
				{
					Iccid:          TestIccidRaw5,
					Msisdn:         TestMsisdn5,
					SmDpAddress:    TestSmDpAddress4,
					ActivationCode: TestActivationCode5,
					QrCode:         TestQrCode5,
					IsPhysical:     TestPhysicalTrueLower,
				},
				{
					Iccid:          TestIccidRaw6,
					Msisdn:         TestMsisdn6,
					SmDpAddress:    TestSmDpAddress5,
					ActivationCode: TestActivationCode6,
					QrCode:         TestQrCode6,
					IsPhysical:     TestPhysicalFalseLower,
				},
			},
			simType: TestSimTypeTest,
			expected: []db.Sim{
				{
					Iccid:          TestIccidRaw5,
					Msisdn:         TestMsisdn5,
					SmDpAddress:    TestSmDpAddress4,
					ActivationCode: TestActivationCode5,
					QrCode:         TestQrCode5,
					SimType:        ukama.SimTypeTest,
					IsPhysical:     false, // Only "TRUE" (uppercase) is considered true
				},
				{
					Iccid:          TestIccidRaw6,
					Msisdn:         TestMsisdn6,
					SmDpAddress:    TestSmDpAddress5,
					ActivationCode: TestActivationCode6,
					QrCode:         TestQrCode6,
					SimType:        ukama.SimTypeTest,
					IsPhysical:     false,
				},
			},
		},
		{
			name: "Unknown sim type",
			input: []RawSim{
				{
					Iccid:          TestIccidRaw7,
					Msisdn:         TestMsisdn7,
					SmDpAddress:    TestSmDpAddress6,
					ActivationCode: TestActivationCode7,
					QrCode:         TestQrCode7,
					IsPhysical:     TestPhysicalTrue,
				},
			},
			simType: TestSimTypeUnknown,
			expected: []db.Sim{
				{
					Iccid:          TestIccidRaw7,
					Msisdn:         TestMsisdn7,
					SmDpAddress:    TestSmDpAddress6,
					ActivationCode: TestActivationCode7,
					QrCode:         TestQrCode7,
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
