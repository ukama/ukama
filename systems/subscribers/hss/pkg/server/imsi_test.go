package server_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukama/ukama/systems/registry/hss/mocks"
	pb "github.com/ukama/ukama/systems/registry/hss/pb/gen"
	"github.com/ukama/ukama/systems/registry/hss/pkg/db"
	"github.com/ukama/ukama/systems/registry/hss/pkg/server"
)

const testImsi = "000000111111111111111"
const testOrg = "testOrg"

func Test_Add(t *testing.T) {
	// Arrange
	imsiRepo := mocks.ImsiRepo{}

	imsiRepo.On("Add", testOrg, mock.Anything).Return(nil)

	var actualImsi string
	gutiRepo := mocks.GutiRepo{}

	sub := mocks.HssSubscriber{}
	sub.On("ImsiAdded", testOrg, mock.MatchedBy(func(i *pb.ImsiRecord) bool {
		actualImsi = i.Imsi
		return true
	})).Return().Once()

	is := server.NewImsiService(&imsiRepo, &gutiRepo, server.NewHssEventsSubscribers(&sub))

	// Act
	userID := uuid.New()
	_, err := is.Add(context.TODO(), &pb.AddImsiRequest{
		Imsi: &pb.ImsiRecord{
			Imsi:   testImsi,
			UserId: userID.String(),
			Apn: &pb.Apn{
				Name: "apn",
			},
		},
		Org: testOrg,
	})

	// give it a sec for go routine to finish
	time.Sleep(1 * time.Second)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, testImsi, actualImsi)

	sub.AssertExpectations(t)
	imsiRepo.AssertExpectations(t)
}

func Test_Delete(t *testing.T) {
	// Arrange
	imsiRepo := mocks.ImsiRepo{}
	imsiRepo.On("Delete", testImsi).Return(nil)
	imsiRepo.On("GetByImsi", testImsi).Return(&db.Imsi{Imsi: testImsi, Org: &db.Org{Name: testOrg}}, nil)

	gutiRepo := mocks.GutiRepo{}

	sub := mocks.HssSubscriber{}
	sub.On("ImsiDeleted", testOrg, testImsi).Return().Once()

	is := server.NewImsiService(&imsiRepo, &gutiRepo, server.NewHssEventsSubscribers(&sub))

	// Act
	_, err := is.Delete(context.TODO(), &pb.DeleteImsiRequest{
		IdOneof: &pb.DeleteImsiRequest_Imsi{
			Imsi: testImsi,
		},
	})

	// give it a sec for go routine to finish
	time.Sleep(1 * time.Second)

	// assert
	assert.NoError(t, err)
	sub.AssertExpectations(t)
	imsiRepo.AssertExpectations(t)
}

func Test_DeleteByIccid(t *testing.T) {
	// Arrange
	userID := uuid.New()
	imsiRepo := mocks.ImsiRepo{}
	imsiRepo.On("Delete", testImsi).Return(nil)
	imsiRepo.On("GetImsiByUserUuid", userID).Return([]*db.Imsi{{Imsi: testImsi, UserUuid: userID, Org: &db.Org{Name: testOrg}}}, nil)

	gutiRepo := mocks.GutiRepo{}

	sub := mocks.HssSubscriber{}
	sub.On("ImsiDeleted", testOrg, testImsi).Return().Once()

	is := server.NewImsiService(&imsiRepo, &gutiRepo, server.NewHssEventsSubscribers(&sub))

	// Act
	_, err := is.Delete(context.TODO(), &pb.DeleteImsiRequest{
		IdOneof: &pb.DeleteImsiRequest_UserId{
			UserId: userID.String(),
		},
	})

	// give it a sec for go routine to finish
	time.Sleep(1 * time.Second)

	// assert
	assert.NoError(t, err)
	sub.AssertExpectations(t)
	imsiRepo.AssertExpectations(t)
}

func Test_AddGuti(t *testing.T) {
	// Arrange
	imsiRepo := mocks.ImsiRepo{}
	imsiRepo.On("GetByImsi", testImsi).Return(&db.Imsi{Imsi: testImsi, Org: &db.Org{Name: testOrg}}, nil)

	gutiRepo := mocks.GutiRepo{}
	gutiRepo.On("Update", mock.Anything).Return(nil).Once()

	sub := mocks.HssSubscriber{}
	sub.On("GutiAdded", testOrg, testImsi, mock.Anything).Return().Once()

	is := server.NewImsiService(&imsiRepo, &gutiRepo, server.NewHssEventsSubscribers(&sub))

	// Act
	_, err := is.AddGuti(context.TODO(), &pb.AddGutiRequest{
		Imsi: testImsi,
		Guti: &pb.Guti{},
	})

	// give it a sec for go routine to finish
	time.Sleep(1 * time.Second)

	// assert
	assert.NoError(t, err)
	sub.AssertExpectations(t)
	gutiRepo.AssertExpectations(t)
}

func Test_UpdateTai(t *testing.T) {
	// Arrange
	imsiRepo := mocks.ImsiRepo{}
	imsiRepo.On("GetByImsi", testImsi).Return(&db.Imsi{Imsi: testImsi, Org: &db.Org{Name: testOrg}}, nil)
	imsiRepo.On("UpdateTai", testImsi, mock.Anything).Return(nil).Once()

	gutiRepo := mocks.GutiRepo{}

	sub := mocks.HssSubscriber{}
	sub.On("TaiUpdated", testOrg, mock.Anything).Return().Once()

	is := server.NewImsiService(&imsiRepo, &gutiRepo, server.NewHssEventsSubscribers(&sub))

	// Act
	_, err := is.UpdateTai(context.TODO(), &pb.UpdateTaiRequest{
		Imsi:   testImsi,
		PlmnId: "001111",
		Tac:    99,
	})

	// give it a sec for go routine to finish
	time.Sleep(1 * time.Second)

	// assert
	assert.NoError(t, err)
	sub.AssertExpectations(t)
	gutiRepo.AssertExpectations(t)
}
