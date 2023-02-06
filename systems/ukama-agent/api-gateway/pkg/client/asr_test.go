package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	//"github.com/ukama/ukama/systems/init/lookup/gen/mocks"

	pb "github.com/ukama/ukama/systems/ukama-agent/asr/pb/gen"
	amocks "github.com/ukama/ukama/systems/ukama-agent/asr/pb/gen/mocks"
)

var iccid = "012345678901234567891"
var network = "40987edb-ebb6-4f84-a27c-99db7c136127"

// var orgId = "880f7c63-eb57-461a-b514-248ce91e9b3e"
var packageId = "8adcdfb4-ed30-405d-b32f-d0b2dda4a1e0"

var sub = pb.ReadResp{
	Record: &pb.Record{
		Iccid:       iccid,
		SimId:       "880f7c63-eb57-461a-b514-248ce91e9b3e",
		Imsi:        "012345678912345",
		Op:          []byte("0123456789012345"),
		Key:         []byte("0123456789012345"),
		Amf:         []byte("800"),
		AlgoType:    1,
		UeDlAmbrBps: 2000000,
		UeUlAmbrBps: 2000000,
		Sqn:         1,
		CsgIdPrsent: false,
		CsgId:       0,
		PackageId:   packageId,
	},
}

func TestAsrClient_Activate(t *testing.T) {
	m := &amocks.AsrRecordServiceClient{}
	l := &Asr{
		client: m,
	}

	pReq := &pb.ActivateReq{
		Iccid:     iccid,
		Network:   network,
		PackageId: packageId,
	}

	m.On("Activate", mock.Anything, pReq).Return(&pb.ActivateResp{}, nil)

	_, err := l.Activate(pReq)
	assert.NoError(t, err)
}

func TestAsrClient_UpdatePackage(t *testing.T) {
	m := &amocks.AsrRecordServiceClient{}

	l := &Asr{
		client: m,
	}

	pReq := &pb.UpdatePackageReq{
		Iccid:     iccid,
		PackageId: packageId,
	}

	m.On("UpdatePackage", mock.Anything, pReq).Return(&pb.UpdatePackageResp{}, nil)

	_, err := l.UpdatePackage(pReq)
	assert.NoError(t, err)
}

func TestAsrClient_Read(t *testing.T) {
	m := &amocks.AsrRecordServiceClient{}

	l := &Asr{
		client: m,
	}

	pReq := &pb.ReadReq{
		Id: &pb.ReadReq_Iccid{
			Iccid: iccid,
		},
	}

	pResp := &pb.ReadResp{
		Record: &pb.Record{
			Iccid:       iccid,
			SimId:       "880f7c63-eb57-461a-b514-248ce91e9b3e",
			Imsi:        "012345678912345",
			Op:          []byte("0123456789012345"),
			Key:         []byte("0123456789012345"),
			Amf:         []byte("800"),
			AlgoType:    1,
			UeDlAmbrBps: 2000000,
			UeUlAmbrBps: 2000000,
			Sqn:         1,
			CsgIdPrsent: false,
			CsgId:       0,
			PackageId:   packageId,
		},
	}

	m.On("Read", mock.Anything, pReq).Return(pResp, nil)

	resp, err := l.Read(pReq)
	if assert.NoError(t, err) {
		m.AssertExpectations(t)
		assert.Equal(t, iccid, resp.Record.Iccid)
	}
}

func TestAsrClient_Inactivate(t *testing.T) {
	m := &amocks.AsrRecordServiceClient{}

	l := &Asr{
		client: m,
	}

	pReq := &pb.InactivateReq{
		Id: &pb.InactivateReq_Iccid{
			Iccid: iccid,
		},
	}

	m.On("Inactivate", mock.Anything, pReq).Return(&pb.InactivateResp{}, nil)

	_, err := l.Inactivate(pReq)
	assert.NoError(t, err)
}
