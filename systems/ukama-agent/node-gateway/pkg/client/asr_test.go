/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package client

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	pb "github.com/ukama/ukama/systems/ukama-agent/asr/pb/gen"
	amocks "github.com/ukama/ukama/systems/ukama-agent/asr/pb/gen/mocks"
)

var iccid = "012345678901234567891"
var imsi = "012345678912345"

// var orgId = "880f7c63-eb57-461a-b514-248ce91e9b3e"
var packageId = "8adcdfb4-ed30-405d-b32f-d0b2dda4a1e0"

func TestAsrClient_UpdateGuti(t *testing.T) {
	m := &amocks.AsrRecordServiceClient{}
	l := &Asr{
		client: m,
	}

	pReq := &pb.UpdateGutiReq{
		Imsi:      imsi,
		UpdatedAt: uint32(time.Now().Unix()),
		Guti: &pb.Guti{
			PlmnId: "00101",
			Mmegi:  3200,
			Mmec:   100,
			Mtmsi:  1,
		},
	}

	m.On("UpdateGuti", mock.Anything, pReq).Return(&pb.UpdateGutiResp{}, nil)

	_, err := l.UpdateGuti(pReq)
	assert.NoError(t, err)
}

func TestAsrClient_UpdateTai(t *testing.T) {
	m := &amocks.AsrRecordServiceClient{}

	l := &Asr{
		client: m,
	}

	pReq := &pb.UpdateTaiReq{
		Imsi:      imsi,
		UpdatedAt: uint32(time.Now().Unix()),
		Tac:       1,
	}

	m.On("UpdateTai", mock.Anything, pReq).Return(&pb.UpdateTaiResp{}, nil)

	_, err := l.UpdateTai(pReq)
	assert.NoError(t, err)
}

func TestAsrClient_Read(t *testing.T) {
	m := &amocks.AsrRecordServiceClient{}

	l := &Asr{
		client: m,
	}

	pReq := &pb.ReadReq{
		Id: &pb.ReadReq_Imsi{
			Imsi: imsi,
		},
	}

	pResp := &pb.ReadResp{
		Record: &pb.Record{
			Iccid:        iccid,
			SimPackageId: "880f7c63-eb57-461a-b514-248ce91e9b3e",
			Imsi:         imsi,
			Op:           []byte("0123456789012345"),
			Key:          []byte("0123456789012345"),
			Amf:          []byte("800"),
			AlgoType:     1,
			UeDlAmbrBps:  2000000,
			UeUlAmbrBps:  2000000,
			Sqn:          1,
			CsgIdPrsent:  false,
			CsgId:        0,
			PackageId:    packageId,
		},
	}

	m.On("Read", mock.Anything, pReq).Return(pResp, nil)

	resp, err := l.Read(pReq)
	if assert.NoError(t, err) {
		m.AssertExpectations(t)
		assert.Equal(t, imsi, resp.Record.Imsi)
	}
}
