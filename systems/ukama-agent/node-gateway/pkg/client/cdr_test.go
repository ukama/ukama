package client

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	//"github.com/ukama/ukama/systems/init/lookup/gen/mocks"

	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	pb "github.com/ukama/ukama/systems/ukama-agent/cdr/pb/gen"
	amocks "github.com/ukama/ukama/systems/ukama-agent/cdr/pb/gen/mocks"
)

var cdr = &pb.CDR{
	Session:       1,
	NodeId:        ukama.NewVirtualHomeNodeId().String(),
	Imsi:          imsi,
	Policy:        uuid.NewV4().String(),
	ApnName:       "ukama.co",
	Ip:            "192.168.8.2",
	StartTime:     uint64(time.Now().Unix() - 100000),
	EndTime:       uint64(time.Now().Unix() - 50000),
	LastUpdatedAt: uint64(time.Now().Unix() - 50000),
	TxBytes:       2048000,
	RxBytes:       1024000,
	TotalBytes:    3072000,
}

func TestCDRClient_PostCDR(t *testing.T) {
	m := &amocks.CDRServiceClient{}
	l := &CDR{
		client: m,
	}

	m.On("PostCDR", mock.Anything, cdr).Return(&pb.CDRResp{}, nil)

	_, err := l.PostCDR(cdr)
	assert.NoError(t, err)
}

func TestCDRClient_GetCDR(t *testing.T) {
	m := &amocks.CDRServiceClient{}
	l := &CDR{
		client: m,
	}
	pReq := &pb.RecordReq{
		Imsi:      cdr.Imsi,
		StartTime: cdr.StartTime,
		EndTime:   cdr.EndTime,
	}

	m.On("GetCDR", mock.Anything, pReq).Return(&pb.RecordResp{Cdr: []*pb.CDR{cdr}}, nil)

	_, err := l.GetCDR(pReq)
	assert.NoError(t, err)
}

func TestAsrClient_GetUsage(t *testing.T) {
	m := &amocks.CDRServiceClient{}
	l := &CDR{
		client: m,
	}

	pReq := &pb.UsageReq{
		Imsi:      cdr.Imsi,
		StartTime: cdr.StartTime,
		EndTime:   cdr.EndTime,
	}

	m.On("GetUsage", mock.Anything, pReq).Return(&pb.UsageResp{Usage: cdr.TotalBytes}, nil)

	resp, err := l.GetUsage(pReq)
	if assert.NoError(t, err) {
		m.AssertExpectations(t)
		assert.Equal(t, cdr.TotalBytes, resp.Usage)
	}
}
