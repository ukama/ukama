/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukama/ukama/systems/common/grpc"
	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/subscriber/sim-pool/mocks"
	pb "github.com/ukama/ukama/systems/subscriber/sim-pool/pb/gen"
	"github.com/ukama/ukama/systems/subscriber/sim-pool/pkg/db"

	"context"
)

const (
	// Organization name for testing
	OrgName = "testOrg"

	// ICCID values
	TestIccid1 = "1234567890123456789"
	TestIccid2 = "9876543210987654321"

	// MSISDN values
	TestMsisdn1 = "2345678901"
	TestMsisdn2 = "3456789012"

	// SimType values
	TestSimTypeData  = ukama.SimTypeUkamaData
	TestSimTypeVoice = ukama.SimTypeOperatorData

	// SmDpAddress values
	TestSmDpAddress1 = "http://localhost:8080"
	TestSmDpAddress2 = "http://localhost:8081"

	// ActivationCode values
	TestActivationCode1 = "123456"
	TestActivationCode2 = "654321"

	// QR Code values
	TestQrCode1 = "QR123"
	TestQrCode2 = "QR456"

	// Test ID
	TestId = uint64(1)

	// Error messages
	ErrorSimPoolRecordNotFound = "SimPool record not found!"
	ErrorDeletingRecord        = "Error while deleting record!"
	ErrorCreatingSims          = "Error creating sims"
	ErrorFetchingSims          = "Error fetching sims"
	ErrorAddingSims            = "Error adding sims"
	ErrorMessageBus            = "Message bus error"

	// CSV headers
	CsvHeader = "ICCID,MSISDN,SmDpAddress,ActivationCode,QrCode,IsPhysical"
)

func TestGetStats(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := &mocks.SimRepo{}
		msgbusClient := &mbmocks.MsgBusServiceClient{}
		simService := NewSimPoolServer(OrgName, mockRepo, msgbusClient)
		reqMock := &pb.GetStatsRequest{
			SimType: TestSimTypeData.String(),
		}
		mockRepo.On("GetSimsByType", mock.Anything).Return([]db.Sim{{
			Iccid:          TestIccid1,
			Msisdn:         TestMsisdn1,
			SimType:        ukama.ParseSimType(TestSimTypeData.String()),
			SmDpAddress:    TestSmDpAddress1,
			IsAllocated:    false,
			ActivationCode: TestActivationCode1,
		}}, nil)
		res, err := simService.GetStats(context.Background(), reqMock)
		assert.NoError(t, err)
		assert.Equal(t, res.Available, uint64(1))
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo := &mocks.SimRepo{}
		msgbusClient := &mbmocks.MsgBusServiceClient{}
		simService := NewSimPoolServer(OrgName, mockRepo, msgbusClient)
		reqMock := &pb.GetStatsRequest{
			SimType: TestSimTypeData.String(),
		}
		mockRepo.On("GetSimsByType", mock.Anything).Return(nil, grpc.SqlErrorToGrpc(errors.New(ErrorSimPoolRecordNotFound), "sim-pool"))
		res, err := simService.GetStats(context.Background(), reqMock)
		assert.Error(t, err)
		assert.Nil(t, res)
		mockRepo.AssertExpectations(t)
	})
}

func TestDelete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := &mocks.SimRepo{}
		msgbusClient := &mbmocks.MsgBusServiceClient{}
		simService := NewSimPoolServer(OrgName, mockRepo, msgbusClient)
		reqMock := &pb.DeleteRequest{
			Id: []uint64{TestId},
		}
		mockRepo.On("Delete", mock.Anything).Return(nil)
		msgbusClient.On("PublishRequest", mock.AnythingOfType("string"), mock.AnythingOfType("*events.SimRemoved")).Return(nil).Once()
		res, err := simService.Delete(context.Background(), reqMock)
		assert.NoError(t, err)
		assert.Equal(t, reqMock.Id[0], res.Id[0])

		mockRepo.AssertExpectations(t)
		msgbusClient.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo := &mocks.SimRepo{}
		msgbusClient := &mbmocks.MsgBusServiceClient{}
		simService := NewSimPoolServer(OrgName, mockRepo, msgbusClient)
		reqMock := &pb.DeleteRequest{
			Id: []uint64{TestId},
		}
		mockRepo.On("Delete", mock.Anything).Return(grpc.SqlErrorToGrpc(errors.New(ErrorDeletingRecord), "sim-pool"))
		res, err := simService.Delete(context.Background(), reqMock)
		assert.Error(t, err)
		assert.Nil(t, res)
		mockRepo.AssertExpectations(t)
	})
}

func TestAdd(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := &mocks.SimRepo{}
		msgbusClient := &mbmocks.MsgBusServiceClient{}
		simService := NewSimPoolServer(OrgName, mockRepo, msgbusClient)
		reqMock := &pb.AddRequest{
			Sim: []*pb.AddSim{
				{
					Iccid:          TestIccid1,
					Msisdn:         TestMsisdn1,
					SimType:        TestSimTypeData.String(),
					SmDpAddress:    TestSmDpAddress1,
					ActivationCode: TestActivationCode1,
					IsPhysical:     false,
				},
			},
		}
		msgbusClient.On("PublishRequest", mock.AnythingOfType("string"), mock.AnythingOfType("*events.SimUploaded")).Return(nil).Once()
		mockRepo.On("Add", mock.Anything).Return(nil)
		res, err := simService.Add(context.Background(), reqMock)
		assert.NoError(t, err)
		assert.Equal(t, reqMock.Sim[0].Iccid, res.Sim[0].Iccid)
		mockRepo.AssertExpectations(t)
		msgbusClient.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo := &mocks.SimRepo{}
		msgbusClient := &mbmocks.MsgBusServiceClient{}
		simService := NewSimPoolServer(OrgName, mockRepo, msgbusClient)
		reqMock := &pb.AddRequest{
			Sim: []*pb.AddSim{
				{
					Iccid:          TestIccid1,
					Msisdn:         TestMsisdn1,
					SimType:        TestSimTypeData.String(),
					SmDpAddress:    TestSmDpAddress1,
					ActivationCode: TestActivationCode1,
					IsPhysical:     false,
				},
			},
		}
		mockRepo.On("Add", mock.Anything).Return(grpc.SqlErrorToGrpc(errors.New(ErrorCreatingSims), "sim-pool"))
		res, err := simService.Add(context.Background(), reqMock)
		assert.Error(t, err)
		assert.Nil(t, res)
		mockRepo.AssertExpectations(t)
	})
}

func TestGet(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := &mocks.SimRepo{}
		msgbusClient := &mbmocks.MsgBusServiceClient{}
		simService := NewSimPoolServer(OrgName, mockRepo, msgbusClient)
		reqMock := &pb.GetRequest{
			IsPhysicalSim: true,
			SimType:       TestSimTypeData.String(),
		}
		mockRepo.On("Get", mock.Anything, mock.Anything).Return(&db.Sim{
			Iccid:          TestIccid1,
			Msisdn:         TestMsisdn1,
			SimType:        ukama.ParseSimType(TestSimTypeData.String()),
			SmDpAddress:    TestSmDpAddress1,
			ActivationCode: TestActivationCode1,
			IsPhysical:     false,
		}, nil)
		res, err := simService.Get(context.Background(), reqMock)
		assert.NoError(t, err)
		assert.Equal(t, TestIccid1, res.Sim.Iccid)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo := &mocks.SimRepo{}
		msgbusClient := &mbmocks.MsgBusServiceClient{}
		simService := NewSimPoolServer(OrgName, mockRepo, msgbusClient)
		reqMock := &pb.GetRequest{
			IsPhysicalSim: true,
			SimType:       TestSimTypeData.String(),
		}
		mockRepo.On("Get", mock.Anything, mock.Anything).Return(nil, grpc.SqlErrorToGrpc(errors.New(ErrorFetchingSims), "sim-pool"))
		res, err := simService.Get(context.Background(), reqMock)
		assert.Error(t, err)
		assert.Nil(t, res)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetByIccid(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := &mocks.SimRepo{}
		msgbusClient := &mbmocks.MsgBusServiceClient{}
		simService := NewSimPoolServer(OrgName, mockRepo, msgbusClient)
		reqMock := &pb.GetByIccidRequest{
			Iccid: TestIccid1,
		}
		mockRepo.On("GetByIccid", reqMock.Iccid).Return(&db.Sim{
			Iccid:          TestIccid1,
			Msisdn:         TestMsisdn1,
			SimType:        ukama.ParseSimType(TestSimTypeData.String()),
			SmDpAddress:    TestSmDpAddress1,
			ActivationCode: TestActivationCode1,
			IsPhysical:     false,
		}, nil)
		res, err := simService.GetByIccid(context.Background(), reqMock)
		assert.NoError(t, err)
		assert.Equal(t, TestIccid1, res.Sim.Iccid)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo := &mocks.SimRepo{}
		msgbusClient := &mbmocks.MsgBusServiceClient{}
		simService := NewSimPoolServer(OrgName, mockRepo, msgbusClient)
		reqMock := &pb.GetByIccidRequest{
			Iccid: TestIccid1,
		}
		mockRepo.On("GetByIccid", mock.Anything).Return(nil, grpc.SqlErrorToGrpc(errors.New(ErrorFetchingSims), "sim-pool"))
		res, err := simService.GetByIccid(context.Background(), reqMock)
		assert.Error(t, err)
		assert.Nil(t, res)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetSims(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := &mocks.SimRepo{}
		msgbusClient := &mbmocks.MsgBusServiceClient{}
		simService := NewSimPoolServer(OrgName, mockRepo, msgbusClient)
		reqMock := &pb.GetSimsRequest{
			SimType: TestSimTypeData.String(),
		}
		mockSims := []db.Sim{
			{
				Iccid:          TestIccid1,
				Msisdn:         TestMsisdn1,
				SimType:        ukama.ParseSimType(TestSimTypeData.String()),
				SmDpAddress:    TestSmDpAddress1,
				IsAllocated:    false,
				ActivationCode: TestActivationCode1,
				IsPhysical:     false,
			},
			{
				Iccid:          TestIccid2,
				Msisdn:         TestMsisdn2,
				SimType:        ukama.ParseSimType(TestSimTypeData.String()),
				SmDpAddress:    TestSmDpAddress1,
				IsAllocated:    true,
				ActivationCode: TestActivationCode2,
				IsPhysical:     true,
			},
		}
		mockRepo.On("GetSims", ukama.ParseSimType(TestSimTypeData.String())).Return(mockSims, nil)
		res, err := simService.GetSims(context.Background(), reqMock)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Len(t, res.Sims, 2)
		assert.Equal(t, TestIccid1, res.Sims[0].Iccid)
		assert.Equal(t, TestIccid2, res.Sims[1].Iccid)
		assert.Equal(t, TestSimTypeData.String(), res.Sims[0].SimType)
		assert.Equal(t, TestSimTypeData.String(), res.Sims[1].SimType)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo := &mocks.SimRepo{}
		msgbusClient := &mbmocks.MsgBusServiceClient{}
		simService := NewSimPoolServer(OrgName, mockRepo, msgbusClient)
		reqMock := &pb.GetSimsRequest{
			SimType: TestSimTypeData.String(),
		}
		mockRepo.On("GetSims", ukama.ParseSimType(TestSimTypeData.String())).Return(nil, grpc.SqlErrorToGrpc(errors.New(ErrorFetchingSims), "sim-pool"))
		res, err := simService.GetSims(context.Background(), reqMock)
		assert.Error(t, err)
		assert.Nil(t, res)
		mockRepo.AssertExpectations(t)
	})

	t.Run("EmptyResult", func(t *testing.T) {
		mockRepo := &mocks.SimRepo{}
		msgbusClient := &mbmocks.MsgBusServiceClient{}
		simService := NewSimPoolServer(OrgName, mockRepo, msgbusClient)
		reqMock := &pb.GetSimsRequest{
			SimType: TestSimTypeVoice.String(),
		}
		mockRepo.On("GetSims", ukama.ParseSimType(TestSimTypeVoice.String())).Return([]db.Sim{}, nil)
		res, err := simService.GetSims(context.Background(), reqMock)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Len(t, res.Sims, 0)
		mockRepo.AssertExpectations(t)
	})
}

func TestUpload(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := &mocks.SimRepo{}
		msgbusClient := &mbmocks.MsgBusServiceClient{}
		simService := NewSimPoolServer(OrgName, mockRepo, msgbusClient)

		csvData := []byte(CsvHeader + "\n" +
			TestIccid1 + "," + TestMsisdn1 + "," + TestSmDpAddress1 + "," + TestActivationCode1 + "," + TestQrCode1 + ",FALSE\n" +
			TestIccid2 + "," + TestMsisdn2 + "," + TestSmDpAddress2 + "," + TestActivationCode2 + "," + TestQrCode2 + ",TRUE")

		reqMock := &pb.UploadRequest{
			SimData: csvData,
			SimType: TestSimTypeData.String(),
		}

		mockRepo.On("Add", mock.AnythingOfType("[]db.Sim")).Return(nil)

		msgbusClient.On("PublishRequest", mock.AnythingOfType("string"), mock.AnythingOfType("*events.SimUploaded")).Return(nil).Once()

		res, err := simService.Upload(context.Background(), reqMock)
		assert.NoError(t, err)
		assert.NotNil(t, res)

		assert.Len(t, res.Iccid, 2)
		assert.Contains(t, res.Iccid, TestIccid1)
		assert.Contains(t, res.Iccid, TestIccid2)

		mockRepo.AssertExpectations(t)
		msgbusClient.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo := &mocks.SimRepo{}
		msgbusClient := &mbmocks.MsgBusServiceClient{}
		simService := NewSimPoolServer(OrgName, mockRepo, msgbusClient)

		csvData := []byte(CsvHeader + "\n" +
			TestIccid1 + "," + TestMsisdn1 + "," + TestSmDpAddress1 + "," + TestActivationCode1 + "," + TestQrCode1 + ",FALSE")

		reqMock := &pb.UploadRequest{
			SimData: csvData,
			SimType: TestSimTypeData.String(),
		}

		mockRepo.On("Add", mock.AnythingOfType("[]db.Sim")).Return(grpc.SqlErrorToGrpc(errors.New(ErrorAddingSims), "sim-pool"))

		res, err := simService.Upload(context.Background(), reqMock)
		assert.Error(t, err)
		assert.Nil(t, res)

		mockRepo.AssertExpectations(t)
	})

	t.Run("EmptyData", func(t *testing.T) {
		mockRepo := &mocks.SimRepo{}
		msgbusClient := &mbmocks.MsgBusServiceClient{}
		simService := NewSimPoolServer(OrgName, mockRepo, msgbusClient)

		csvData := []byte(CsvHeader)

		reqMock := &pb.UploadRequest{
			SimData: csvData,
			SimType: TestSimTypeVoice.String(),
		}

		mockRepo.On("Add", mock.AnythingOfType("[]db.Sim")).Return(nil)

		msgbusClient.On("PublishRequest", mock.AnythingOfType("string"), mock.AnythingOfType("*events.SimUploaded")).Return(nil).Once()

		res, err := simService.Upload(context.Background(), reqMock)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Len(t, res.Iccid, 0)

		mockRepo.AssertExpectations(t)
		msgbusClient.AssertExpectations(t)
	})

	t.Run("InvalidCSVData", func(t *testing.T) {
		mockRepo := &mocks.SimRepo{}
		msgbusClient := &mbmocks.MsgBusServiceClient{}
		simService := NewSimPoolServer(OrgName, mockRepo, msgbusClient)

		invalidData := []byte(`Invalid CSV format`)

		reqMock := &pb.UploadRequest{
			SimData: invalidData,
			SimType: TestSimTypeData.String(),
		}

		mockRepo.On("Add", mock.AnythingOfType("[]db.Sim")).Return(nil)

		msgbusClient.On("PublishRequest", mock.AnythingOfType("string"), mock.AnythingOfType("*events.SimUploaded")).Return(nil).Once()

		res, err := simService.Upload(context.Background(), reqMock)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Len(t, res.Iccid, 0)

		mockRepo.AssertExpectations(t)
		msgbusClient.AssertExpectations(t)
	})

	t.Run("MessageBusError", func(t *testing.T) {
		mockRepo := &mocks.SimRepo{}
		msgbusClient := &mbmocks.MsgBusServiceClient{}
		simService := NewSimPoolServer(OrgName, mockRepo, msgbusClient)

		csvData := []byte(CsvHeader + "\n" +
			TestIccid1 + "," + TestMsisdn1 + "," + TestSmDpAddress1 + "," + TestActivationCode1 + "," + TestQrCode1 + ",FALSE")

		reqMock := &pb.UploadRequest{
			SimData: csvData,
			SimType: TestSimTypeData.String(),
		}

		mockRepo.On("Add", mock.AnythingOfType("[]db.Sim")).Return(nil)

		msgbusClient.On("PublishRequest", mock.AnythingOfType("string"), mock.AnythingOfType("*events.SimUploaded")).Return(errors.New(ErrorMessageBus)).Once()

		res, err := simService.Upload(context.Background(), reqMock)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Len(t, res.Iccid, 1)
		assert.Contains(t, res.Iccid, TestIccid1)
		assert.Equal(t, TestIccid1, res.Iccid[0])

		mockRepo.AssertExpectations(t)
		msgbusClient.AssertExpectations(t)
	})
}
