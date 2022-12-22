package server

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/subscriber/simPool/mocks"
	pb "github.com/ukama/ukama/systems/subscriber/simPool/pb/gen"
	"github.com/ukama/ukama/systems/subscriber/simPool/pkg/db"

	"context"
)

func TestGetStats_Success(t *testing.T) {
	mockRepo := &mocks.SimRepo{}
	simService := NewSimServer(mockRepo)
	reqMock := &pb.GetStatsRequest{
		SimType: pb.SimType_inter_mno_data,
	}
	mockRepo.On("GetStats", mock.Anything).Return([]db.Sim{{
		Iccid:          "1234567890123456789",
		Msisdn:         "2345678901",
		QrCode:         "http://localhost:8080",
		Sim_type:       "inter_mno_data",
		SmDpAddress:    "http://localhost:8080",
		Is_allocated:   false,
		ActivationCode: "123456",
	}}, nil)
	res, err := simService.GetStats(context.Background(), reqMock)
	assert.NoError(t, err)
	assert.Equal(t, res.Available, uint64(1))
}

func TestGetStats_Error(t *testing.T) {
	mockRepo := &mocks.SimRepo{}
	simService := NewSimServer(mockRepo)
	reqMock := &pb.GetStatsRequest{
		SimType: pb.SimType_inter_mno_data,
	}
	mockRepo.On("GetStats", mock.Anything).Return(nil, grpc.SqlErrorToGrpc(errors.New("SimPool record not found!"), "simPool"))
	res, err := simService.GetStats(context.Background(), reqMock)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestDelete_Success(t *testing.T) {
	mockRepo := &mocks.SimRepo{}
	simService := NewSimServer(mockRepo)
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
	simService := NewSimServer(mockRepo)
	reqMock := &pb.DeleteRequest{
		Id: []uint64{1},
	}
	mockRepo.On("Delete", mock.Anything).Return(grpc.SqlErrorToGrpc(errors.New("Error while deleting record!"), "simPool"))
	res, err := simService.Delete(context.Background(), reqMock)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestAdd_Success(t *testing.T) {
	mockRepo := &mocks.SimRepo{}
	simService := NewSimServer(mockRepo)
	reqMock := &pb.AddRequest{
		Sim: []*pb.AddSim{
			{
				Iccid:          "1234567890123456789",
				Msisdn:         "2345678901",
				QrCode:         "http://localhost:8080",
				SimType:        pb.SimType_inter_mno_data,
				SmDpAddress:    "http://localhost:8080",
				ActivationCode: "123456",
				IsPhysicalSim:  false,
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
	simService := NewSimServer(mockRepo)
	reqMock := &pb.AddRequest{
		Sim: []*pb.AddSim{
			{
				Iccid:          "1234567890123456789",
				Msisdn:         "2345678901",
				QrCode:         "http://localhost:8080",
				SimType:        pb.SimType_inter_mno_data,
				SmDpAddress:    "http://localhost:8080",
				ActivationCode: "123456",
				IsPhysicalSim:  false,
			},
		},
	}
	mockRepo.On("Add", mock.Anything).Return(grpc.SqlErrorToGrpc(errors.New("Error creating sims"), "simPool"))
	res, err := simService.Add(context.Background(), reqMock)
	assert.Error(t, err)
	assert.Nil(t, res)
}
