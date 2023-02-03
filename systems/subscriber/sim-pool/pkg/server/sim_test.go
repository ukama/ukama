package server

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukama/ukama/systems/common/grpc"
<<<<<<< HEAD
	mbmocks "github.com/ukama/ukama/systems/common/mocks"
=======
>>>>>>> subscriber-sys_sim-manager
	"github.com/ukama/ukama/systems/subscriber/sim-pool/mocks"
	pb "github.com/ukama/ukama/systems/subscriber/sim-pool/pb/gen"
	"github.com/ukama/ukama/systems/subscriber/sim-pool/pkg/db"

	"context"
)

func TestGetStats_Success(t *testing.T) {
	mockRepo := &mocks.SimRepo{}
<<<<<<< HEAD
	msgbusClient := &mbmocks.MsgBusServiceClient{}
	simService := NewSimPoolServer(mockRepo, msgbusClient)
	reqMock := &pb.GetStatsRequest{
		SimType: pb.SimType_INTER_MNO_DATA,
	}
	mockRepo.On("GetSimsByType", mock.Anything).Return([]db.Sim{{
		Iccid:          "1234567890123456789",
		Msisdn:         "2345678901",
		SimType:         "inter_mno_data",
		SmDpAddress:    "http://localhost:8080",
		IsAllocated:     false,
=======
	simService := NewSimServer(mockRepo)
	reqMock := &pb.GetStatsRequest{
		SimType: pb.SimType_INTER_MNO_DATA,
	}
	mockRepo.On("GetStats", mock.Anything).Return([]db.Sim{{
		Iccid:          "1234567890123456789",
		Msisdn:         "2345678901",
		Sim_type:       "inter_mno_data",
		SmDpAddress:    "http://localhost:8080",
		Is_allocated:   false,
>>>>>>> subscriber-sys_sim-manager
		ActivationCode: "123456",
	}}, nil)
	res, err := simService.GetStats(context.Background(), reqMock)
	assert.NoError(t, err)
	assert.Equal(t, res.Available, uint64(1))
}

func TestGetStats_Error(t *testing.T) {
	mockRepo := &mocks.SimRepo{}
<<<<<<< HEAD
	msgbusClient := &mbmocks.MsgBusServiceClient{}
	simService := NewSimPoolServer(mockRepo, msgbusClient)
	reqMock := &pb.GetStatsRequest{
		SimType: pb.SimType_INTER_MNO_DATA,
	}
	mockRepo.On("GetSimsByType", mock.Anything).Return(nil, grpc.SqlErrorToGrpc(errors.New("SimPool record not found!"), "sim-pool"))
=======
	simService := NewSimServer(mockRepo)
	reqMock := &pb.GetStatsRequest{
		SimType: pb.SimType_INTER_MNO_DATA,
	}
	mockRepo.On("GetStats", mock.Anything).Return(nil, grpc.SqlErrorToGrpc(errors.New("SimPool record not found!"), "sim-pool"))
>>>>>>> subscriber-sys_sim-manager
	res, err := simService.GetStats(context.Background(), reqMock)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestDelete_Success(t *testing.T) {
	mockRepo := &mocks.SimRepo{}
<<<<<<< HEAD
	msgbusClient := &mbmocks.MsgBusServiceClient{}
	simService := NewSimPoolServer(mockRepo, msgbusClient)
=======
	simService := NewSimServer(mockRepo)
>>>>>>> subscriber-sys_sim-manager
	reqMock := &pb.DeleteRequest{
		Id: []uint64{1},
	}
	mockRepo.On("Delete", mock.Anything).Return(nil)
	res, err := simService.Delete(context.Background(), reqMock)
	assert.NoError(t, err)
	assert.Equal(t, reqMock.Id[0], res.Id[0])
}

func TestDelete_Error(t *testing.T) {
	mockRepo := &mocks.SimRepo{}
<<<<<<< HEAD
	msgbusClient := &mbmocks.MsgBusServiceClient{}
	simService := NewSimPoolServer(mockRepo, msgbusClient)
=======
	simService := NewSimServer(mockRepo)
>>>>>>> subscriber-sys_sim-manager
	reqMock := &pb.DeleteRequest{
		Id: []uint64{1},
	}
	mockRepo.On("Delete", mock.Anything).Return(grpc.SqlErrorToGrpc(errors.New("Error while deleting record!"), "sim-pool"))
	res, err := simService.Delete(context.Background(), reqMock)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestAdd_Success(t *testing.T) {
	mockRepo := &mocks.SimRepo{}
<<<<<<< HEAD
	msgbusClient := &mbmocks.MsgBusServiceClient{}
	simService := NewSimPoolServer(mockRepo, msgbusClient)
=======
	simService := NewSimServer(mockRepo)
>>>>>>> subscriber-sys_sim-manager
	reqMock := &pb.AddRequest{
		Sim: []*pb.AddSim{
			{
				Iccid:          "1234567890123456789",
				Msisdn:         "2345678901",
				SimType:        pb.SimType_INTER_MNO_DATA,
				SmDpAddress:    "http://localhost:8080",
				ActivationCode: "123456",
				IsPhysical:     false,
			},
		},
	}
	mockRepo.On("Add", mock.Anything).Return(nil)
	res, err := simService.Add(context.Background(), reqMock)
	assert.NoError(t, err)
	assert.Equal(t, reqMock.Sim[0].Iccid, res.Sim[0].Iccid)
}

func TestAdd_Error(t *testing.T) {
	mockRepo := &mocks.SimRepo{}
<<<<<<< HEAD
	msgbusClient := &mbmocks.MsgBusServiceClient{}
	simService := NewSimPoolServer(mockRepo, msgbusClient)
=======
	simService := NewSimServer(mockRepo)
>>>>>>> subscriber-sys_sim-manager
	reqMock := &pb.AddRequest{
		Sim: []*pb.AddSim{
			{
				Iccid:          "1234567890123456789",
				Msisdn:         "2345678901",
				SimType:        pb.SimType_INTER_MNO_DATA,
				SmDpAddress:    "http://localhost:8080",
				ActivationCode: "123456",
				IsPhysical:     false,
			},
		},
	}
	mockRepo.On("Add", mock.Anything).Return(grpc.SqlErrorToGrpc(errors.New("Error creating sims"), "sim-pool"))
	res, err := simService.Add(context.Background(), reqMock)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestGet_Success(t *testing.T) {
	mockRepo := &mocks.SimRepo{}
<<<<<<< HEAD
	msgbusClient := &mbmocks.MsgBusServiceClient{}
	simService := NewSimPoolServer(mockRepo, msgbusClient)
=======
	simService := NewSimServer(mockRepo)
>>>>>>> subscriber-sys_sim-manager
	reqMock := &pb.GetRequest{
		IsPhysicalSim: true,
		SimType:       pb.SimType_INTER_MNO_DATA,
	}
	mockRepo.On("Get", mock.Anything, mock.Anything).Return(&db.Sim{
		Iccid:          "1234567890123456789",
		Msisdn:         "2345678901",
<<<<<<< HEAD
		SimType:         "inter_mno_data",
		SmDpAddress:    "http://localhost:8080",
		ActivationCode: "123456",
		IsPhysical:      false,
=======
		Sim_type:       "inter_mno_data",
		SmDpAddress:    "http://localhost:8080",
		ActivationCode: "123456",
		Is_physical:    false,
>>>>>>> subscriber-sys_sim-manager
	}, nil)
	res, err := simService.Get(context.Background(), reqMock)
	assert.NoError(t, err)
	assert.Equal(t, "1234567890123456789", res.Sim.Iccid)
}

func TestGet_Error(t *testing.T) {
	mockRepo := &mocks.SimRepo{}
<<<<<<< HEAD
	msgbusClient := &mbmocks.MsgBusServiceClient{}
	simService := NewSimPoolServer(mockRepo, msgbusClient)
=======
	simService := NewSimServer(mockRepo)
>>>>>>> subscriber-sys_sim-manager
	reqMock := &pb.GetRequest{
		IsPhysicalSim: true,
		SimType:       pb.SimType_INTER_MNO_DATA,
	}
	mockRepo.On("Get", mock.Anything, mock.Anything).Return(nil, grpc.SqlErrorToGrpc(errors.New("Error fetching sims"), "sim-pool"))
	res, err := simService.Get(context.Background(), reqMock)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestGetByIccid_Success(t *testing.T) {
	mockRepo := &mocks.SimRepo{}
<<<<<<< HEAD
	msgbusClient := &mbmocks.MsgBusServiceClient{}
	simService := NewSimPoolServer(mockRepo, msgbusClient)
=======
	simService := NewSimServer(mockRepo)
>>>>>>> subscriber-sys_sim-manager
	reqMock := &pb.GetByIccidRequest{
		Iccid: "1234567890123456789",
	}
	mockRepo.On("GetByIccid", reqMock.Iccid).Return(&db.Sim{
		Iccid:          "1234567890123456789",
		Msisdn:         "2345678901",
<<<<<<< HEAD
		SimType:         "inter_mno_data",
		SmDpAddress:    "http://localhost:8080",
		ActivationCode: "123456",
		IsPhysical:      false,
=======
		Sim_type:       "inter_mno_data",
		SmDpAddress:    "http://localhost:8080",
		ActivationCode: "123456",
		Is_physical:    false,
>>>>>>> subscriber-sys_sim-manager
	}, nil)
	res, err := simService.GetByIccid(context.Background(), reqMock)
	assert.NoError(t, err)
	assert.Equal(t, "1234567890123456789", res.Sim.Iccid)
}

func TestGetByIccid_Error(t *testing.T) {
	mockRepo := &mocks.SimRepo{}
<<<<<<< HEAD
	msgbusClient := &mbmocks.MsgBusServiceClient{}
	simService := NewSimPoolServer(mockRepo, msgbusClient)
=======
	simService := NewSimServer(mockRepo)
>>>>>>> subscriber-sys_sim-manager
	reqMock := &pb.GetByIccidRequest{
		Iccid: "1234567890123456789",
	}
	mockRepo.On("GetByIccid", mock.Anything).Return(nil, grpc.SqlErrorToGrpc(errors.New("Error fetching sims"), "sim-pool"))
	res, err := simService.GetByIccid(context.Background(), reqMock)
	assert.Error(t, err)
	assert.Nil(t, res)
}
