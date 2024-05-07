package server

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/ukama"

	cmocks "github.com/ukama/ukama/systems/common/mocks"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/common/uuid"
	mocks "github.com/ukama/ukama/systems/ukama-agent/cdr/mocks"
	pb "github.com/ukama/ukama/systems/ukama-agent/cdr/pb/gen"
	"github.com/ukama/ukama/systems/ukama-agent/cdr/pkg/db"
)

var cdr = db.CDR{
	Session:       1,
	NodeId:        ukama.NewVirtualHomeNodeId().String(),
	Imsi:          "123456789012345678",
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

var usage = db.Usage{
	Imsi:             "123456789012345678",
	Historical:       0,
	Usage:            0,
	LastSessionUsage: 0,
	LastSessionId:    0,
}

var OrgName = "ukama"
var OrgId = "40987edb-ebb6-4f84-a27c-99db7c136127"

func TestCDR_PostCDR(t *testing.T) {
	cdrRepo := &mocks.CDRRepo{}
	usageRepo := &mocks.UsageRepo{}
	mbC := &cmocks.MsgBusServiceClient{}
	s, err := NewCDRServer(cdrRepo, usageRepo, OrgId, OrgName, mbC)
	assert.NoError(t, err)

	req := &pb.CDR{
		Session:       cdr.Session,
		NodeId:        cdr.NodeId,
		Imsi:          cdr.Imsi,
		Policy:        cdr.Policy,
		ApnName:       cdr.ApnName,
		Ip:            cdr.Ip,
		StartTime:     cdr.StartTime,
		EndTime:       cdr.EndTime,
		LastUpdatedAt: cdr.LastUpdatedAt,
		TxBytes:       cdr.TxBytes,
		RxBytes:       cdr.RxBytes,
		TotalBytes:    cdr.TotalBytes,
	}

	cdrRepo.On("Add", &cdr).Return(nil).Once()
	usageRepo.On("Get", cdr.Imsi).Return(&usage, nil).Once()
	cdrRepo.On("GetByTimeAndNodeId", cdr.Imsi, cdr.StartTime, mock.Anything, cdr.NodeId).Return(&[]db.CDR{cdr}, nil).Once()
	usageRepo.On("Add", mock.MatchedBy(func(u *db.Usage) bool {
		return u.Imsi == cdr.Imsi
	})).Return(nil).Once()
	mbC.On("PublishRequest", "event.cloud.local.ukama.ukamaagent.cdr.cdr.create", mock.MatchedBy(func(e *epb.CDRReported) bool {
		return e.Imsi == cdr.Imsi
	})).Return(nil).Once()
	_, err = s.PostCDR(context.TODO(), req)
	assert.NoError(t, err)
	cdrRepo.AssertExpectations(t)
	usageRepo.AssertExpectations(t)
	assert.NoError(t, err)
}

func TestCDR_InitUsage(t *testing.T) {
	cdrRepo := &mocks.CDRRepo{}
	usageRepo := &mocks.UsageRepo{}
	mbC := &cmocks.MsgBusServiceClient{}
	s, err := NewCDRServer(cdrRepo, usageRepo, OrgId, OrgName, mbC)
	assert.NoError(t, err)

	usageRepo.On("Get", usage.Imsi).Return(nil, gorm.ErrRecordNotFound).Once()
	usageRepo.On("Add", mock.MatchedBy(func(u *db.Usage) bool {
		return u.Imsi == cdr.Imsi
	})).Return(nil).Once()

	err = s.InitUsage(cdr.Imsi, cdr.Policy)
	assert.NoError(t, err)
	usageRepo.AssertExpectations(t)
	assert.NoError(t, err)
}

func TestCDR_GetCDR(t *testing.T) {
	cdrRepo := &mocks.CDRRepo{}
	usageRepo := &mocks.UsageRepo{}
	mbC := &cmocks.MsgBusServiceClient{}
	s, err := NewCDRServer(cdrRepo, usageRepo, OrgId, OrgName, mbC)
	assert.NoError(t, err)

	req := &pb.RecordReq{
		Imsi:      cdr.Imsi,
		StartTime: cdr.StartTime,
		EndTime:   cdr.EndTime,
	}

	cdrRepo.On("GetByFilters", req.Imsi, req.SessionId, req.Policy, req.StartTime, req.EndTime).Return(&[]db.CDR{cdr}, nil).Once()

	resp, err := s.GetCDR(context.TODO(), req)
	assert.NoError(t, err)

	assert.Equal(t, cdr.Imsi, resp.Cdr[0].Imsi)
	assert.Equal(t, cdr.Session, resp.Cdr[0].Session)
	usageRepo.AssertExpectations(t)
	assert.NoError(t, err)
}

func TestCDR_GetUsage(t *testing.T) {
	cdrRepo := &mocks.CDRRepo{}
	usageRepo := &mocks.UsageRepo{}
	mbC := &cmocks.MsgBusServiceClient{}
	s, err := NewCDRServer(cdrRepo, usageRepo, OrgId, OrgName, mbC)
	assert.NoError(t, err)

	usageRepo.On("Get", cdr.Imsi).Return(&usage, nil).Once()

	req := &pb.UsageReq{
		Imsi: usage.Imsi,
	}
	resp, err := s.GetUsage(context.TODO(), req)
	assert.NoError(t, err)

	assert.Equal(t, usage.Imsi, resp.Imsi)
	assert.Equal(t, usage.Usage, resp.Usage)

	usageRepo.AssertExpectations(t)
	assert.NoError(t, err)
}

func TestCDR_GetUsageDetails(t *testing.T) {
	cdrRepo := &mocks.CDRRepo{}
	usageRepo := &mocks.UsageRepo{}
	mbC := &cmocks.MsgBusServiceClient{}
	s, err := NewCDRServer(cdrRepo, usageRepo, OrgId, OrgName, mbC)
	assert.NoError(t, err)

	usageRepo.On("Get", cdr.Imsi).Return(&usage, nil).Once()

	req := &pb.CycleUsageReq{
		Imsi: usage.Imsi,
	}
	resp, err := s.GetUsageDetails(context.TODO(), req)
	assert.NoError(t, err)

	assert.Equal(t, usage.Imsi, resp.Imsi)
	assert.Equal(t, usage.Usage, resp.Usage)
	assert.Equal(t, usage.Historical, resp.Historical)
	usageRepo.AssertExpectations(t)
	assert.NoError(t, err)
}

func TestCDR_GetUsageForPeriod(t *testing.T) {
	cdrRepo := &mocks.CDRRepo{}
	usageRepo := &mocks.UsageRepo{}
	mbC := &cmocks.MsgBusServiceClient{}
	s, err := NewCDRServer(cdrRepo, usageRepo, OrgId, OrgName, mbC)
	assert.NoError(t, err)

	req := &pb.UsageForPeriodReq{
		Imsi:      cdr.Imsi,
		StartTime: cdr.StartTime,
		EndTime:   cdr.EndTime,
	}

	cdrRepo.On("GetByTime", req.Imsi, req.StartTime, req.EndTime).Return(&[]db.CDR{cdr}, nil).Once()

	resp, err := s.GetUsageForPeriod(context.TODO(), req)
	assert.NoError(t, err)

	assert.Equal(t, cdr.TotalBytes, resp.Usage)

	usageRepo.AssertExpectations(t)
	assert.NoError(t, err)
}

func TestCDR_ResetPackageUsage(t *testing.T) {
	cdrRepo := &mocks.CDRRepo{}
	usageRepo := &mocks.UsageRepo{}
	mbC := &cmocks.MsgBusServiceClient{}
	s, err := NewCDRServer(cdrRepo, usageRepo, OrgId, OrgName, mbC)
	assert.NoError(t, err)

	usageRepo.On("Get", cdr.Imsi).Return(&usage, nil).Once()
	usageRepo.On("Add", mock.MatchedBy(func(u *db.Usage) bool {
		return u.Imsi == usage.Imsi
	})).Return(nil).Once()

	err = s.ResetPackageUsage(cdr.Imsi, cdr.Policy)
	assert.NoError(t, err)

	usageRepo.AssertExpectations(t)
	assert.NoError(t, err)
}
